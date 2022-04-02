package svc

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type IHelp interface {
	CreateHelp(req *domain.Help) (*domain.Help, *errors.RestErr)
	GetAll() ([]*domain.Help, *errors.RestErr)
}
