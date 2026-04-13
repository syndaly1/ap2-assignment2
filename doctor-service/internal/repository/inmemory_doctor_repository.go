package repository

import (
	"errors"
	"sync"

	"github.com/syndaly1/ap2-assignment2/doctor-service/internal/model"
)

type InMemoryDoctorRepository struct {
	mu      sync.RWMutex
	doctors map[string]model.Doctor
}

func NewInMemoryDoctorRepository() *InMemoryDoctorRepository {
	return &InMemoryDoctorRepository{
		doctors: make(map[string]model.Doctor),
	}
}

func (r *InMemoryDoctorRepository) Create(doctor model.Doctor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.doctors[doctor.ID] = doctor
	return nil
}

func (r *InMemoryDoctorRepository) GetByID(id string) (model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	doctor, ok := r.doctors[id]
	if !ok {
		return model.Doctor{}, errors.New("doctor not found")
	}

	return doctor, nil
}

func (r *InMemoryDoctorRepository) ExistsByEmail(email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, doctor := range r.doctors {
		if doctor.Email == email {
			return true, nil
		}
	}

	return false, nil
}

func (r *InMemoryDoctorRepository) GetAll() ([]model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	doctors := make([]model.Doctor, 0, len(r.doctors))
	for _, doctor := range r.doctors {
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}
