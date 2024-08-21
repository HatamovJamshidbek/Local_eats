package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"order_serive/config"
	"order_serive/config/logger"
	pb "order_serive/genproto"
	"order_serive/service"
	"order_serive/storage/postgres"
)

func main() {
	lg := logger.NewLogger()
	storage, err := postgres.NewPostgresStorage()
	if err != nil {
		lg.Fatal("Error creating PostgreSQL storage: %v")
	}

	cnf := config.Load()
	listen, err := net.Listen("tcp", cnf.HTTPPort)
	if err != nil {
		fmt.Println("e---------", err)
		lg.Fatal("Error starting HTTP listener: %v")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, service.NewOrderService(storage))
	err = grpcServer.Serve(listen)
	if err != nil {
		lg.Fatal("Error serving gRPC: %v")
	}

	fmt.Println("Server started successfully")
}
