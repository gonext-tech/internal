package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/project_views"
	"github.com/labstack/echo/v4"
)

type ProjectService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Project, models.Meta, error)
	GetID(id string) (models.Project, error)
	Create(*models.Project) (*models.Project, error)
	Update(models.Project) (models.Project, error)
	Delete(models.Project) error
}

type ProjectHandler struct {
	ProjectServices ProjectService
	ServerServices  ServerService
	AdminServices   AdminService
	UploadServices  UploadService
}

func NewProjectHandler(ps ProjectService, us UploadService, ss ServerService, u AdminService) *ProjectHandler {
	return &ProjectHandler{
		ProjectServices: ps,
		UploadServices:  us,
		AdminServices:   u,
		ServerServices:  ss,
	}
}

func (ph *ProjectHandler) ListPage(c echo.Context) error {
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
	projects, meta, err := ph.ProjectServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch projects"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	if err != nil {
		isError = false
		errorMsg = "can't fetch projects"
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
		"Project List(%d)", meta.TotalCount)
	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.List(titlePage, projects, meta, params),
	))
}

func (ph *ProjectHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	project, err := ph.ProjectServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch project with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/project")
	}
	titlePage := fmt.Sprintf(
		"Project | %s", project.Name)
	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.View(project),
	))
}

func (ph *ProjectHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Project | Create"
	servers, _, _ := ph.ServerServices.GetALL(50, 1, "desc", "id", "", "UP")
	leads, _, _ := ph.AdminServices.GetALL(50, 1, "desc", "id", "", "ACTIVE")
	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.Create(servers, leads),
	))
}

func (ph *ProjectHandler) CreateHandler(c echo.Context) error {
	var project models.Project
	if err := c.Bind(&project); err != nil {
		setFlashmessages(c, "error", "cannot parse project body")
		return ph.CreatePage(c)
	}
	_, err := ph.ProjectServices.Create(&project)
	if err != nil {
		setFlashmessages(c, "error", err.Error())
		return ph.CreatePage(c)
	}
	imageURLs := UploadImage(c, ph.UploadServices, "internal", fmt.Sprintf("project/%d", project.ID))

	if len(imageURLs) > 0 {
		project.File = imageURLs[0]
		_, err = ph.ProjectServices.Update(project)
		if err != nil {
			setFlashmessages(c, "error", "Can't create project")
			return ph.CreatePage(c)
		}
	}
	setFlashmessages(c, "success", "project created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/project")
}

func (ph *ProjectHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Project | Update"
	id := c.Param("id")
	project, err := ph.ProjectServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("project with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	servers, _, _ := ph.ServerServices.GetALL(50, 1, "desc", "id", "", "UP")
	leads, _, _ := ph.AdminServices.GetALL(50, 1, "desc", "id", "", "ACTIVE")

	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.Update(project, servers, leads),
	))
}

func (ph *ProjectHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")

	project, err := ph.ProjectServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("project with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ph.UpdatePage(c)
	}
	if err := c.Bind(&project); err != nil {
		errorMsg = "cannot parse the project body"
		setFlashmessages(c, "error", errorMsg)
		return ph.UpdatePage(c)
	}
	imageURLs := UploadImage(c, ph.UploadServices, "internal", fmt.Sprintf("project/%d", project.ID))

	if len(imageURLs) > 0 {
		project.File = imageURLs[0]
	}

	project, err = ph.ProjectServices.Update(project)
	if err != nil {
		errorMsg = fmt.Sprintf("project with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ph.UpdatePage(c)
	}
	setFlashmessages(c, "success", "project updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/project")
}

func (ph *ProjectHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	project, err := ph.ProjectServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("project with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/project")
	}
	err = ph.ProjectServices.Delete(project)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete project with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/project")
	}
	setFlashmessages(c, "success", "Project successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/project")
}
