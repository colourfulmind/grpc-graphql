package graph

import (
	"database/sql"
	"log/slog"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB  *sql.DB
	log *slog.Logger
}

const timeout = 30 * time.Second
