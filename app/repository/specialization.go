package repository

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type ISpecialization interface {
	Save(req *domain.Specialization) (*domain.Specialization, *errors.RestErr)
	Update(id uint, specialization string) *errors.RestErr
	GetAll() ([]*domain.Specialization, *errors.RestErr)
}
