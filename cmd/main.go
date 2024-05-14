package main

import (
	"database/sql"
	"log/slog"
	"mall/baskets"
	"mall/customers"
	"mall/database/migrations"
	"mall/depot"
	"mall/internal/config"
	"mall/internal/logger"
	"mall/internal/monolith"
	"mall/internal/waiter"
	"mall/notifications"
	"mall/ordering"
	"mall/payments"
	"mall/stores"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// parse configuration
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	// connect database
	db, err := sql.Open("pgx", cfg.Postgres.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	// migrate database
	err = migrateDB(db)
	if err != nil {
		return err
	}

	// init grpc
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// init infrastructure
	m := app{
		cfg:    cfg,
		db:     db,
		logger: logger.New(cfg.LogLevel),
		mux:    chi.NewMux(),
		rpc:    grpcServer,
		waiter: waiter.New(waiter.CatchSignals()),
	}

	// init modules
	m.modules = []monolith.Module{
		&baskets.Module{},
		&customers.Module{},
		&depot.Module{},
		&notifications.Module{},
		&ordering.Module{},
		&payments.Module{},
		&stores.Module{},
	}

	// start modules
	if err := m.startupModules(); err != nil {
		return err
	}

	slog.Info("started mallbots application")
	defer slog.Info("stopped mallbots application")

	m.waiter.Add(
		m.waitForWeb,
		m.waitForRPC,
	)

	return m.waiter.Wait()
}

func migrateDB(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}
