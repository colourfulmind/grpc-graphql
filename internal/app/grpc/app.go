package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"

	"ozon/internal/grpc/comments"
	"ozon/internal/grpc/posts"
	"ozon/internal/grpc/sso"
)

type App struct {
	log        *slog.Logger
	GRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, ssoService sso.SSO, postService posts.Posts, commentsService comments.Comments, port int) *App {
	GRPCSever := grpc.NewServer()
	sso.Register(GRPCSever, ssoService)
	posts.Register(GRPCSever, postService)
	comments.Register(GRPCSever, commentsService)
	reflection.Register(GRPCSever)
	return &App{
		log:        log,
		GRPCServer: GRPCSever,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "internal/app/grpc.Run"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err = a.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "internal/app/grpc.Stop"

	a.log.With(slog.String("op", op)).Info(
		"stopping gRPC server", slog.Int("port", a.port),
	)

	a.GRPCServer.GracefulStop()
}
