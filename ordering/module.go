package ordering

import (
	"context"
	"mall/internal/monolith"
	"mall/ordering/internal/application"
	"mall/ordering/internal/logging"
	"mall/ordering/internal/postgres"
	"mall/ordering/internal/rest"
	"mall/ordering/internal/rpc"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	orders := postgres.NewOrderRepository("ordering.orders", mono.DB())
	conn, err := rpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	customers := rpc.NewCustomerRepository(conn)
	payments := rpc.NewPaymentRepository(conn)
	invoices := rpc.NewInvoiceRepository(conn)
	shopping := rpc.NewShoppingListRepository(conn)
	notifications := rpc.NewNotificationRepository(conn)

	// setup application
	var app application.App
	app = application.New(orders, customers, payments, invoices, shopping, notifications)
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
