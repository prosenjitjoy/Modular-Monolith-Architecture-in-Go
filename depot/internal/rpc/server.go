package rpc

import (
	"context"
	"mall/depot/depotpb"
	"mall/depot/internal/application"
	"mall/depot/internal/domain"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	depotpb.UnimplementedDepotServiceServer
}

var _ depotpb.DepotServiceServer = (*server)(nil)

func Register(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	depotpb.RegisterDepotServiceServer(registrar, server{app: app})

	return nil
}

func (s server) CreateShoppingList(ctx context.Context, request *depotpb.CreateShoppingListRequest) (*depotpb.CreateShoppingListResponse, error) {
	id := uuid.NewString()

	items := make([]application.OrderItem, 0, len(request.GetItems()))
	for _, item := range request.GetItems() {
		items = append(items, s.itemToDomain(item))
	}

	err := s.app.CreateShoppingList(ctx, application.CreateShoppingList{
		ID:      id,
		OrderID: request.GetOrderId(),
		Items:   items,
	})
	if err != nil {
		return nil, err
	}

	return &depotpb.CreateShoppingListResponse{Id: id}, nil
}

func (s server) CancelShoppingList(ctx context.Context, request *depotpb.CancelShoppingListRequest) (*depotpb.CancelShoppingListResponse, error) {
	err := s.app.CancelShoppingList(ctx, application.CancelShoppingList{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &depotpb.CancelShoppingListResponse{}, nil
}

func (s server) AssignShoppingList(ctx context.Context, request *depotpb.AssignShoppingListRequest) (*depotpb.AssignShoppingListResponse, error) {
	err := s.app.AssignShoppingList(ctx, application.AssignShoppingList{
		ID:    request.GetId(),
		BotID: request.GetBotId(),
	})
	if err != nil {
		return nil, err
	}

	return &depotpb.AssignShoppingListResponse{}, nil
}

func (s server) CompleteShoppingList(ctx context.Context, request *depotpb.CompleteShoppingListRequest) (*depotpb.CompleteShoppingListResponse, error) {
	err := s.app.CompleteShoppingList(ctx, application.CompleteShoppingList{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &depotpb.CompleteShoppingListResponse{}, nil
}

func (s server) GetShoppingList(ctx context.Context, request *depotpb.GetShoppingListRequest) (*depotpb.GetShoppingListResponse, error) {
	shoppingList, err := s.app.GetShoppingList(ctx, application.GetShoppingList{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &depotpb.GetShoppingListResponse{ShoppingList: s.listFromDomain(shoppingList)}, nil
}

func (s server) itemToDomain(item *depotpb.OrderItem) application.OrderItem {
	return application.OrderItem{
		StoreID:   item.GetStoreId(),
		ProductID: item.GetProductId(),
		Quantity:  int(item.GetQuantity()),
	}
}

func (s server) listFromDomain(list *domain.ShoppingList) *depotpb.ShoppingList {
	stops := make(map[string]*depotpb.Stop)
	for storeID, stop := range list.Stops {
		items := make(map[string]*depotpb.Item)
		for productID, item := range stop.Items {
			items[productID] = &depotpb.Item{
				Name:     item.ProductName,
				Quantity: int32(item.Quantity),
			}
		}

		stops[storeID] = &depotpb.Stop{
			StoreName:     stop.StoreName,
			StoreLocation: stop.StoreLocation,
			Items:         items,
		}
	}

	return &depotpb.ShoppingList{
		Id:            list.ID,
		OrderId:       list.OrderID,
		Stops:         stops,
		AssignedBotId: list.AssignedBotID,
		Status:        list.Status.String(),
	}
}
