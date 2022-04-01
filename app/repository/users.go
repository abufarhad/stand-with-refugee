package repository

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type IUsers interface {
	Save(user *domain.User) (*domain.User, *errors.RestErr)
	Update(user *domain.User) *errors.RestErr
	GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr)
	GetUserByID(userID uint) (*domain.User, *errors.RestErr)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserRankListByPoint() ([]*domain.User, *errors.RestErr)

	UpdatePassword(userID uint, updateValue map[string]interface{}) *errors.RestErr
	GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr)
	SetLastLoginAt(user *domain.User) error
	HasRole(userID, roleID uint) bool
	ResetPassword(userID int, hashedPass []byte) error
	GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr)
	GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr)

	SaveCommitments(commitments domain.Commitments) (*domain.Commitments, *errors.RestErr)
	AllCommitments(cid uint) ([]*domain.Commitments, *errors.RestErr)
	DeleteCommitments(cid uint) *errors.RestErr

	SavePlace(place domain.Place) (*domain.Place, *errors.RestErr)
	AllPlaces() ([]*domain.Place, *errors.RestErr)
	DeletePLace(pid uint) *errors.RestErr

	SaveHelp(help *domain.Help) (*domain.Help, *errors.RestErr)
	UpdateHelp(help *domain.Help) *errors.RestErr
}
