package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/client_views"
	"github.com/gonext-tech/internal/views/domain_views"
	"github.com/labstack/echo/v4"
)

type ClientService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Client, models.Meta, error)
	GetID(id string) (models.Client, error)
	Create(models.Client) (models.Client, error)
	Update(models.Client) (models.Client, error)
	Delete(models.Client) error
}

type ClientHandler struct {
	ClientServices ClientService
	UploadServices UploadService
}

func NewClientHandler(cs ClientService, us UploadService) *ClientHandler {
	return &ClientHandler{
		ClientServices: cs,
		UploadServices: us,
	}
}

func (ch *ClientHandler) ListPage(c echo.Context) error {
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
	clients, meta, err := ch.ClientServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch clients"
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
		"Client List(%d)", meta.TotalCount)
	return renderView(c, domain_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		client_views.List(titlePage, clients, meta, params),
	))
}

func (ch *ClientHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	client, err := ch.ClientServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch client with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/client")
	}
	titlePage := fmt.Sprintf(
		"Client | %s", client.Name)
	return renderView(c, client_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		client_views.View(client),
	))
}

func (ch *ClientHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Client | Create"

	return renderView(c, domain_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		client_views.Create(),
	))
}

func (ch *ClientHandler) CreateHandler(c echo.Context) error {
	var client models.Client
	if err := c.Bind(&client); err != nil {
		return err
	}
	_, err := ch.ClientServices.Create(client)
	if err != nil {
		return err
	}
	imageURLs := UploadImage(c, ch.UploadServices, "", fmt.Sprintf("client/%d", client.ID))

	if len(imageURLs) > 0 {
		client.Image = imageURLs[0]
		_, err = ch.ClientServices.Update(client)
		if err != nil {
			setFlashmessages(c, "error", "cannot upload client image")
			return ch.CreatePage(c)
		}
	}
	setFlashmessages(c, "success", "client created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/client")
}

func (ch *ClientHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Client | Update"
	id := c.Param("id")
	client, err := ch.ClientServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("client with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	return renderView(c, client_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		client_views.Update(client),
	))
}

func (ch *ClientHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")

	client, err := ch.ClientServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("client with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}

	if err := c.Bind(&client); err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}

	client, err = ch.ClientServices.Update(client)
	if err != nil {
		errorMsg = fmt.Sprintf("cannot update client with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}

	imageURLs := UploadImage(c, ch.UploadServices, "", fmt.Sprintf("client/%d", client.ID))

	if len(imageURLs) > 0 {
		client.Image = imageURLs[0]
		_, err = ch.ClientServices.Update(client)
		if err != nil {
			setFlashmessages(c, "error", "cannot upload client image")
			return ch.CreatePage(c)
		}
	}
	setFlashmessages(c, "success", "domain updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/client")
}

func (ch *ClientHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	client, err := ch.ClientServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("client with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/client")
	}
	err = ch.ClientServices.Delete(client)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete client with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/client")
	}
	setFlashmessages(c, "success", "Client successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/client")
}
