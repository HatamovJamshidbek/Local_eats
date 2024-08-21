package storage

import pb "order_serive/genproto"

type InitRoot interface {
	Review() Reviews
	Meal() Meals
	Order() Orders
	Payment() Payments
}
type Reviews interface {
	CreateReview(request *pb.CreateReviewRequest) (*pb.ReviewResponse, error)
	GetReviews(request *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error)
}
type Meals interface {
	CreateMeal(request *pb.CreateMealRequest) (*pb.MealResponse, error)
	UpdateMeal(request *pb.UpdateMealRequest) (*pb.MealResponse, error)
	DeleteMeal(request *pb.IdRequest) (*pb.Void, error)
	GetAllMeal(request *pb.GetAllMealRequest) (*pb.MealsResponse, error)
	UpdateNutritionInfo(request *pb.UpdateNutritionInfoRequest) (*pb.Dish, error)
}
type Orders interface {
	CreateOrder(request *pb.CreateOrderRequest) (*pb.OrderResponse, error)
	UpdateOrderStatus(request *pb.UpdateOrderStatusRequest) (*pb.OrderStatusResponse, error)
	GetOrdersForCustomer(request *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error)
	GetOrdersForChef(request *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error)
	GetOrderById(request *pb.GetOrderRequest) (*pb.OrderResponse, error)
	GetUserActivity(request *pb.GetUserActivityRequest) (*pb.GetUserActivityResponse, error)
	GetKitchenStatistics(request *pb.GetKitchenStatisticRequest) (*pb.KitchenStatisticsResponse, error)
}
type Payments interface {
	CreatePayment(request *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error)
}
