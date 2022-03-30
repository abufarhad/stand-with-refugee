package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/config"
	"clean/infra/errors"
	"clean/infra/logger"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type users struct {
	urepo repository.IUsers
	rSvc  svc.ICache
}

func (u *users) DeleteCommitments(cId uint) *errors.RestErr {
	return u.urepo.DeleteCommitments(cId)
}

func (u *users) PostCommitments(commitments domain.Commitments) (*domain.Commitments, *errors.RestErr) {
	resp, err := u.urepo.SaveCommitments(commitments)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *users) GetCommitments(cid uint) ([]*domain.Commitments, *errors.RestErr) {
	resp, err := u.urepo.AllCommitments(cid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewUsersService(urepo repository.IUsers, rSvc svc.ICache) svc.IUsers {
	return &users{
		urepo: urepo,
		rSvc:  rSvc,
	}
}

func (u *users) CreateUser(user domain.User) (*domain.User, *errors.RestErr) {
	resp, saveErr := u.urepo.Save(&user)
	if saveErr != nil {
		return nil, saveErr
	}
	return resp, nil
}

func (u *users) GetUserById(userId uint) (*domain.User, *errors.RestErr) {
	resp, getErr := u.urepo.GetUserByID(userId)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}

func (u *users) GetUserByEmail(userName string) (*domain.User, error) {
	resp, getErr := u.urepo.GetUserByEmail(userName)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}

func (u *users) UpdateUser(userID uint, req serializers.UserReq) *errors.RestErr {
	var user domain.User

	err := methodsutil.StructToStruct(req, &user)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("update user"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	user.ID = userID

	if updateErr := u.urepo.Update(&user); updateErr != nil {
		return updateErr
	}

	if err := u.deleteUserCache(int(userID)); err != nil {
		restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		return restErr
	}
	return nil
}

func (u *users) ChangePassword(id int, data *serializers.ChangePasswordReq) error {
	user, getErr := u.urepo.GetUserByID(uint(id))
	if getErr != nil {
		return errors.NewError(getErr.Message)
	}

	currentPass := []byte(user.Password)
	if err := bcrypt.CompareHashAndPassword(currentPass, []byte(data.OldPassword)); err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("comparing hash and old password"), err)
		return errors.ErrInvalidPassword
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 8)

	updates := map[string]interface{}{
		"password":    hashedPass,
		"first_login": false,
	}

	upErr := u.urepo.UpdatePassword(user.ID, updates)
	if upErr != nil {
		return errors.NewError(upErr.Message)
	}

	return nil
}

func (u *users) ForgotPassword(email string) error {
	user, err := u.urepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	secret := passwordResetSecret(user)

	payload := jwt.MapClaims{}
	payload["email"] = user.Email

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		logger.Error("error occur when getting complete signed token", err)
		return err
	}

	// TODO: Send Mail
	logger.Info(signedToken)
	// fpassReq := &serializers.ForgetPasswordMailReq{
	// 	To:     user.Email,
	// 	UserID: user.ID,
	// 	Token:  signedToken,
	// }

	// if err := u.msvc.SendForgotPasswordEmail(*fpassReq); err != nil {
	// 	return errors.ErrSendingEmail
	// }

	return nil
}

func (u *users) VerifyResetPassword(req *serializers.VerifyResetPasswordReq) error {
	user, getErr := u.urepo.GetUserByID(uint(req.ID))
	if getErr != nil {
		return errors.NewError(getErr.Message)
	}

	secret := passwordResetSecret(user)

	parsedToken, err := methodsutil.ParseJwtToken(req.Token, secret)
	if err != nil {
		logger.Error("error occur when parse jwt token with secret", err)
		return errors.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok && !parsedToken.Valid {
		return errors.ErrInvalidPasswordResetToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.ErrInvalidPasswordResetToken
	}

	parsedEmail := claims["email"].(string)
	if user.Email != parsedEmail {
		return errors.ErrInvalidPasswordResetToken
	}

	return nil
}

func (u *users) ResetPassword(req *serializers.ResetPasswordReq) error {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 8)

	if err := u.urepo.ResetPassword(req.ID, hashedPass); err != nil {
		return err
	}

	return nil
}

func (u *users) deleteUserCache(userID int) error {
	if err := u.rSvc.Del(
		config.Redis().UserPrefix+strconv.Itoa(userID),
		config.Redis().TokenPrefix+strconv.Itoa(userID),
	); err != nil {
		logger.Error("error occur when deleting cached user after user update", err)
		return err
	}

	return nil
}

func passwordResetSecret(user *domain.User) string {
	return user.Password + strconv.Itoa(int(user.CreatedAt.Unix()))
}
