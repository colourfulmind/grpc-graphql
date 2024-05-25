package app

import (
	"fmt"
	"log/slog"

	grpcserver "ozon/internal/app/grpc"
	"ozon/internal/config"
	"ozon/internal/services/comments"
	"ozon/internal/services/posts"
	"ozon/internal/services/sso"
	resolvers "ozon/internal/storage/graphql"
	"ozon/pkg/logger/sl"
)

type App struct {
	Server *grpcserver.App
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	const op = "internal/app.New"
	storage, err := resolvers.New(cfg.Postgres, log)
	if err != nil {
		log.Error("error occurred while connecting to database", op, sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ssoService := sso.New(log, storage, storage, storage, cfg.TokenTTL)
	postService := posts.New(log, storage, storage)
	commentService := comments.New(log, storage, storage)

	return &App{
		Server: grpcserver.New(log, ssoService, postService, commentService, cfg.GRPC.Port),
	}, nil
}
