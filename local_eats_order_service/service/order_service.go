package service

import (
	"context"
	pb "order_serive/genproto"
	"order_serive/storage"
)

type OrderService struct {
	InitRepo storage.InitRoot
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(init storage.InitRoot) *OrderService {
	return &OrderService{InitRepo: init}
}
func (service *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	return service.InitRepo.Order().CreateOrder(in)
}
func (service *OrderService) UpdateOrderStatus(ctx context.Context, in *pb.UpdateOrderStatusRequest) (*pb.OrderStatusResponse, error) {
	return service.InitRepo.Order().UpdateOrderStatus(in)
}
func (service *OrderService) GetOrdersForCustomer(ctx context.Context, in *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	return service.InitRepo.Order().GetOrdersForCustomer(in)
}
func (service *OrderService) GetOrdersForChef(ctx context.Context, in *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	return service.InitRepo.Order().GetOrdersForChef(in)
}
func (service *OrderService) GetOrderById(ctx context.Context, in *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	return service.InitRepo.Order().GetOrderById(in)
}

func (service *OrderService) ActivityUser(ctx context.Context, in *pb.GetUserActivityRequest) (*pb.GetUserActivityResponse, error) {
	return service.InitRepo.Order().GetUserActivity(in)
}
func (service *OrderService) KitchenStatistic(ctx context.Context, in *pb.GetKitchenStatisticRequest) (*pb.KitchenStatisticsResponse, error) {
	return service.InitRepo.Order().GetKitchenStatistics(in)
}
