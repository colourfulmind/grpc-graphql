package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"ozon/internal/client/graphql"
	"ozon/internal/client/graphql/graph"
	"strconv"
	"syscall"

	"ozon/internal/config"
	"ozon/pkg/logger/logsetup"
	"ozon/pkg/logger/sl"
)

func main() {
	const op = "cmd/test_client.main"
	cfg := config.MustLoad()

	log := logsetup.SetupLogger(cfg.Env)
	log.Info("starting client", slog.Any("config", cfg.Clients.GRPCClient))

	cc, err := graphql.NewConnection(
		context.Background(),
		log,
		":"+strconv.Itoa(cfg.GRPC.Port),
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

	log.Info("", slog.String("connect to http://localhost:%s/ for GraphQL playground", cfg.Clients.GRPCClient.Port))

	go http.ListenAndServe(":"+cfg.Clients.GRPCClient.Port, nil)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sgl := <-stop
	log.Info("stopping test_client", slog.String("signal", sgl.String()))
	log.Info("test_client stopped")
}
