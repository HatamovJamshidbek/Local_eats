package service

import (
	"context"
	pb "order_serive/genproto"
)

func (service *OrderService) CreatePayment(ctx context.Context, in *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	return service.InitRepo.Payment().CreatePayment(in)
}
