package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"ozon/internal/client/graph"
	"strconv"
	"syscall"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"ozon/internal/config"
	"ozon/pkg/logger/logsetup"
	"ozon/pkg/logger/sl"
)

const defaultPort = "8080"

func main() {
	const op = "cmd/client.main"
	cfg := config.MustLoad()

	log := logsetup.SetupLogger(cfg.Env)
	log.Info("starting client", slog.Any("config", cfg.Clients.GRPCClient))

	cc, err := NewConnection(
		context.Background(),
		log,
		cfg.GRPC.Host+":"+strconv.Itoa(cfg.GRPC.Port),
		cfg.Clients.GRPCClient.RetriesCount,
		cfg.Clients.GRPCClient.Timeout,
	)
	if err != nil {
		log.Error("failed to connect to server", op, sl.Err(err))
		panic(err)
	}
	defer cc.Close()

	resolver, err := graph.New(cfg.Postgres, cc, log)
	if err != nil {
		log.Error("failed to create resolver", op, sl.Err(err))
		panic(err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info("", slog.String("connect to http://localhost:%s/ for GraphQL playground", "8080"))

	err = http.ListenAndServe(":"+"8080", nil)
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sgl := <-stop
	log.Info("stopping client", slog.String("signal", sgl.String()))
	log.Info("client stopped")
}

func StringToTime(s string) time.Time {
	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}

func NewConnection(ctx context.Context, log *slog.Logger, addr string, retriesCount int, timeout time.Duration) (*grpc.ClientConn, error) {
	const op = "internal/clients/blog/NewConnection"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cc, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}
