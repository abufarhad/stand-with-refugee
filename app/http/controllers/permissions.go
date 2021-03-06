package controllers

import (
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/consts"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"
	"net/http"
	"regexp"
	"strconv"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type permissions struct {
	psvc svc.IPermissions
}

// NewPermissionsController will initialize the controllers
func NewPermissionsController(grp interface{}, ACL func(string) echo.MiddlewareFunc, psvc svc.IPermissions) {
	rc := &permissions{
		psvc: psvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/permission", rc.CreatePermission, ACL(consts.PermissionPermissionCrud))
	g.PATCH("/v1/permission/:permission_id", rc.UpdatePermission, ACL(consts.PermissionPermissionCrud))
	g.DELETE("/v1/permission/:permission_id", rc.DeletePermission, ACL(consts.PermissionPermissionCrud))
	g.GET("/v1/permission", rc.ListPermission, ACL(consts.PermissionPermissionCrud))
}

func (ctr *permissions) CreatePermission(c echo.Context) error {
	var permissionToCreate *serializers.PermissionReq
	var err error

	if err = c.Bind(&permissionToCreate); err != nil {
		logger.Error(msgutil.EntityBindToStructFailedMsg("permission"), err)
		restErr := errors.NewBadRequestError(errors.ErrCheckParamBodyHeader)
		return c.JSON(restErr.Status, restErr)
	}

	if err = permissionToCreate.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	resp, createErr := ctr.psvc.CreatePermission(permissionToCreate)
	if createErr != nil {
		return c.JSON(createErr.Status, createErr)
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctr *permissions) UpdatePermission(c echo.Context) error {
	var permissionToUpdate serializers.PermissionReq

	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("permission id"), err)
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err = c.Bind(&permissionToUpdate); err != nil {
		logger.Error(msgutil.EntityBindToStructFailedMsg("update permission"), err)
		restErr := errors.NewBadRequestError(errors.ErrCheckParamBodyHeader)
		return c.JSON(restErr.Status, restErr)
	}

	if err = permissionToUpdate.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if upErr := ctr.psvc.UpdatePermission(uint(permissionID), permissionToUpdate); upErr != nil {
		return c.JSON(upErr.Status, upErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityUpdateSuccessMsg("permission")})
}

func (ctr *permissions) DeletePermission(c echo.Context) error {
	id := c.Param("permission_id")

	valErr := v.Validate(id, v.Required, v.Match(regexp.MustCompile("^[0-9]+$")).Error("invalid permission id"))
	if valErr != nil {
		restErr := errors.NewBadRequestError(valErr.Error())
		return c.JSON(restErr.Status, restErr)
	}

	roleID, _ := strconv.Atoi(id)

	if delErr := ctr.psvc.DeletePermission(uint(roleID)); delErr != nil {
		return c.JSON(delErr.Status, delErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityDeleteSuccessMsg("permission")})
}

func (ctr *permissions) ListPermission(c echo.Context) error {
	res, err := ctr.psvc.ListPermission()

	if err != nil {
		return c.JSON(err.Status, err)
	}

	return c.JSON(http.StatusOK, res)
}
