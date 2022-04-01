package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/infra/errors"
	"clean/infra/logger"
	"gorm.io/gorm"
)

type symptom struct {
	*gorm.DB
}

// NewSymptomRepository will create an object that represent the Symptom.Repository implementations
func NewSymptomRepository(db *gorm.DB) repository.ISymptom {
	return &symptom{
		DB: db,
	}
}

func (s symptom) Save(req *domain.Symptom) (*domain.Symptom, *errors.RestErr) {
	res := s.DB.Model(&domain.Symptom{}).Create(&req)

	if res.Error != nil {
		logger.Error("error occurred when creating specialization", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return req, nil
}

func (s symptom) Update(id uint, symptom string) *errors.RestErr {
	res := s.DB.Model(&domain.Symptom{}).Where("id = ?", id).Updates(map[string]interface{}{"symptom": symptom})
	if res.Error != nil {

		logger.Error("error occurred when updating specialization", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return nil
}

func (s symptom) GetAll() ([]*domain.Symptom, *errors.RestErr) {
	symptoms := []*domain.Symptom{}

	err := s.DB.Model(&domain.Symptom{}).Find(&symptoms).Error
	if err != nil {
		logger.Error("error occurred when getting all specialization", err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return symptoms, nil
}
