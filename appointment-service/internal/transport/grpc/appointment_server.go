package grpc

import (
	"context"
	"time"

	"github.com/syndaly1/ap2-assignment2/appointment-service/internal/model"
	"github.com/syndaly1/ap2-assignment2/appointment-service/internal/usecase"
	appointmentpb "github.com/syndaly1/ap2-assignment2/appointment-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppointmentUsecase interface {
	CreateAppointment(ctx context.Context, title, description, doctorID string) (model.Appointment, error)
	GetAppointment(id string) (model.Appointment, error)
	GetAllAppointments() ([]model.Appointment, error)
	UpdateStatus(ctx context.Context, id string, newStatus model.Status) (model.Appointment, error)
}

type AppointmentServer struct {
	appointmentpb.UnimplementedAppointmentServiceServer
	uc AppointmentUsecase
}

func NewAppointmentServer(uc AppointmentUsecase) *AppointmentServer {
	return &AppointmentServer{uc: uc}
}

func (s *AppointmentServer) CreateAppointment(ctx context.Context, req *appointmentpb.CreateAppointmentRequest) (*appointmentpb.AppointmentResponse, error) {
	appointment, err := s.uc.CreateAppointment(ctx, req.GetTitle(), req.GetDescription(), req.GetDoctorId())
	if err != nil {
		switch err {
		case usecase.ErrTitleRequired, usecase.ErrDoctorIDRequired:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case usecase.ErrDoctorNotFound:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case usecase.ErrDoctorUnavailable:
			return nil, status.Error(codes.Unavailable, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return toAppointmentResponse(appointment), nil
}

func (s *AppointmentServer) GetAppointment(ctx context.Context, req *appointmentpb.GetAppointmentRequest) (*appointmentpb.AppointmentResponse, error) {
	appointment, err := s.uc.GetAppointment(req.GetId())
	if err != nil {
		if err == usecase.ErrAppointmentNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return toAppointmentResponse(appointment), nil
}

func (s *AppointmentServer) ListAppointments(ctx context.Context, req *appointmentpb.ListAppointmentsRequest) (*appointmentpb.ListAppointmentsResponse, error) {
	appointments, err := s.uc.GetAllAppointments()
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &appointmentpb.ListAppointmentsResponse{
		Appointments: make([]*appointmentpb.AppointmentResponse, 0, len(appointments)),
	}

	for _, appointment := range appointments {
		resp.Appointments = append(resp.Appointments, toAppointmentResponse(appointment))
	}

	return resp, nil
}

func (s *AppointmentServer) UpdateAppointmentStatus(ctx context.Context, req *appointmentpb.UpdateStatusRequest) (*appointmentpb.AppointmentResponse, error) {
	appointment, err := s.uc.UpdateStatus(ctx, req.GetId(), model.Status(req.GetStatus()))
	if err != nil {
		switch err {
		case usecase.ErrInvalidStatus:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case usecase.ErrInvalidTransition:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case usecase.ErrAppointmentNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case usecase.ErrDoctorNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case usecase.ErrDoctorUnavailable:
			return nil, status.Error(codes.Unavailable, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return toAppointmentResponse(appointment), nil
}

func toAppointmentResponse(appointment model.Appointment) *appointmentpb.AppointmentResponse {
	return &appointmentpb.AppointmentResponse{
		Id:          appointment.ID,
		Title:       appointment.Title,
		Description: appointment.Description,
		DoctorId:    appointment.DoctorID,
		Status:      string(appointment.Status),
		CreatedAt:   appointment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   appointment.UpdatedAt.Format(time.RFC3339),
	}
}
