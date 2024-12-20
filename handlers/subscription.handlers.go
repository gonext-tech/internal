package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/shop_views/subscriptions"
	"github.com/gonext-tech/internal/views/subscription_views"
	"github.com/labstack/echo/v4"
)

type SubscriptionService interface {
	GetALL(limit, page int, orderBy, sortBy, project, shop, status, searchTerm string) ([]models.Subscription, models.Meta, error)
	GetID(id, name string) (models.Subscription, error)
	Create(models.Subscription) (models.Subscription, error)
	Update(models.Subscription) (models.Subscription, error)
	Delete(models.Subscription) (models.Subscription, error)
}

type SubscriptionHandler struct {
	SubscriptionServices SubscriptionService
	ShopServices         ShopService
	ProjectServices      ProjectService
	MembershipServices   MembershipService
	StatsServices        StatsService
}

func NewSubscriptionHandler(ms SubscriptionService, ps ProjectService, mems MembershipService, ss ShopService, sss StatsService) *SubscriptionHandler {
	return &SubscriptionHandler{
		SubscriptionServices: ms,
		ProjectServices:      ps,
		ShopServices:         ss,
		MembershipServices:   mems,
		StatsServices:        sss,
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
	shop := c.QueryParam("shop")
	if project == "" {
		project = ""
	}

	searchTerm := c.QueryParam("searchTerm")
	response, meta, err := sh.SubscriptionServices.GetALL(limit, page, orderBy, sortBy, project, shop, status, searchTerm)
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
		"Subscriptions (%d)", meta.TotalCount)
	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscription_views.List(titlePage, response, meta, params),
	))
}

func (sh *SubscriptionHandler) ShopListPage(c echo.Context) error {
	isError = false
	shop := c.Param("id")
	project := c.Param("name")

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
	response, meta, err := sh.SubscriptionServices.GetALL(limit, page, orderBy, sortBy, project, shop, status, searchTerm)
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
		"Memberships (%d)", meta.TotalCount)
	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscriptions.List(titlePage, response, meta, params),
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
		return c.Redirect(http.StatusSeeOther, "/subscription")
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
	projects, _, _ := sh.ProjectServices.GetALL(50, 1, "desc", "id", "", "ACTIVE")
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
	var subscription models.Subscription
	if err := c.Bind(&subscription); err != nil {
		setFlashmessages(c, "error", err.Error())
		return sh.CreatePage(c)
	}
	if subscription.PaymentStatus == "PAID" {
		now := time.Now()
		subscription.PaidAt = &now
	}
	_, err := sh.SubscriptionServices.Create(subscription)
	if err != nil {
		setFlashmessages(c, "error", "Can't create subscription")
		return sh.CreatePage(c)
	}
	//NOTE: --> Handle statistic stuff <--
	err = sh.StatsServices.HandleStatsSubscription(models.Subscription{}, subscription)
	if err != nil {
		setFlashmessages(c, "error", "Can't update subscriptions stats,contact the support")
		return sh.CreatePage(c)
	}
	// --> END <--

	//NOTE: --> Get shop and update the nextBillingDate <--
	shop, err := sh.ShopServices.GetID(fmt.Sprintf("%d", subscription.ShopID), subscription.ProjectName)
	if err != nil {
		setFlashmessages(c, "error", "Can't get shop")
		return sh.CreatePage(c)
	}
	if shop.NextBillingDate == nil || shop.NextBillingDate.Before(subscription.EndDate) {
		shop.NextBillingDate = &subscription.EndDate
		//NOTE: --> we should do this so we don't get error when updating the shop <--
		shop.Owner = models.User{}
		shop.Workers = []models.User{}
		_, err = sh.ShopServices.Update(shop)
		if err != nil {
			setFlashmessages(c, "error", "Can't update shop nextBillingDate")
			return sh.CreatePage(c)
		}
	}
	// --> END <--

	setFlashmessages(c, "success", "subscription created successfully!!")
	return c.Redirect(http.StatusSeeOther, "/subscription")
}

func (sh *SubscriptionHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Subscription | Update"
	id := c.Param("id")
	projectName := c.Param("name")
	subscription, err := sh.SubscriptionServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("membership with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	projects, _, _ := sh.ProjectServices.GetALL(50, 1, "desc", "id", "", "ACTIVE")
	memberships, _, _ := sh.MembershipServices.GetALL(50, 1, "desc", "id", projectName, "", "")
	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscription_views.Update(subscription, projects, memberships),
	))
}

func (sh *SubscriptionHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	subscription, err := sh.SubscriptionServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("subscription with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}
	oldSub := subscription
	if err := c.Bind(&subscription); err != nil {
		errorMsg = "cannot parse the subscription body"
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}

	if oldSub.PaymentStatus != "PAID" && subscription.PaymentStatus == "PAID" {
		now := time.Now()
		subscription.PaidAt = &now
	}

	_, err = sh.SubscriptionServices.Update(subscription)
	if err != nil {
		errorMsg = fmt.Sprintf("subscription with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}

	//NOTE: --> Handle statistic stuff <--
	err = sh.StatsServices.HandleStatsSubscription(oldSub, subscription)
	if err != nil {
		errorMsg = fmt.Sprintf("subscription with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}
	// --> END <--

	//NOTE: --> Get shop and update the nextBillingDate <--
	shop, err := sh.ShopServices.GetID(fmt.Sprintf("%d", subscription.ShopID), projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("subscription with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}
	//NOTE: --> we should do this so we don't get error when updating the shop <--
	shop.Owner = models.User{}
	shop.Workers = []models.User{}
	if shop.NextBillingDate == nil || shop.NextBillingDate.Before(subscription.EndDate) {
		shop.NextBillingDate = &subscription.EndDate
		_, err = sh.ShopServices.Update(shop)
		if err != nil {
			setFlashmessages(c, "error", "Can't create subscription")
			return sh.CreatePage(c)
		}
	}
	// --> END <--

	setFlashmessages(c, "success", "subscription updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/subscription")
}

func (sh *SubscriptionHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	subscription, err := sh.SubscriptionServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("subscription with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/subscription")
	}
	oldSubscription := subscription
	oldSubscription.ID = 0
	subscription.PaymentStatus = "NOT_PAID"
	subscriptionCopy := subscription
	_, err = sh.SubscriptionServices.Delete(subscriptionCopy)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete subscription with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/subscription")
	}

	//NOTE: --> we should do this so we don't get error when updating the shop <--
	err = sh.StatsServices.HandleStatsSubscription(oldSubscription, subscription)
	if err != nil {
		errorMsg = fmt.Sprintf("subscription with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}
	// --> END <--
	setFlashmessages(c, "success", "subscription successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/subscription")
}
