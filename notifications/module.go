package notifications

import (
	"context"
	"mall/internal/monolith"
	"mall/notifications/internal/application"
	"mall/notifications/internal/logging"
	"mall/notifications/internal/rpc"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	conn, err := rpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	customers := rpc.NewCustomerRepository(conn)

	// setup application
	var app application.App
	app = application.New(customers)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup driver adapters
	if err := rpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}

	return nil
}
