package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ramyjaber1/internal/models"
	"github.com/ramyjaber1/internal/views/subscription_views"
)

type SubscriptionService interface {
	GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) (models.PaginationResponse, error)
	GetID(id, name string) (models.Subscription, error)
	Create(models.Subscription) (models.Subscription, error)
	Update(models.Subscription) (models.Subscription, error)
}

type SubscriptionHandler struct {
	SubscriptionServices SubscriptionService
	ProjectServices      ProjectService
}

func NewSubscriptionHandler(cs SubscriptionService, ps ProjectService) *SubscriptionHandler {
	return &SubscriptionHandler{
		SubscriptionServices: cs,
		ProjectServices:      ps,
	}
}

func (sh *SubscriptionHandler) ListPage(c echo.Context) error {
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
	project := c.QueryParam("project")
	if project == "" {
		project = ""
	}

	searchTerm := c.QueryParam("searchTerm")
	response, err := sh.SubscriptionServices.GetALL(limit, page, orderBy, sortBy, project, status, searchTerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		"Subscriptions (%d)", response.Meta.TotalCount)

	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscription_views.List(titlePage, response, params),
	))
}

func (sh *SubscriptionHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	subscription, err := sh.SubscriptionServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch hospital with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/hospital")
	}
	titlePage := fmt.Sprintf(
		"Subscription | %s", subscription.ProjectName)
	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscription_views.View(subscription),
	))
}

func (sh *SubscriptionHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Subscription | Create"
	projects, _, _ := sh.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscription_views.Create(projects),
	))
}

func (sh *SubscriptionHandler) CreateHandler(c echo.Context) error {
	var susbcription models.Subscription
	if err := c.Bind(&susbcription); err != nil {
		setFlashmessages(c, "error", errorMsg)
		return sh.CreatePage(c)
	}
	_, err := sh.SubscriptionServices.Create(susbcription)
	if err != nil {
		setFlashmessages(c, "error", "Can't create hospital")
		return sh.CreatePage(c)
	}
	setFlashmessages(c, "success", "hospital created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/hospital")
}

func (sh *SubscriptionHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Subscription | Update"
	id := c.Param("id")
	projectName := c.Param("name")
	subscription, err := sh.SubscriptionServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("subscription with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	projects, _, _ := sh.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscription_views.Update(subscription, projects),
	))
}

func (sh *SubscriptionHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	subscription, err := sh.SubscriptionServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("hospital with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/hospital")
	}
	subscription.Status = "NOT_ACTIVE"
	_, err = sh.SubscriptionServices.Update(subscription)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete hospital with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/hospital")
	}
	setFlashmessages(c, "success", "Donor successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/hospital")
}
