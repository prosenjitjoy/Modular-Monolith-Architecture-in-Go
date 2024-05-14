package monolith

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/internal/config"
	"mall/internal/waiter"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type Monolith interface {
	Config() *config.AppConfig
	DB() *sql.DB
	Logger() *slog.Logger
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
}

type Module interface {
	Startup(ctx context.Context, mono Monolith) error
}
