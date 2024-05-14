package stores

import (
	"context"
	"mall/internal/monolith"
	"mall/stores/internal/application"
	"mall/stores/internal/logging"
	"mall/stores/internal/postgres"
	"mall/stores/internal/rest"
	"mall/stores/internal/rpc"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	stores := postgres.NewStoreRepository("stores.stores", mono.DB())
	participatingStore := postgres.NewParticipatingStoreRepository("stores.stores", mono.DB())
	products := postgres.NewProductRepository("stores.products", mono.DB())

	// setup application
	var app application.App
	app = application.New(stores, participatingStore, products)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup driver adapters
	if err := rpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}

	return nil
}
