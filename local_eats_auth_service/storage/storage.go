package storage

import (
	pb "auth_serice/genproto"
)

type IStorage interface {
	Users() IUserStorage
	Kitchens() IKitchenStorage
}

type IUserStorage interface {
	Register(request *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Login(request *pb.LoginRequest) (*pb.Claims, error)
	GetProfile(request *pb.IdRequest) (*pb.UserResponse, error)
	EditProfile(request *pb.UpdateUserProfileRequest) (*pb.UserResponse, error)
}

type IKitchenStorage interface {
	CreateKitchen(request *pb.CreateKitchenRequest) (*pb.CreateKitchenResponse, error)
	UpdateKitchen(request *pb.UpdateKitchenRequest) (*pb.UpdateKitchenResponse, error)
	GetByIdKitchen(request *pb.IdRequest) (*pb.KitchenResponse, error)
	GetAllKitchens(request *pb.LimitOffset) (*pb.KitchensResponse, error)
	SearchKitchen(request *pb.SearchKitchenRequest) (*pb.KitchensResponse, error)
	UpdateWorkingHours(req *pb.UpdateWorkingHoursRequest) (*pb.UpdateWorkingHoursResponse, error)
}
