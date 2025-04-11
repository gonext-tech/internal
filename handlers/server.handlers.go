package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/server_views"
	"github.com/labstack/echo/v4"
)

type ServerService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.MonitoredServer, models.Meta, error)
	GetID(id string) (models.MonitoredServer, error)
	Create(models.MonitoredServer) (models.MonitoredServer, error)
	Update(models.MonitoredServer) (models.MonitoredServer, error)
	Delete(models.MonitoredServer) error
}

type ServerHandler struct {
	ServerServices ServerService
	UploadServices UploadService
}

func NewServerHandler(ss ServerService, us UploadService) *ServerHandler {
	return &ServerHandler{
		ServerServices: ss,
		UploadServices: us,
	}
}

func (sh *ServerHandler) ListPage(c echo.Context) error {
	log.Println("did we enter here first!!")
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
	servers, meta, err := sh.ServerServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch servers"
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
		"Server List(%d)", meta.TotalCount)
	return renderView(c, server_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		server_views.List(titlePage, servers, meta, params),
	))
}

func (sh *ServerHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	server, err := sh.ServerServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch server with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/server")
	}
	titlePage := fmt.Sprintf(
		"Server | %s", server.Name)
	return renderView(c, server_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		server_views.View(server),
	))
}

func (sh *ServerHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Server | Create"
	return renderView(c, server_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		server_views.Create(),
	))
}

func (sh *ServerHandler) CreateHandler(c echo.Context) error {
	var server models.MonitoredServer
	if err := c.Bind(&server); err != nil {
		return err
	}
	log.Println("serverrr", server)
	_, err := sh.ServerServices.Create(server)
	if err != nil {
		return err
	}
	setFlashmessages(c, "success", "server created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/server")
}

func (sh *ServerHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Server | Update"
	id := c.Param("id")
	server, err := sh.ServerServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("server with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	return renderView(c, server_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		server_views.Update(server),
	))
}

func (sh *ServerHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")

	server, err := sh.ServerServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("server with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}

	if err := c.Bind(&server); err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}

	server, err = sh.ServerServices.Update(server)
	if err != nil {
		errorMsg = fmt.Sprintf("server with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "server updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/server")
}

func (sh *ServerHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	project, err := sh.ServerServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("server with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/server")
	}
	err = sh.ServerServices.Delete(project)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete project with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/server")
	}
	setFlashmessages(c, "success", "Server successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/server")
}
