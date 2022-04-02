package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/svc"
	"clean/infra/errors"
)

type help struct {
	hRepo repository.IHelp
}

func NewHelpService(hRepo repository.IHelp) svc.IHelp {
	return &help{
		hRepo: hRepo,
	}
}

func (h help) CreateHelp(req *domain.Help) (*domain.Help, *errors.RestErr) {
	return h.hRepo.Save(req)
}

func (h help) GetAll() ([]*domain.Help, *errors.RestErr) {
	return h.hRepo.Getall()
}
