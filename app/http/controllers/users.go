package controllers

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/consts"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type users struct {
	uSvc svc.IUsers
}

// NewUsersController will initialize the controllers
func NewUsersController(grp interface{}, uSvc svc.IUsers) {
	uc := &users{
		uSvc: uSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/user/signup", uc.Create)
	g.PATCH("/v1/user", uc.Update)
	g.GET("/v1/user/ranklist", uc.GetUserRankList)

	g.POST("/v1/user/commitments", uc.PostCommitments)
	g.GET("/v1/user/commitments/:d_id", uc.GetCommitments)
	g.DELETE("/v1/user/commitment/delete/:c_id", uc.DeleteCommitments)

	g.POST("/v1/place/create", uc.PlaceCreate)
	g.GET("/v1/place/getall", uc.GetPlaces)
	g.DELETE("/v1/place/delete/:p_id", uc.DeletePlaces)

	g.POST("/v1/password/change", uc.ChangePassword)
	g.POST("/v1/password/forgot", uc.ForgotPassword)
	g.POST("/v1/password/verifyreset", uc.VerifyResetPassword)
	g.POST("/v1/password/reset", uc.ResetPassword)
}

func (ctr *users) Create(c echo.Context) error {
	var user domain.User

	if err := c.Bind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(hashedPass)
	user.RoleId = consts.RoleIDAdmin

	result, saveErr := ctr.uSvc.CreateUser(user)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}
	var resp serializers.UserResp
	respErr := methodsutil.StructToStruct(result, &resp)
	if respErr != nil {
		return respErr
	}

	return c.JSON(http.StatusCreated, resp)
}

func (ctr *users) Update(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		logger.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}

	var user serializers.UserReq
	if err := c.Bind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	updateErr := ctr.uSvc.UpdateUser(uint(loggedInUser.ID), user)
	if updateErr != nil {
		return c.JSON(updateErr.Status, updateErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityUpdateSuccessMsg("user")})
}

func (ctr *users) ChangePassword(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		logger.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}
	body := &serializers.ChangePasswordReq{}
	if err := c.Bind(&body); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}
	if err = body.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}
	if body.OldPassword == body.NewPassword {
		restErr := errors.NewBadRequestError("password can't be same as old one")
		return c.JSON(restErr.Status, restErr)
	}
	if err := ctr.uSvc.ChangePassword(loggedInUser.ID, body); err != nil {
		switch err {
		case errors.ErrInvalidPassword:
			restErr := errors.NewBadRequestError("old password didn't match")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityChangedSuccessMsg("password")})
}

func (ctr *users) ForgotPassword(c echo.Context) error {
	body := &serializers.ForgotPasswordReq{}

	if err := c.Bind(&body); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := body.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := ctr.uSvc.ForgotPassword(body.Email); err != nil && err == errors.ErrSendingEmail {
		restErr := errors.NewInternalServerError("failed to send password reset email")
		return c.JSON(restErr.Status, restErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Password reset link sent to email"})
}

func (ctr *users) VerifyResetPassword(c echo.Context) error {
	req := &serializers.VerifyResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := req.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := ctr.uSvc.VerifyResetPassword(req); err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidPasswordResetToken:
			restErr := errors.NewUnauthorizedError("failed to send reset_token email")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	return c.JSON(http.StatusOK, "reset token verified")
}

func (ctr *users) ResetPassword(c echo.Context) error {
	req := &serializers.ResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := req.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	verifyReq := &serializers.VerifyResetPasswordReq{
		Token: req.Token,
		ID:    req.ID,
	}

	if err := ctr.uSvc.VerifyResetPassword(verifyReq); err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidPasswordResetToken:
			restErr := errors.NewUnauthorizedError("failed to send reset_token email")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	if err := ctr.uSvc.ResetPassword(req); err != nil {
		restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		return c.JSON(restErr.Status, restErr)
	}

	return c.JSON(http.StatusOK, "password reset successful")
}

//Creating  commitments stuff
func (ctr *users) PostCommitments(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		logger.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}

	var commitment domain.Commitments

	if err := c.Bind(&commitment); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	commitment.DoctorId = uint(loggedInUser.ID)
	commitment.Point = 50

	resp, postErr := ctr.uSvc.PostCommitments(commitment)
	if postErr != nil {
		return c.JSON(postErr.Status, postErr)
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctr *users) GetCommitments(c echo.Context) error {
	cId, cErr := strconv.Atoi(c.Param("d_id"))
	if cErr != nil {
		restErr := errors.NewBadRequestError(cErr.Error())
		return c.JSON(restErr.Status, restErr)
	}

	resp, err := ctr.uSvc.GetCommitments(uint(cId))
	if err != nil {
		return c.JSON(err.Status, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctr *users) DeleteCommitments(c echo.Context) error {
	cId, cErr := strconv.Atoi(c.Param("c_id"))
	if cErr != nil {
		restErr := errors.NewBadRequestError(cErr.Error())
		return c.JSON(restErr.Status, restErr)
	}

	err := ctr.uSvc.DeleteCommitments(uint(cId))
	if err != nil {
		return c.JSON(err.Status, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Deleted successfully"})
}

// Creating Help stuff
func (ctr *users) CreateHelp(c echo.Context) error {
	var help domain.Help

	if err := c.Bind(&help); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	result, saveErr := ctr.uSvc.CreateHelp(help)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}
	var resp serializers.UserResp
	respErr := methodsutil.StructToStruct(result, &resp)
	if respErr != nil {
		return respErr
	}

	return c.JSON(http.StatusCreated, resp)
}

func (ctr *users) UpdateHelp(c echo.Context) error {

	var help domain.Help
	hId, cErr := strconv.Atoi(c.Param("h_id"))
	if cErr != nil {
		restErr := errors.NewBadRequestError(cErr.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := c.Bind(&help); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	help.ID = uint(hId)
	updateErr := ctr.uSvc.UpdateHelp(help)
	if updateErr != nil {
		return c.JSON(updateErr.Status, updateErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityUpdateSuccessMsg("user")})
}

//Creating  place stuff
func (ctr *users) PlaceCreate(c echo.Context) error {

	var place domain.Place

	if err := c.Bind(&place); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	resp, postErr := ctr.uSvc.PlaceCreate(place)
	if postErr != nil {
		return c.JSON(postErr.Status, postErr)
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctr *users) GetPlaces(c echo.Context) error {

	resp, err := ctr.uSvc.GetPlaces()
	if err != nil {
		return c.JSON(err.Status, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctr *users) DeletePlaces(c echo.Context) error {
	pId, cErr := strconv.Atoi(c.Param("p_id"))
	if cErr != nil {
		restErr := errors.NewBadRequestError(cErr.Error())
		return c.JSON(restErr.Status, restErr)
	}

	err := ctr.uSvc.DeletePlaces(uint(pId))
	if err != nil {
		return c.JSON(err.Status, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Deleted successfully"})
}

// Userlist point wise
func (ctr *users) GetUserRankList(c echo.Context) error {

	result, err := ctr.uSvc.GetUserRankList()
	if err != nil {
		return c.JSON(err.Status, err)
	}
	var resp []*serializers.UserResp
	respErr := methodsutil.StructToStruct(result, &resp)
	if respErr != nil {
		return respErr
	}

	return c.JSON(http.StatusOK, resp)
}
