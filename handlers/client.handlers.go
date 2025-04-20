package handlers

import (
	"fmt"
	"net/http"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"github.com/gonext-tech/internal/views/client_views"
	"github.com/gonext-tech/internal/views/components"
	"github.com/gonext-tech/internal/views/domain_views"
	"github.com/labstack/echo/v4"
)

type ClientService interface {
	GetALL(queries.InvoiceQueryParams) ([]models.Client, models.Meta, error)
	GetID(id string) (*models.Client, error)
	Create(*models.Client) error
	Update(*models.Client) error
	Delete(*models.Client) error
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
	var query queries.InvoiceQueryParams
	if err := c.Bind(&query); err != nil {
		errorMsg = "can't read query params"
		setFlashmessages(c, "error", errorMsg)
	}
	query.SetDefaults()
	clients, meta, err := ch.ClientServices.GetALL(query)
	if err != nil {
		isError = false
		errorMsg = "can't fetch clients"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	if c.Request().Header.Get("X-Partial-Content") == "true" {
		// Return only the table content
		return renderView(c, client_views.List(
			fmt.Sprintf("Client List(%d)", meta.TotalCount),
			clients,
			meta,
			query,
		))
	}

	titlePage := fmt.Sprintf(
		"Client List(%d)", meta.TotalCount)
	return renderView(c, client_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		client_views.List(titlePage, clients, meta, query),
	))
}

func (ch *ClientHandler) Search(c echo.Context) error {
	isError = false
	var query queries.InvoiceQueryParams
	if err := c.Bind(&query); err != nil {
		errorMsg = "can't read query params"
		setFlashmessages(c, "error", errorMsg)
	}
	query.Page = 1
	query.Limit = 5
	query.SortBy = "id"
	query.OrderBy = "desc"

	clients, _, err := ch.ClientServices.GetALL(query)
	if err != nil {
		isError = false
		errorMsg = "can't fetch clients"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		renderView(c, components.ClientResult([]models.Client{}))

	}

	return renderView(c, components.ClientResult(clients))
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
	err := ch.ClientServices.Create(&client)
	if err != nil {
		return err
	}
	imageURLs := UploadImage(c, ch.UploadServices, "", fmt.Sprintf("client/%d", client.ID))

	if len(imageURLs) > 0 {
		client.Image = imageURLs[0]
		err = ch.ClientServices.Update(&client)
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

	if err := c.Bind(client); err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}

	err = ch.ClientServices.Update(client)
	if err != nil {
		errorMsg = fmt.Sprintf("cannot update client with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}

	imageURLs := UploadImage(c, ch.UploadServices, "", fmt.Sprintf("client/%d", client.ID))

	if len(imageURLs) > 0 {
		client.Image = imageURLs[0]
		err = ch.ClientServices.Update(client)
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
