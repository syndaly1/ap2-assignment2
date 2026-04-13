package app

import (
	repository "github.com/syndaly1/ap2-assignment2/doctor-service/internal/repository"
	transportgrpc "github.com/syndaly1/ap2-assignment2/doctor-service/internal/transport/grpc"
	"github.com/syndaly1/ap2-assignment2/doctor-service/internal/usecase"
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
