package graph

import (
	"database/sql"
	"log/slog"
	"ozon/protos/gen/go/comments"
	"ozon/protos/gen/go/posts"
	"ozon/protos/gen/go/sso"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB       *sql.DB
	SSO      sso.SSOClient
	Posts    posts.PostsClient
	Comments comments.CommentsClient
	log      *slog.Logger
}
