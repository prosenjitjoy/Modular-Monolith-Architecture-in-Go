package depot

import (
	"context"
	"mall/depot/internal/application"
	"mall/depot/internal/logging"
	"mall/depot/internal/postgres"
	"mall/depot/internal/rest"
	"mall/depot/internal/rpc"
	"mall/internal/monolith"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	shoppingLists := postgres.NewShoppingListRepository("depot.shopping_lists", mono.DB())
	conn, err := rpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	stores := rpc.NewStoreRepository(conn)
	products := rpc.NewProductRepository(conn)
	orders := rpc.NewOrderRepository(conn)

	// setup application
	var app application.App
	app = application.New(shoppingLists, stores, products, orders)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup driver adapters
	if err := rpc.Register(ctx, app, mono.RPC()); err != nil {
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}

	return nil
}
