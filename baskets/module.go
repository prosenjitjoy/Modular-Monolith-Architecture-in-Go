package baskets

import (
	"context"
	"mall/baskets/internal/application"
	"mall/baskets/internal/logging"
	"mall/baskets/internal/postgres"
	"mall/baskets/internal/rest"
	"mall/baskets/internal/rpc"
	"mall/internal/monolith"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	//setup driven adapters
	baskets := postgres.NewBasketRepository("baskets.baskets", mono.DB())
	conn, err := rpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	stores := rpc.NewStoreRepository(conn)
	products := rpc.NewProductRepository(conn)
	orders := rpc.NewOrderRepository(conn)

	// setup application
	var app application.App
	app = application.New(baskets, stores, products, orders)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup driver adapters
	if err := rpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}

	return nil
}
