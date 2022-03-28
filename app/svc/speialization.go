package svc

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type ISpecialization interface {
	CreateSp(req *domain.Specialization) (*domain.Specialization, *errors.RestErr)
	UpdateSp(id uint, specialization string) *errors.RestErr
	GetAllSp() ([]*domain.Specialization, *errors.RestErr)
}
