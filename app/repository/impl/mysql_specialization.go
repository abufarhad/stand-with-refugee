package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/infra/errors"
	"clean/infra/logger"
	"gorm.io/gorm"
)

type specialization struct {
	*gorm.DB
}

// NewSpecializationRepository will create an object that represent the Specialization.Repository implementations
func NewSpecializationRepository(db *gorm.DB) repository.ISpecialization {
	return &specialization{
		DB: db,
	}
}

func (s specialization) Save(req *domain.Specialization) (*domain.Specialization, *errors.RestErr) {
	res := s.DB.Model(&domain.Specialization{}).Create(&req)

	if res.Error != nil {
		logger.Error("error occurred when creating specialization", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return req, nil
}

func (s specialization) Update(id uint, specialization string) *errors.RestErr {

	res := s.DB.Model(&domain.Specialization{}).Where("id = ?", id).Updates(map[string]interface{}{"specialization": specialization})
	if res.Error != nil {
		logger.Error("error occurred when updating specialization", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return nil
}

func (s specialization) GetAll() ([]*domain.Specialization, *errors.RestErr) {
	specializations := []*domain.Specialization{}

	err := s.DB.Model(&domain.Specialization{}).Find(&specializations).Error
	if err != nil {
		logger.Error("error occurred when getting all specialization", err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return specializations, nil
}
