package service

import (
	"context"
	pb "order_serive/genproto"
)

func (service *OrderService) CreateMeal(ctx context.Context, in *pb.CreateMealRequest) (*pb.MealResponse, error) {
	return service.InitRepo.Meal().CreateMeal(in)
}
func (service *OrderService) UpdateMeal(ctx context.Context, in *pb.UpdateMealRequest) (*pb.MealResponse, error) {
	return service.InitRepo.Meal().UpdateMeal(in)
}
func (service *OrderService) Delete(ctx context.Context, in *pb.IdRequest) (*pb.Void, error) {
	return service.InitRepo.Meal().DeleteMeal(in)
}
func (service *OrderService) GetAllMeal(ctx context.Context, in *pb.GetAllMealRequest) (*pb.MealsResponse, error) {
	return service.InitRepo.Meal().GetAllMeal(in)
}

func (service *OrderService) UpdateNutritionInfo(ctx context.Context, in *pb.UpdateNutritionInfoRequest) (*pb.Dish, error) {
	return service.InitRepo.Meal().UpdateNutritionInfo(in)
}
