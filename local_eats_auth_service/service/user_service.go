package service

import (
	pb "auth_serice/genproto"
	"auth_serice/storage"
	"auth_serice/token"
	"context"
)

type UserService struct {
	InitRepo storage.IStorage
	*pb.UnimplementedUserServiceServer
}

func (service *UserService) Users() storage.IUserStorage {
	return service.InitRepo.Users()
}

func (service *UserService) Kitchens() storage.IKitchenStorage {
	//TODO implement me
	panic("implement me")
}

func (service *UserService) mustEmbedUnimplementedKitchenServiceServer() {
	panic("implement me")
}

func NewUserService(init storage.IStorage) *UserService {
	return &UserService{InitRepo: init}
}
func (service *UserService) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return service.InitRepo.Users().Register(in)
}
func (service *UserService) UpdateToken(ctx context.Context, in *pb.UpdateTokenRequest) (*pb.LoginResponse, error) {
	return token.UpdateToken(in)
}
func (service *UserService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	login, err := service.InitRepo.Users().Login(in)
	if err != nil {
		return nil, err
	}
	return token.GenerateJwtToken(login)
}
func (service *UserService) UserProfile(ctx context.Context, in *pb.IdRequest) (*pb.UserResponse, error) {
	return service.InitRepo.Users().GetProfile(in)
}
func (service *UserService) UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UserResponse, error) {
	return service.InitRepo.Users().EditProfile(in)
}
