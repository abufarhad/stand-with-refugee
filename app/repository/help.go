package repository

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type IHelp interface {
	Save(req *domain.Help) (*domain.Help, *errors.RestErr)
	Getall() ([]*domain.Help, *errors.RestErr)
}
