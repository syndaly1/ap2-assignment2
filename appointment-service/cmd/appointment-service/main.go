package main

import (
	"appointment-service/internal/app"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer, doctorConn, err := app.NewGRPCServer()
	if err != nil {
		log.Fatalf("failed to create grpc server: %v", err)
	}
	defer doctorConn.Close()

	log.Println("Appointment gRPC service started at :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
