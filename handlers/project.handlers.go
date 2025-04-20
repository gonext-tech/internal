package handlers

import (
	"fmt"
	"net/http"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"github.com/gonext-tech/internal/views/project_views"
	"github.com/labstack/echo/v4"
)

type ProjectService interface {
	GetALL(queries.InvoiceQueryParams) ([]models.Project, models.Meta, error)
	GetID(id string) (*models.Project, error)
	Create(*models.Project) error
	Update(*models.Project) error
	Delete(*models.Project) error
}

type ProjectHandler struct {
	ProjectServices ProjectService
	ServerServices  ServerService
	AdminServices   AdminService
	ClientServices  ClientService
	UploadServices  UploadService
}

func NewProjectHandler(ps ProjectService, us UploadService, ss ServerService, u AdminService, cs ClientService) *ProjectHandler {
	return &ProjectHandler{
		ProjectServices: ps,
		UploadServices:  us,
		AdminServices:   u,
		ClientServices:  cs,
		ServerServices:  ss,
	}
}

func (ph *ProjectHandler) ListPage(c echo.Context) error {
	isError = false
	var query queries.InvoiceQueryParams
	if err := c.Bind(&query); err != nil {
		errorMsg = "can't read query params"
		setFlashmessages(c, "error", errorMsg)
	}
	query.SetDefaults()

	projects, meta, err := ph.ProjectServices.GetALL(query)

	if c.Request().Header.Get("X-Partial-Content") == "true" {
		// Return only the table content
		return renderView(c, project_views.List(
			fmt.Sprintf("Project List(%d)", meta.TotalCount),
			projects,
			meta,
			query,
		))
	}
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

	titlePage := fmt.Sprintf(
		"Project List(%d)", meta.TotalCount)
	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.List(titlePage, projects, meta, query),
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
		project_views.View(*project),
	))
}

func (ph *ProjectHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Project | Create"
	queries := queries.InvoiceQueryParams{}
	queries.SearchDefaults()
	leads, _, _ := ph.AdminServices.GetALL(queries)
	clients, _, _ := ph.ClientServices.GetALL(queries)
	queries.Status = "UP"
	servers, _, _ := ph.ServerServices.GetALL(queries)
	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.Create(servers, leads, clients),
	))
}

func (ph *ProjectHandler) CreateHandler(c echo.Context) error {
	var project models.Project
	if err := c.Bind(&project); err != nil {
		setFlashmessages(c, "error", "cannot parse project body")
		return ph.CreatePage(c)
	}
	err := ph.ProjectServices.Create(&project)
	if err != nil {
		setFlashmessages(c, "error", err.Error())
		return ph.CreatePage(c)
	}
	imageURLs := UploadImage(c, ph.UploadServices, "internal", fmt.Sprintf("project/%d", project.ID))

	if len(imageURLs) > 0 {
		project.File = imageURLs[0]
		err = ph.ProjectServices.Update(&project)
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

	queries := queries.InvoiceQueryParams{}
	queries.SearchDefaults()

	leads, _, _ := ph.AdminServices.GetALL(queries)
	clients, _, _ := ph.ClientServices.GetALL(queries)
	queries.Status = "UP"
	servers, _, _ := ph.ServerServices.GetALL(queries)

	return renderView(c, project_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		project_views.Update(*project, servers, leads, clients),
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
	if err := c.Bind(project); err != nil {
		errorMsg = "cannot parse the project body"
		setFlashmessages(c, "error", errorMsg)
		return ph.UpdatePage(c)
	}
	imageURLs := UploadImage(c, ph.UploadServices, "internal", fmt.Sprintf("project/%d", project.ID))

	if len(imageURLs) > 0 {
		project.File = imageURLs[0]
	}

	err = ph.ProjectServices.Update(project)
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
