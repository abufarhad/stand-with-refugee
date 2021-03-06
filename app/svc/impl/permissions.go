package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"
)

type permissions struct {
	prepo repository.IPermissions
}

func NewPermissionsService(prepo repository.IPermissions) svc.IPermissions {
	return &permissions{
		prepo: prepo,
	}
}

func (p *permissions) CreatePermission(req *serializers.PermissionReq) (*domain.Permission, *errors.RestErr) {
	permission := &domain.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	resp, err := p.prepo.Create(permission)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *permissions) GetPermission(id uint) (*serializers.PermissionResp, *errors.RestErr) {
	var resp serializers.RoleResp
	rule, getErr := p.prepo.Get(id)
	if getErr != nil {
		return nil, getErr
	}

	err := methodsutil.StructToStruct(rule, resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("get permission"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return nil, nil
}

func (p *permissions) UpdatePermission(permissionID uint, req serializers.PermissionReq) *errors.RestErr {
	permission := &domain.Permission{
		ID:          permissionID,
		Name:        req.Name,
		Description: req.Description,
	}

	upErr := p.prepo.Update(permission)
	if upErr != nil {
		return upErr
	}

	return nil
}

func (p *permissions) DeletePermission(id uint) *errors.RestErr {
	err := p.prepo.Remove(id)
	if err != nil {
		return err
	}

	return nil
}

func (p *permissions) ListPermission() ([]*serializers.PermissionResp, *errors.RestErr) {
	var resp []*serializers.PermissionResp

	permissions, lstErr := p.prepo.List()
	if lstErr != nil {
		return nil, lstErr
	}

	err := methodsutil.StructToStruct(permissions, &resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("get all permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return resp, nil
}
