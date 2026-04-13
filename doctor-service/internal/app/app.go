package app

import (
	repository "doctor-service/internal/repository"
	transportgrpc "doctor-service/internal/transport/grpc"
	"doctor-service/internal/usecase"

	doctorpb "github.com/syndaly1/ap2-assignment2/doctor-service/proto"

	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	repo := repository.NewInMemoryDoctorRepository()
	uc := usecase.NewDoctorUsecase(repo)
	server := transportgrpc.NewDoctorServer(uc)

	grpcServer := grpc.NewServer()
	doctorpb.RegisterDoctorServiceServer(grpcServer, server)

	return grpcServer
}
