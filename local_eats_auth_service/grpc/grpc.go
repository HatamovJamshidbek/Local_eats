package grpc

import (
	pb "auth_serice/genproto"
	"auth_serice/service"
	"auth_serice/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(storage storage.IStorage) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, service.NewUserService(storage))
	pb.RegisterKitchenServiceServer(grpcServer, service.NewKitchenService(storage))

	reflection.Register(grpcServer)
	return grpcServer
}
