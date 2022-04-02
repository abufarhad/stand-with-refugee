package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/infra/errors"
	"clean/infra/logger"
	"gorm.io/gorm"
)

type help struct {
	*gorm.DB
}

// NewSymptomRepository will create an object that represent the Symptom.Repository implementations
func NewHelpRepository(db *gorm.DB) repository.IHelp {
	return &help{
		DB: db,
	}
}

func (s help) Save(req *domain.Help) (*domain.Help, *errors.RestErr) {
	res := s.DB.Model(&domain.Help{}).Create(&req)

	if res.Error != nil {
		logger.Error("error occurred when creating help", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return req, nil
}

func (s help) Getall() ([]*domain.Help, *errors.RestErr) {
	helps := []*domain.Help{}

	err := s.DB.Model(&domain.Help{}).Find(&helps).Error
	if err != nil {
		logger.Error("error occurred when getting all helps", err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return helps, nil
}
