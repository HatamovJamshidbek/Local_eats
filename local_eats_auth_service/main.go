package main

import (
	"auth_serice/api"
	"auth_serice/api/handlers"
	"auth_serice/config"
	grpcServer "auth_serice/grpc"
	"auth_serice/service"
	"auth_serice/storage/postgres"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"sync"
)

func main() {
	cnf := config.Load()

	store, err := postgres.NewPostgresStorage()
	if err != nil {
		log.Fatalf("Failed to set up storage: %v", err)
	}

	userService := service.NewUserService(store)
	handler := handlers.NewHandler(userService)

	r := gin.Default()
	api.RouterApi(handler)

	grpcSrv := grpcServer.SetUpServer(store)
	lis, err := net.Listen("tcp", cnf.AuthServiceGrpcHost+":"+cnf.AuthServiceGrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", cnf.AuthServiceGrpcPort, err)
	}

	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	log.Println("Starting HTTP server on port", cnf.HTTPPort)
	if err := r.Run(":" + cnf.HTTPPort); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
