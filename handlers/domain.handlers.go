package handlers

import (
	"fmt"
	"net/http"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"github.com/gonext-tech/internal/views/domain_views"
	"github.com/labstack/echo/v4"
)

type DomainService interface {
	GetALL(queries.InvoiceQueryParams) ([]models.Domain, models.Meta, error)
	GetID(id string) (*models.Domain, error)
	Create(*models.Domain) error
	Update(*models.Domain) error
	Delete(*models.Domain) error
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
	var query queries.InvoiceQueryParams
	if err := c.Bind(&query); err != nil {
		errorMsg = "can't read query params"
		setFlashmessages(c, "error", errorMsg)
	}
	query.SetDefaults()
	domains, meta, err := dh.DomainServices.GetALL(query)
	if err != nil {
		isError = false
		errorMsg = "can't fetch clients"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	if c.Request().Header.Get("X-Partial-Content") == "true" {
		// Return only the table content
		return renderView(c, domain_views.List(
			fmt.Sprintf("Domain List(%d)", meta.TotalCount),
			domains,
			meta,
			query,
		))
	}

	titlePage := fmt.Sprintf(
		"Domain List(%d)", meta.TotalCount)
	return renderView(c, domain_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		domain_views.List(titlePage, domains, meta, query),
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
	servers, _, _ := dh.ServerServices.GetALL(queries.InvoiceQueryParams{
		Status:  "UP",
		OrderBy: "desc",
		SortBy:  "id",
		Page:    1,
		Limit:   50,
	})
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
	err := dh.DomainServices.Create(&domain)
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

	servers, _, _ := dh.ServerServices.GetALL(queries.InvoiceQueryParams{
		Status:  "UP",
		OrderBy: "desc",
		SortBy:  "id",
		Page:    1,
		Limit:   50,
	})
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

	if err := c.Bind(domain); err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return dh.UpdatePage(c)
	}

	err = dh.DomainServices.Update(domain)
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
