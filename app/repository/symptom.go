package repository

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type ISymptom interface {
	Save(req *domain.Symptom) (*domain.Symptom, *errors.RestErr)
	Update(id uint, symptom string) *errors.RestErr
	GetAll() ([]*domain.Symptom, *errors.RestErr)
}
