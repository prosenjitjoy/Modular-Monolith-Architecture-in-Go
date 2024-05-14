package rpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
	"mall/ordering/orderingpb"
)

type server struct {
	app application.App
	orderingpb.UnimplementedOrderingServiceServer
}

var _ orderingpb.OrderingServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	orderingpb.RegisterOrderingServiceServer(registrar, server{app: app})

	return nil
}

func (s server) CreateOrder(ctx context.Context, request *orderingpb.CreateOrderRequest) (*orderingpb.CreateOrderResponse, error) {
	id := uuid.NewString()

	items := make([]*domain.Item, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, s.itemToDomain(item))
	}

	err := s.app.CreateOrder(ctx, application.CreateOrder{
		ID:         id,
		CustomerID: request.GetCustomerId(),
		PaymentID:  request.GetPaymentId(),
		Items:      items,
	})
	if err != nil {
		return nil, err
	}

	return &orderingpb.CreateOrderResponse{Id: id}, nil
}

func (s server) CancelOrder(ctx context.Context, request *orderingpb.CancelOrderRequest) (*orderingpb.CancelOrderResponse, error) {
	err := s.app.CancelOrder(ctx, application.CancelOrder{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &orderingpb.CancelOrderResponse{}, nil
}

func (s server) ReadyOrder(ctx context.Context, request *orderingpb.ReadyOrderRequest) (*orderingpb.ReadyOrderResponse, error) {
	err := s.app.ReadyOrder(ctx, application.ReadyOrder{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &orderingpb.ReadyOrderResponse{}, nil
}

func (s server) CompleteOrder(ctx context.Context, request *orderingpb.CompleteOrderRequest) (*orderingpb.CompleteOrderResponse, error) {
	err := s.app.CompleteOrder(ctx, application.CompleteOrder{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &orderingpb.CompleteOrderResponse{}, nil
}

func (s server) GetOrder(ctx context.Context, request *orderingpb.GetOrderRequest) (*orderingpb.GetOrderResponse, error) {
	order, err := s.app.GetOrder(ctx, application.GetOrder{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &orderingpb.GetOrderResponse{Order: s.orderFromDomain(order)}, nil
}

func (s server) orderFromDomain(order *domain.Order) *orderingpb.Order {
	items := make([]*orderingpb.Item, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, s.itemFromDomain(item))
	}

	return &orderingpb.Order{
		Id:         order.ID,
		CustomerId: order.CustomerID,
		PaymentId:  order.PaymentID,
		Items:      items,
		Status:     order.Status.String(),
	}
}

func (s server) itemToDomain(item *orderingpb.Item) *domain.Item {
	return &domain.Item{
		ProductID:    item.GetProductId(),
		StoreID:      item.GetStoreId(),
		StoreName:    item.GetStoreName(),
		ProductName:  item.GetProductName(),
		ProductPrice: item.GetProductPrice(),
		Quantity:     int(item.GetQuantity()),
	}
}

func (s server) itemFromDomain(item *domain.Item) *orderingpb.Item {
	return &orderingpb.Item{
		StoreId:      item.StoreID,
		ProductId:    item.ProductID,
		StoreName:    item.StoreName,
		ProductName:  item.ProductName,
		ProductPrice: item.ProductPrice,
		Quantity:     int32(item.Quantity),
	}
}
