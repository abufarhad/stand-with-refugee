package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/svc"
	"clean/infra/errors"
)

type specialization struct {
	spRepo repository.ISpecialization
}

func NewSpecializationService(spRepo repository.ISpecialization) svc.ISpecialization {
	return &specialization{
		spRepo: spRepo,
	}
}

func (s specialization) CreateSp(req *domain.Specialization) (*domain.Specialization, *errors.RestErr) {
	return s.spRepo.Save(req)
}

func (s specialization) UpdateSp(id uint, specialization string) *errors.RestErr {
	return s.spRepo.Update(id, specialization)
}

func (s specialization) GetAllSp() ([]*domain.Specialization, *errors.RestErr) {
	return s.spRepo.GetAll()
}
