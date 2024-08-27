package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/project_views"
	"github.com/gonext-tech/internal/views/shop_views/appointments"
	"github.com/labstack/echo/v4"
)

type AppointmentService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, shopID, status string) ([]models.Appointment, models.Meta, error)
	GetID(id string) (models.Appointment, error)
	Create(models.Appointment) (models.Appointment, error)
	Update(models.Appointment) (models.Appointment, error)
}

type AppointmentHandler struct {
	AppointmentServices AppointmentService
}

func NewAppointmentHandler(as AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		AppointmentServices: as,
	}
}

func (ah *AppointmentHandler) ListPage(c echo.Context) error {
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
	shopID := c.Param("id")

	status := c.QueryParam("status")
	searchTerm := c.QueryParam("searchTerm")
	apts, meta, err := ah.AppointmentServices.GetALL(limit, page, orderBy, sortBy, searchTerm, shopID, status)
	if err != nil {
		log.Println("apppointmnett-errr", err)
		isError = false
		errorMsg = "can't fetch appointments"
	}
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
		"Appointment List(%d)", meta.TotalCount)
	return renderView(c, appointments.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		appointments.List(titlePage, apts, meta, params, shopID),
	))
}

func (ah *AppointmentHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Project | Create"
	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.Create(),
	))
}
