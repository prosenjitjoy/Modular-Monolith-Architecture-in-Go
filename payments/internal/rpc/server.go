package rpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"mall/payments/internal/application"
	"mall/payments/paymentspb"
)

type server struct {
	app application.App
	paymentspb.UnimplementedPaymentsServiceServer
}

var _ paymentspb.PaymentsServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	paymentspb.RegisterPaymentsServiceServer(registrar, server{app: app})

	return nil
}

func (s server) AuthorizePayment(ctx context.Context, request *paymentspb.AuthorizePaymentRequest) (*paymentspb.AuthorizePaymentResponse, error) {
	id := uuid.NewString()
	err := s.app.AuthorizePayment(ctx, application.AuthorizePayment{
		ID:         id,
		CustomerID: request.GetCustomerId(),
		Amount:     request.GetAmount(),
	})
	if err != nil {
		return nil, err
	}

	return &paymentspb.AuthorizePaymentResponse{Id: id}, nil
}

func (s server) ConfirmPayment(ctx context.Context, request *paymentspb.ConfirmPaymentRequest) (*paymentspb.ConfirmPaymentResponse, error) {
	err := s.app.ConfirmPayment(ctx, application.ConfirmPayment{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &paymentspb.ConfirmPaymentResponse{}, nil
}

func (s server) CreateInvoice(ctx context.Context, request *paymentspb.CreateInvoiceRequest) (*paymentspb.CreateInvoiceResponse, error) {
	id := uuid.NewString()

	err := s.app.CreateInvoice(ctx, application.CreateInvoice{
		ID:      id,
		OrderID: request.GetOrderId(),
		Amount:  request.GetAmount(),
	})
	if err != nil {
		return nil, err
	}

	return &paymentspb.CreateInvoiceResponse{Id: id}, nil
}

func (s server) AdjustInvoice(ctx context.Context, request *paymentspb.AdjustInvoiceRequest) (*paymentspb.AdjustInvoiceResponse, error) {
	err := s.app.AdjustInvoice(ctx, application.AdjustInvoice{
		ID:     request.GetId(),
		Amount: request.GetAmount(),
	})
	if err != nil {
		return nil, err
	}

	return &paymentspb.AdjustInvoiceResponse{}, nil
}

func (s server) PayInvoice(ctx context.Context, request *paymentspb.PayInvoiceRequest) (*paymentspb.PayInvoiceResponse, error) {
	err := s.app.PayInvoice(ctx, application.PayInvoice{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &paymentspb.PayInvoiceResponse{}, nil
}

func (s server) CancelInvoice(ctx context.Context, request *paymentspb.CancelInvoiceRequest) (*paymentspb.CancelInvoiceResponse, error) {
	err := s.app.CancelInvoice(ctx, application.CancelInvoice{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &paymentspb.CancelInvoiceResponse{}, nil
}
