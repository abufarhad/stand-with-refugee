package controllers

import (
	"clean/app/domain"
	"clean/app/svc"
	"clean/infra/errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type symptoms struct {
	symSvc svc.ISymptom
}

// NewSpecializationController will initialize the controllers
func NewSymptomController(grp interface{}, symSvc svc.ISymptom) {
	sp := &symptoms{
		symSvc: symSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/symptom/create", sp.Create)
	g.PATCH("/v1/symptom/:sym_id", sp.Update)
	g.GET("/v1/symptom/all", sp.GetAll)

}

func (ctr *symptoms) Create(c echo.Context) error {
	var symptom *domain.Symptom

	if err := c.Bind(&symptom); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	result, saveErr := ctr.symSvc.CreateSym(symptom)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusCreated, result)
}

func (ctr *symptoms) Update(c echo.Context) error {

	spID, err := strconv.Atoi(c.Param("sym_id"))
	if err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	var symptom *domain.Symptom
	if err := c.Bind(&symptom); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}
	symp := symptom.Symptom
	updateErr := ctr.symSvc.UpdateSym(uint(spID), symp)
	if updateErr != nil {
		return c.JSON(updateErr.Status, updateErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Updated successfully"})
}

func (ctr *symptoms) GetAll(c echo.Context) error {

	resp, err := ctr.symSvc.GetAllSym()
	if err != nil {
		return c.JSON(err.Status, err)
	}
	return c.JSON(http.StatusOK, resp)
}
