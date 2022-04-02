package controllers

import (
	"clean/app/domain"
	"clean/app/svc"
	"clean/infra/errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type help struct {
	hSvc svc.IHelp
}

// NewSpecializationController will initialize the controllers
func NewHelpController(grp interface{}, hSvc svc.IHelp) {
	sp := &help{
		hSvc: hSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/help/create", sp.Create)
	g.GET("/v1/help/all", sp.GetAll)

}

func (ctr *help) Create(c echo.Context) error {
	var help *domain.Help

	if err := c.Bind(&help); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	result, saveErr := ctr.hSvc.CreateHelp(help)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusCreated, result)
}

func (ctr *help) GetAll(c echo.Context) error {

	result, saveErr := ctr.hSvc.GetAll()
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusOK, result)
}
