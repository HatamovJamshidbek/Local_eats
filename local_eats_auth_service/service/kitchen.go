package service

import (
	pb "auth_serice/genproto"
	"auth_serice/storage"
	"context"
)

type KitchenService struct {
	InitRepo storage.IStorage
	*pb.UnimplementedKitchenServiceServer
}

func NewKitchenService(init storage.IStorage) *KitchenService {
	return &KitchenService{InitRepo: init}
}

func (service *KitchenService) UpdateKitchen(ctx context.Context, in *pb.UpdateKitchenRequest) (*pb.UpdateKitchenResponse, error) {
	return service.InitRepo.Kitchens().UpdateKitchen(in)
}
func (service *KitchenService) CreateKitchen(ctx context.Context, in *pb.CreateKitchenRequest) (*pb.CreateKitchenResponse, error) {
	return service.InitRepo.Kitchens().CreateKitchen(in)
}
func (service *KitchenService) GetByIdKitchen(ctx context.Context, in *pb.IdRequest) (*pb.KitchenResponse, error) {
	return service.InitRepo.Kitchens().GetByIdKitchen(in)
}
func (service *KitchenService) GetAll(ctx context.Context, in *pb.LimitOffset) (*pb.KitchensResponse, error) {
	return service.InitRepo.Kitchens().GetAllKitchens(in)
}
func (service *KitchenService) SearchKitchen(ctx context.Context, in *pb.SearchKitchenRequest) (*pb.KitchensResponse, error) {
	return service.InitRepo.Kitchens().SearchKitchen(in)
}
func (service *KitchenService) UpdateWorkingHours(ctx context.Context, in *pb.UpdateWorkingHoursRequest) (*pb.UpdateWorkingHoursResponse, error) {
	return service.InitRepo.Kitchens().UpdateWorkingHours(in)
}
func (service *KitchenService) ActivityUser(ctx context.Context, in *pb.GetUserActivityRequest) (*pb.GetUserActivityResponse, error) {
	return nil, nil
}
