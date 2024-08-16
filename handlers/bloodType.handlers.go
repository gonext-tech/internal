package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ramyjaber1/internal/models"
	blood_views "github.com/ramyjaber1/internal/views/bloodType_views"
)

type BloodService interface {
	GetALL() ([]models.BloodType, error)
	GetID(id string) (models.BloodType, error)
	Create(models.BloodType) (models.BloodType, error)
	Update(models.BloodType) (models.BloodType, error)
	Delete(models.BloodType) error
}

type BloodHandler struct {
	BloodServices BloodService
}

func NewBloodHandler(bs BloodService) *BloodHandler {
	return &BloodHandler{
		BloodServices: bs,
	}
}

func (bh *BloodHandler) ListPage(c echo.Context) error {
	isError = false
	bloods, err := bh.BloodServices.GetALL()
	if err != nil {
		isError = true
		errorMsg = "can't fetch projects"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}
	titlePage := fmt.Sprintf(
		"Blood List(%d)", len(bloods))
	return renderView(c, blood_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		blood_views.List(titlePage, bloods),
	))
}

func (bh *BloodHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	blood, err := bh.BloodServices.GetID(id)
	if err != nil {
		isError = true
		errorMsg = fmt.Sprintf("can't fetch project with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/blood-type")
	}
	titlePage := fmt.Sprintf(
		"Blood | %s", blood.Type)

	return renderView(c, blood_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		blood_views.View(blood),
	))
}

func (bh *BloodHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Bloodtype | Create"
	return renderView(c, blood_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		blood_views.Create(),
	))
}

func (bh *BloodHandler) CreateHandler(c echo.Context) error {
	isError = false
	var blood models.BloodType
	if err := c.Bind(&blood); err != nil {
		setFlashmessages(c, "error", err.Error())
		return bh.CreatePage(c)
	}
	_, err := bh.BloodServices.Create(blood)
	if err != nil {
		setFlashmessages(c, "error", err.Error())
		return bh.CreatePage(c)
	}
	setFlashmessages(c, "success", "blood type created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/blood-type")
}

func (bh *BloodHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Project | Update"
	id := c.Param("id")
	blood, err := bh.BloodServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("blood with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	return renderView(c, blood_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		blood_views.Update(blood),
	))
}

func (bh *BloodHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	blood, err := bh.BloodServices.GetID(id)
	if err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return bh.UpdatePage(c)
	}
	if err := c.Bind(&blood); err != nil {
		isError = true
		errorMsg = "cannot parse the blood body"
		setFlashmessages(c, "error", errorMsg)
		return bh.UpdatePage(c)
	}
	blood, err = bh.BloodServices.Update(blood)
	if err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return bh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "blood type updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/blood-type")
}

func (bh *BloodHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	blood, err := bh.BloodServices.GetID(id)
	if err != nil {
		isError = true
		errorMsg = fmt.Sprintf("project with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/blood-type")
	}
	err = bh.BloodServices.Delete(blood)
	if err != nil {
		isError = true
		errorMsg = fmt.Sprintf("couldnt delete blood type with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/blood-type")
	}
	setFlashmessages(c, "success", "Blood type successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/blood-type")
}
