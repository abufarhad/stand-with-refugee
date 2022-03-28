package controllers

import (
	"clean/app/domain"
	"clean/app/svc"
	"clean/infra/errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type specializations struct {
	spSvc svc.ISpecialization
}

// NewSpecializationController will initialize the controllers
func NewSpecializationController(grp interface{}, spSvc svc.ISpecialization) {
	sp := &specializations{
		spSvc: spSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/specialization/create", sp.Create)
	g.PATCH("/v1/specialization/:sp_id", sp.Update)
	g.GET("/v1/specialization/all", sp.GetAll)

}

func (ctr *specializations) Create(c echo.Context) error {
	var specialization *domain.Specialization

	if err := c.Bind(&specialization); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	result, saveErr := ctr.spSvc.CreateSp(specialization)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusCreated, result)
}

func (ctr *specializations) Update(c echo.Context) error {

	spID, err := strconv.Atoi(c.Param("sp_id"))
	if err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	var specializationBody *domain.Specialization
	if err := c.Bind(&specializationBody); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}
	specialization := specializationBody.Specialization
	updateErr := ctr.spSvc.UpdateSp(uint(spID), specialization)
	if updateErr != nil {
		return c.JSON(updateErr.Status, updateErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Updated successfully"})
}

func (ctr *specializations) GetAll(c echo.Context) error {

	resp, err := ctr.spSvc.GetAllSp()
	if err != nil {
		return c.JSON(err.Status, err)
	}
	return c.JSON(http.StatusOK, resp)
}
