package app

import (
	"appointment-service/internal/client"
	repository "appointment-service/internal/repository"
	transportgrpc "appointment-service/internal/transport/grpc"
	"appointment-service/internal/usecase"

	appointmentpb "github.com/syndaly1/ap2-assignment2/appointment-service/proto"
	doctorpb "github.com/syndaly1/ap2-assignment2/doctor-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCServer() (*grpc.Server, *grpc.ClientConn, error) {
	repo := repository.NewInMemoryAppointmentRepository()

	doctorConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	doctorGRPCClient := doctorpb.NewDoctorServiceClient(doctorConn)
	doctorClient := client.NewDoctorClient(doctorGRPCClient)

	uc := usecase.NewAppointmentUsecase(repo, doctorClient)
	server := transportgrpc.NewAppointmentServer(uc)

	grpcServer := grpc.NewServer()
	appointmentpb.RegisterAppointmentServiceServer(grpcServer, server)

	return grpcServer, doctorConn, nil
}
