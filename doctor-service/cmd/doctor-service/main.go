package main

import (
	"log"
	"net"

	"github.com/syndaly1/ap2-assignment2/doctor-service/internal/app"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := app.NewGRPCServer()

	log.Println("Doctor gRPC service started at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
