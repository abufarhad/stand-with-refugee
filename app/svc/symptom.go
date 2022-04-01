package svc

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type ISymptom interface {
	CreateSym(req *domain.Symptom) (*domain.Symptom, *errors.RestErr)
	UpdateSym(id uint, Symptom string) *errors.RestErr
	GetAllSym() ([]*domain.Symptom, *errors.RestErr)
}
