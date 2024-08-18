package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/components"
	"github.com/gonext-tech/internal/views/membership_views"
	"github.com/labstack/echo/v4"
)

type MembershipService interface {
	GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) ([]models.Membership, models.Meta, error)
	Fetch(string) ([]models.Membership, error)
	GetID(id, name string) (models.Membership, error)
	Create(models.Membership) (models.Membership, error)
	Update(models.Membership) (models.Membership, error)
	Delete(models.Membership) (models.Membership, error)
}

type MembershipHandler struct {
	MembershipServices MembershipService
	ProjectServices    ProjectService
}

func NewMembershipHandler(ms MembershipService, ps ProjectService) *MembershipHandler {
	return &MembershipHandler{
		MembershipServices: ms,
		ProjectServices:    ps,
	}
}

func (mh *MembershipHandler) Fetch(c echo.Context) error {
	projectName := c.QueryParam("project_name")
	memberships, _ := mh.MembershipServices.Fetch(projectName)
	return renderView(c, components.MembershipResult(memberships))
}

func (mh *MembershipHandler) ListPage(c echo.Context) error {
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
	response, meta, err := mh.MembershipServices.GetALL(limit, page, orderBy, sortBy, project, status, searchTerm)
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
	return renderView(c, membership_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		membership_views.List(titlePage, response, meta, params),
	))
}

func (mh *MembershipHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	membership, err := mh.MembershipServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch hospital with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/membership")
	}
	titlePage := fmt.Sprintf(
		"Membership | %s", membership.Name)
	return renderView(c, membership_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		membership_views.View(membership),
	))
}

func (mh *MembershipHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Membership | Create"
	projects, _, _ := mh.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, membership_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		membership_views.Create(projects),
	))
}

func (mh *MembershipHandler) CreateHandler(c echo.Context) error {
	var membership models.Membership
	if err := c.Bind(&membership); err != nil {
		setFlashmessages(c, "error", errorMsg)
		return mh.CreatePage(c)
	}
	_, err := mh.MembershipServices.Create(membership)
	if err != nil {
		setFlashmessages(c, "error", "Can't create membership")
		return mh.CreatePage(c)
	}
	setFlashmessages(c, "success", "membership created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/membership")
}

func (mh *MembershipHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Membership | Update"
	id := c.Param("id")
	projectName := c.Param("name")
	membership, err := mh.MembershipServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("membership with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	projects, _, _ := mh.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, membership_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		membership_views.Update(membership, projects),
	))
}

func (mh *MembershipHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("projectName")
	membership, err := mh.MembershipServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("membership with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return mh.UpdatePage(c)
	}

	if err := c.Bind(&membership); err != nil {
		errorMsg = "cannot parse the membership body"
		setFlashmessages(c, "error", errorMsg)
		return mh.UpdatePage(c)
	}

	_, err = mh.MembershipServices.Update(membership)
	if err != nil {
		errorMsg = fmt.Sprintf("membership with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return mh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "membership updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/membership")
}

func (mh *MembershipHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	subscription, err := mh.MembershipServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("membership with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/membership")
	}
	_, err = mh.MembershipServices.Delete(subscription)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete membership with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/membership")
	}
	setFlashmessages(c, "success", "Donor successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/membership")
}
