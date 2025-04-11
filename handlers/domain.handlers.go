package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/domain_views"
	"github.com/labstack/echo/v4"
)

type DomainService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Domain, models.Meta, error)
	GetID(id string) (models.Domain, error)
	Create(models.Domain) (models.Domain, error)
	Update(models.Domain) (models.Domain, error)
	Delete(models.Domain) error
}

type DomainHandler struct {
	DomainServices DomainService
	ServerServices ServerService
}

func NewDomainHandler(ds DomainService, ss ServerService) *DomainHandler {
	return &DomainHandler{
		DomainServices: ds,
		ServerServices: ss,
	}
}

func (dh *DomainHandler) ListPage(c echo.Context) error {
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
	domains, meta, err := dh.DomainServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
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
		"Domain List(%d)", meta.TotalCount)
	return renderView(c, domain_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		domain_views.List(titlePage, domains, meta, params),
	))
}

func (dh *DomainHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	domain, err := dh.DomainServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch domain with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/domain")
	}
	titlePage := fmt.Sprintf(
		"Domain | %s", domain.Name)
	return renderView(c, domain_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		domain_views.View(domain),
	))
}

func (dh *DomainHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Domain | Create"
	servers, _, _ := dh.ServerServices.GetALL(50, 1, "desc", "id", "", "UP")
	return renderView(c, domain_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		domain_views.Create(servers),
	))
}

func (dh *DomainHandler) CreateHandler(c echo.Context) error {
	var domain models.Domain
	if err := c.Bind(&domain); err != nil {
		return err
	}
	_, err := dh.DomainServices.Create(domain)
	if err != nil {
		return err
	}
	setFlashmessages(c, "success", "domain created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/domain")
}

func (dh *DomainHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Domain | Update"
	id := c.Param("id")
	domain, err := dh.DomainServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("domain with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	servers, _, _ := dh.ServerServices.GetALL(50, 1, "desc", "id", "", "UP")
	return renderView(c, domain_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		domain_views.Update(domain, servers),
	))
}

func (dh *DomainHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")

	domain, err := dh.DomainServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("domain with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return dh.UpdatePage(c)
	}

	if err := c.Bind(&domain); err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return dh.UpdatePage(c)
	}

	domain, err = dh.DomainServices.Update(domain)
	if err != nil {
		errorMsg = fmt.Sprintf("domain with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return dh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "domain updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/domain")
}

func (dh *DomainHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	domain, err := dh.DomainServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("domain with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/domain")
	}
	err = dh.DomainServices.Delete(domain)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete domain with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/domain")
	}
	setFlashmessages(c, "success", "Domain successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/domain")
}
