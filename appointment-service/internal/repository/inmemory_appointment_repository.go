package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/syndaly1/ap2-assignment2/appointment-service/internal/model"
)

type InMemoryAppointmentRepository struct {
	mu           sync.RWMutex
	appointments map[string]model.Appointment
}

func NewInMemoryAppointmentRepository() *InMemoryAppointmentRepository {
	return &InMemoryAppointmentRepository{
		appointments: make(map[string]model.Appointment),
	}
}

func (r *InMemoryAppointmentRepository) Create(appointment model.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.appointments[appointment.ID] = appointment
	return nil
}

func (r *InMemoryAppointmentRepository) GetByID(id string) (model.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	appointment, ok := r.appointments[id]
	if !ok {
		return model.Appointment{}, errors.New("appointment not found")
	}

	return appointment, nil
}

func (r *InMemoryAppointmentRepository) GetAll() ([]model.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	appointments := make([]model.Appointment, 0, len(r.appointments))
	for _, appointment := range r.appointments {
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *InMemoryAppointmentRepository) UpdateStatus(id string, status model.Status, updatedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	appointment, ok := r.appointments[id]
	if !ok {
		return errors.New("appointment not found")
	}

	appointment.Status = status
	appointment.UpdatedAt = updatedAt
	r.appointments[id] = appointment

	return nil
}
