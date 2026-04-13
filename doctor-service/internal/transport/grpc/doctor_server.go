package grpc

import (
	"context"

	"github.com/syndaly1/ap2-assignment2/doctor-service/internal/model"
	"github.com/syndaly1/ap2-assignment2/doctor-service/internal/usecase"
	doctorpb "github.com/syndaly1/ap2-assignment2/doctor-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DoctorUsecase interface {
	CreateDoctor(fullName, specialization, email string) (model.Doctor, error)
	GetDoctor(id string) (model.Doctor, error)
	GetAllDoctors() ([]model.Doctor, error)
}

type DoctorServer struct {
	doctorpb.UnimplementedDoctorServiceServer
	uc DoctorUsecase
}

func NewDoctorServer(uc DoctorUsecase) *DoctorServer {
	return &DoctorServer{uc: uc}
}

func (s *DoctorServer) CreateDoctor(ctx context.Context, req *doctorpb.CreateDoctorRequest) (*doctorpb.DoctorResponse, error) {
	doctor, err := s.uc.CreateDoctor(req.GetFullName(), req.GetSpecialization(), req.GetEmail())
	if err != nil {
		switch err {
		case usecase.ErrFullNameRequired, usecase.ErrEmailRequired:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case usecase.ErrEmailTaken:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return toDoctorResponse(doctor), nil
}

func (s *DoctorServer) GetDoctor(ctx context.Context, req *doctorpb.GetDoctorRequest) (*doctorpb.DoctorResponse, error) {
	doctor, err := s.uc.GetDoctor(req.GetId())
	if err != nil {
		if err == usecase.ErrDoctorNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return toDoctorResponse(doctor), nil
}

func (s *DoctorServer) ListDoctors(ctx context.Context, req *doctorpb.ListDoctorsRequest) (*doctorpb.ListDoctorsResponse, error) {
	doctors, err := s.uc.GetAllDoctors()
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &doctorpb.ListDoctorsResponse{
		Doctors: make([]*doctorpb.DoctorResponse, 0, len(doctors)),
	}

	for _, doctor := range doctors {
		resp.Doctors = append(resp.Doctors, toDoctorResponse(doctor))
	}

	return resp, nil
}

func toDoctorResponse(doctor model.Doctor) *doctorpb.DoctorResponse {
	return &doctorpb.DoctorResponse{
		Id:             doctor.ID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		Email:          doctor.Email,
	}
}
