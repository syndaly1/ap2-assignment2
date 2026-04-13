package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/syndaly1/ap2-assignment2/doctor-service/internal/model"
)

var (
	ErrFullNameRequired = errors.New("full_name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailTaken       = errors.New("email already exists")
	ErrDoctorNotFound   = errors.New("doctor not found")
)

type DoctorRepository interface {
	Create(doctor model.Doctor) error
	GetByID(id string) (model.Doctor, error)
	ExistsByEmail(email string) (bool, error)
	GetAll() ([]model.Doctor, error)
}

type DoctorUsecase struct {
	repo DoctorRepository
}

func NewDoctorUsecase(repo DoctorRepository) *DoctorUsecase {
	return &DoctorUsecase{repo: repo}
}

func (u *DoctorUsecase) CreateDoctor(fullName, specialization, email string) (model.Doctor, error) {
	if fullName == "" {
		return model.Doctor{}, ErrFullNameRequired
	}
	if email == "" {
		return model.Doctor{}, ErrEmailRequired
	}

	exists, err := u.repo.ExistsByEmail(email)
	if err != nil {
		return model.Doctor{}, err
	}
	if exists {
		return model.Doctor{}, ErrEmailTaken
	}

	doctor := model.Doctor{
		ID:             uuid.NewString(),
		FullName:       fullName,
		Specialization: specialization,
		Email:          email,
	}

	if err := u.repo.Create(doctor); err != nil {
		return model.Doctor{}, err
	}

	return doctor, nil
}

func (u *DoctorUsecase) GetDoctor(id string) (model.Doctor, error) {
	doctor, err := u.repo.GetByID(id)
	if err != nil {
		return model.Doctor{}, ErrDoctorNotFound
	}
	return doctor, nil
}

func (u *DoctorUsecase) GetAllDoctors() ([]model.Doctor, error) {
	return u.repo.GetAll()
}
