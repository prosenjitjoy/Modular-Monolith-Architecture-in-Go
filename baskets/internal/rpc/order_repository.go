package rpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"mall/baskets/internal/domain"
	"mall/ordering/orderingpb"
)

type OrderRepository struct {
	client orderingpb.OrderingServiceClient
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(conn *grpc.ClientConn) OrderRepository {
	return OrderRepository{client: orderingpb.NewOrderingServiceClient(conn)}
}

func (r OrderRepository) Save(ctx context.Context, basket *domain.Basket) (string, error) {
	items := make([]*orderingpb.Item, 0, len(basket.Items))
	for _, item := range basket.Items {
		items = append(items, &orderingpb.Item{
			StoreId:      item.StoreID,
			ProductId:    item.ProductID,
			StoreName:    item.StoreName,
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			Quantity:     int32(item.Quantity),
		})
	}

	resp, err := r.client.CreateOrder(ctx, &orderingpb.CreateOrderRequest{
		Items:      items,
		CustomerId: basket.CustomerID,
		PaymentId:  basket.PaymentID,
	})
	if err != nil {
		return "", fmt.Errorf("saving order: %w", err)
	}

	return resp.GetId(), nil
}
