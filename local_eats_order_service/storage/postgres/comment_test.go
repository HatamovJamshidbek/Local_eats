package postgres

import (
	pb "order_serive/genproto"
	"reflect"
	"testing"
)

func TestReviewRepository_CreateReview(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Review()

	request := &pb.CreateReviewRequest{
		OrderId:   "order123",
		UserId:    "user456",
		KitchenId: "kitchen789",
		Rating:    4.5,
		Comment:   "Great service!",
	}

	response, err := repo.CreateReview(request)
	if err != nil {
		t.Fatalf("Error creating review: %v", err)
	}

	expectedResponse := &pb.ReviewResponse{
		Id:        response.Id,
		OrderId:   request.OrderId,
		UserId:    request.UserId,
		KitchenId: request.KitchenId,
		Rating:    request.Rating,
		Comment:   request.Comment,
		CreatedAt: response.CreatedAt,
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("CreateReview response does not match expected.\nGot: %+v\nExpected: %+v", response, expectedResponse)
	}
}

func TestReviewRepository_GetReviews(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Review()
	request := &pb.GetReviewsRequest{
		KitchenId: "kitchen789",
		LimitOffset: &pb.LimitOffset{
			Limit:  10,
			Offset: 0,
		},
	}

	response, err := repo.GetReviews(request)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}

	expectedResponse := &pb.GetReviewsResponse{
		Reviews:       response.Reviews,
		Total:         response.Total,
		AverageRating: response.AverageRating,
		LimitOffset:   response.LimitOffset,
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("get response does not match expected.\nGot: %+v\nExpected: %+v", response, expectedResponse)
	}
}
