package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/svc"
	"clean/infra/errors"
)

type symptom struct {
	symRepo repository.ISymptom
}

func NewSymptomService(symRepo repository.ISymptom) svc.ISymptom {
	return &symptom{
		symRepo: symRepo,
	}
}

func (s symptom) CreateSym(req *domain.Symptom) (*domain.Symptom, *errors.RestErr) {
	return s.symRepo.Save(req)
}

func (s symptom) UpdateSym(id uint, symptom string) *errors.RestErr {
	return s.symRepo.Update(id, symptom)
}

func (s symptom) GetAllSym() ([]*domain.Symptom, *errors.RestErr) {
	return s.symRepo.GetAll()
}
