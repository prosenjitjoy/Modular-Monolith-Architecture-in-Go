package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mall/internal/config"
	"mall/internal/monolith"
	"mall/internal/waiter"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type app struct {
	cfg     *config.AppConfig
	db      *sql.DB
	logger  *slog.Logger
	mux     *chi.Mux
	rpc     *grpc.Server
	waiter  waiter.Waiter
	modules []monolith.Module
}

func (a *app) Config() *config.AppConfig {
	return a.cfg
}

func (a *app) DB() *sql.DB {
	return a.db
}

func (a *app) Logger() *slog.Logger {
	return a.logger
}

func (a *app) Mux() *chi.Mux {
	return a.mux
}

func (a *app) RPC() *grpc.Server {
	return a.rpc
}

func (a *app) Waiter() waiter.Waiter {
	return a.waiter
}

func (a *app) startupModules() error {
	for _, module := range a.modules {
		if err := module.Startup(a.waiter.Context(), a); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) waitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    a.cfg.Web.Address(),
		Handler: a.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		slog.Info("web server started", "address", a.cfg.Web.Address())
		defer slog.Info("web server shutdown")

		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		slog.Info("web server to be shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
		defer cancel()

		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	})

	return group.Wait()
}

func (a *app) waitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		slog.Info("rpc server started", "address", a.cfg.Rpc.Address())
		defer slog.Info("rpc server shutdown")

		if err := a.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		slog.Info("rpc server to be shutdown")

		stopped := make(chan struct{})
		go func() {
			a.RPC().GracefulStop()
			close(stopped)
		}()

		timeout := time.NewTimer(a.cfg.Timeout)
		select {
		case <-timeout.C:
			a.RPC().Stop()
			return errors.New("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}
