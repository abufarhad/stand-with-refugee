package svc

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infra/errors"
)

type IUsers interface {
	CreateUser(domain.User) (*domain.User, *errors.RestErr)
	GetUserById(uid uint) (*domain.User, *errors.RestErr)

	PostCommitments(commitments domain.Commitments) (*domain.Commitments, *errors.RestErr)
	GetCommitments(cid uint) ([]*domain.Commitments, *errors.RestErr)
	DeleteCommitments(cid uint) *errors.RestErr

	GetUserByEmail(useremail string) (*domain.User, error)
	UpdateUser(userID uint, req serializers.UserReq) *errors.RestErr
	ChangePassword(id int, data *serializers.ChangePasswordReq) error
	ForgotPassword(email string) error
	VerifyResetPassword(req *serializers.VerifyResetPasswordReq) error
	ResetPassword(req *serializers.ResetPasswordReq) error
}
