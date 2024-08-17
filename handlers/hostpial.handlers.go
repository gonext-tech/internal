package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/hospital_views"
)

type HospitalService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Hospital, models.Meta, error)
	GetID(id string) (models.Hospital, error)
	Create(models.Hospital) (models.Hospital, error)
	Update(models.Hospital) (models.Hospital, error)
	Delete(models.Hospital) error
}

type HospitalHandler struct {
	HospitalServices HospitalService
	CityServices     CityService
}

func NewHospitalHandler(hs HospitalService, cs CityService) *HospitalHandler {
	return &HospitalHandler{
		HospitalServices: hs,
		CityServices:     cs,
	}
}

func (hh *HospitalHandler) ListPage(c echo.Context) error {
	isError = false
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}
	orderBy := c.QueryParam("orderBy")
	if orderBy == "" {
		orderBy = "desc"
	}
	sortBy := c.QueryParam("sortBy")
	if sortBy == "" {
		sortBy = "id"
	}
	status := c.QueryParam("status")
	searchTerm := c.QueryParam("searchTerm")
	hospitals, meta, err := hh.HospitalServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch hospitals"
	}
	log.Println("errr", err)
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}
	var params models.ParamResponse
	if searchTerm != "" {
		params.Search = searchTerm
	}
	if status != "" {
		params.Status = status
	}
	params.Page = page
	params.Limit = limit
	params.SortBy = sortBy
	params.OrderBy = orderBy
	titlePage := fmt.Sprintf(
		"Hospitals (%d)", meta.TotalCount)

	return renderView(c, hospital_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		hospital_views.List(titlePage, hospitals, meta, params),
	))
}

func (hh *HospitalHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	hospital, err := hh.HospitalServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch hospital with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/hospital")
	}
	titlePage := fmt.Sprintf(
		"Hospital | %s", hospital.Name)
	return renderView(c, hospital_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		hospital_views.View(hospital),
	))
}

func (hh *HospitalHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Hospital | Create"
	cities, _ := hh.CityServices.GetALL()
	return renderView(c, hospital_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		hospital_views.Create(cities),
	))
}

func (hh *HospitalHandler) CreateHandler(c echo.Context) error {
	var hospital models.Hospital
	if err := c.Bind(&hospital); err != nil {
		setFlashmessages(c, "error", errorMsg)
		return hh.CreatePage(c)
	}
	_, err := hh.HospitalServices.Create(hospital)
	if err != nil {
		setFlashmessages(c, "error", "Can't create hospital")
		return hh.CreatePage(c)
	}
	setFlashmessages(c, "success", "hospital created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/hospital")
}

func (hh *HospitalHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Hospital | Update"
	id := c.Param("id")
	hospital, err := hh.HospitalServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("hospital with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	cities, _ := hh.CityServices.GetALL()
	return renderView(c, hospital_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		hospital_views.Update(hospital, cities),
	))
}

func (hh *HospitalHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	hospital, err := hh.HospitalServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("hospital with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return hh.UpdatePage(c)
	}

	if err := c.Bind(&hospital); err != nil {
		errorMsg = "cannot parse the hospital body"
		setFlashmessages(c, "error", errorMsg)
		return hh.UpdatePage(c)
	}

	_, err = hh.HospitalServices.Update(hospital)
	if err != nil {
		errorMsg = fmt.Sprintf("hospital with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return hh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "hospital updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/hospital")
}

func (hh *HospitalHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	hospital, err := hh.HospitalServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("hospital with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/hospital")
	}
	err = hh.HospitalServices.Delete(hospital)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete hospital with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/hospital")
	}
	setFlashmessages(c, "success", "Donor successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/hospital")
}
