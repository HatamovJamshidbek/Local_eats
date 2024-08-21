package service

import (
	"context"
	pb "order_serive/genproto"
)

func (service *OrderService) CreateReview(ctx context.Context, in *pb.CreateReviewRequest) (*pb.ReviewResponse, error) {
	return service.InitRepo.Review().CreateReview(in)
}
func (service *OrderService) GetReviews(ctx context.Context, in *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	return service.InitRepo.Review().GetReviews(in)
}
