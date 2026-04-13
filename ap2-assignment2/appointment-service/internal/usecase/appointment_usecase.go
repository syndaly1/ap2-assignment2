package usecase

import (
	"appointment-service/internal/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTitleRequired       = errors.New("title is required")
	ErrDoctorIDRequired    = errors.New("doctor_id is required")
	ErrDoctorNotFound      = errors.New("doctor not found")
	ErrDoctorUnavailable   = errors.New("doctor service unavailable")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidTransition   = errors.New("cannot transition from done to new")
	ErrAppointmentNotFound = errors.New("appointment not found")
)

type AppointmentRepository interface {
	Create(appointment model.Appointment) error
	GetByID(id string) (model.Appointment, error)
	GetAll() ([]model.Appointment, error)
	UpdateStatus(id string, status model.Status, updatedAt time.Time) error
}

type DoctorClient interface {
	GetDoctor(ctx context.Context, id string) error
}

type AppointmentUsecase struct {
	repo         AppointmentRepository
	doctorClient DoctorClient
}

func NewAppointmentUsecase(repo AppointmentRepository, doctorClient DoctorClient) *AppointmentUsecase {
	return &AppointmentUsecase{
		repo:         repo,
		doctorClient: doctorClient,
	}
}

func (u *AppointmentUsecase) CreateAppointment(ctx context.Context, title, description, doctorID string) (model.Appointment, error) {
	if title == "" {
		return model.Appointment{}, ErrTitleRequired
	}
	if doctorID == "" {
		return model.Appointment{}, ErrDoctorIDRequired
	}

	if err := u.doctorClient.GetDoctor(ctx, doctorID); err != nil {
		return model.Appointment{}, err
	}

	now := time.Now()
	appointment := model.Appointment{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		DoctorID:    doctorID,
		Status:      model.StatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.repo.Create(appointment); err != nil {
		return model.Appointment{}, err
	}

	return appointment, nil
}

func (u *AppointmentUsecase) GetAppointment(id string) (model.Appointment, error) {
	appointment, err := u.repo.GetByID(id)
	if err != nil {
		return model.Appointment{}, ErrAppointmentNotFound
	}

	return appointment, nil
}

func (u *AppointmentUsecase) GetAllAppointments() ([]model.Appointment, error) {
	return u.repo.GetAll()
}

func (u *AppointmentUsecase) UpdateStatus(ctx context.Context, id string, newStatus model.Status) (model.Appointment, error) {
	if newStatus != model.StatusNew && newStatus != model.StatusInProgress && newStatus != model.StatusDone {
		return model.Appointment{}, ErrInvalidStatus
	}

	appointment, err := u.repo.GetByID(id)
	if err != nil {
		return model.Appointment{}, ErrAppointmentNotFound
	}

	if err := u.doctorClient.GetDoctor(ctx, appointment.DoctorID); err != nil {
		return model.Appointment{}, err
	}

	if appointment.Status == model.StatusDone && newStatus == model.StatusNew {
		return model.Appointment{}, ErrInvalidTransition
	}

	now := time.Now()
	if err := u.repo.UpdateStatus(id, newStatus, now); err != nil {
		return model.Appointment{}, err
	}

	appointment.Status = newStatus
	appointment.UpdatedAt = now

	return appointment, nil
}
