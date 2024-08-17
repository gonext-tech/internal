package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/ticket_views"
)

type TicketService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Ticket, models.Meta, error)
	GetID(id string) (models.Ticket, error)
	Create(models.Ticket) (models.Ticket, error)
	Update(models.Ticket) (models.Ticket, error)
	Delete(models.Ticket) error
}

type TicketHandler struct {
	TicketServices       TicketService
	ProjectServices      ProjectService
	NotificationServices NotificationService
}

func NewTicketHandler(ts TicketService, ps ProjectService, ns NotificationService) *TicketHandler {
	return &TicketHandler{
		TicketServices:       ts,
		ProjectServices:      ps,
		NotificationServices: ns,
	}
}

func (th *TicketHandler) ListPage(c echo.Context) error {
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
	tickets, meta, err := th.TicketServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		errorMsg = "can't fetch projects"
		setFlashmessages(c, "error", errorMsg)
	}

	titlePage := fmt.Sprintf(
		"Ticket List(%d)", meta.TotalCount)
	return renderView(c, ticket_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		ticket_views.List(titlePage, tickets, meta),
	))
}

func (th *TicketHandler) Search(c echo.Context) error {
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
	tickets, _, err := th.TicketServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch projects"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}
	log.Println("ticketss", tickets)

	return renderView(c, ticket_views.TableRows(tickets))
}

func (th *TicketHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	ticket, err := th.TicketServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch ticket with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/ticket")
	}
	titlePage := fmt.Sprintf(
		"Ticket | %d", ticket.ID)
	return renderView(c, ticket_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		ticket_views.View(ticket),
	))
}

func (th *TicketHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Ticket | Create"
	return renderView(c, ticket_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		ticket_views.Create(),
	))
}

func (th *TicketHandler) ApiCreate(c echo.Context) error {
	var ticket models.Ticket
	if err := c.Bind(&ticket); err != nil {
		return c.JSON(404, map[string]interface{}{
			"ok":      false,
			"message": "failed to save",
		})
	}
	ticket, err := th.TicketServices.Create(ticket)
	if err != nil {
		return c.JSON(404, map[string]interface{}{
			"ok":      false,
			"message": "failed to save",
		})
	}
	if err := th.NotificationServices.Create(ticket, ""); err != nil {
		return c.JSON(404, map[string]interface{}{
			"ok":      false,
			"message": "failed to send notification",
		})
	}
	return c.JSON(200, map[string]interface{}{
		"ok":      true,
		"message": "ticket send successfully",
	})
}

func (th *TicketHandler) CreateHandler(c echo.Context) error {
	var ticket models.Ticket
	if err := c.Bind(&ticket); err != nil {
		return err
	}
	_, err := th.TicketServices.Create(ticket)
	if err != nil {
		return err
	}
	err = th.NotificationServices.Create(ticket, "")
	if err != nil {
		setFlashmessages(c, "error", "couldn't create notification")
		th.CreatePage(c)
	}
	setFlashmessages(c, "success", "Ticket created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/ticket")
}

func (th *TicketHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Ticket | Update"
	id := c.Param("id")
	ticket, err := th.TicketServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("project with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	projects, _, _ := th.ProjectServices.GetALL(100, 1, "desc", "id", "", "")
	return renderView(c, ticket_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		ticket_views.Update(ticket, projects),
	))
}

func (th *TicketHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	ticket, err := th.TicketServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("ticket with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return th.UpdatePage(c)
	}
	if err := c.Bind(&ticket); err != nil {
		errorMsg = "cannot parse the ticket body"
		setFlashmessages(c, "error", errorMsg)
		return th.UpdatePage(c)
	}
	_, err = th.TicketServices.Update(ticket)
	if err != nil {
		errorMsg = fmt.Sprintf("ticket with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return th.UpdatePage(c)
	}
	setFlashmessages(c, "success", "Ticket updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/ticket")
}

func (th *TicketHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	project, err := th.TicketServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("project with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/ticket")
	}
	err = th.TicketServices.Delete(project)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete ticket with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/ticket")
	}
	setFlashmessages(c, "success", "Ticket successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/ticket")
}
