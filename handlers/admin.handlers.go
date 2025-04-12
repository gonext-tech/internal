package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/admin_views"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Admin, models.Meta, error)
	GetID(id string) (models.Admin, error)
	Create(models.Admin) (models.Admin, error)
	Update(models.Admin) (models.Admin, error)
	Delete(models.Admin) error
}

type AdminHandler struct {
	AdminServices  AdminService
	UploadServices UploadService
}

func NewAdminHandler(as AdminService, us UploadService) *AdminHandler {
	return &AdminHandler{
		AdminServices:  as,
		UploadServices: us,
	}
}

func (ah *AdminHandler) ListPage(c echo.Context) error {
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
	admins, meta, err := ah.AdminServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch admins"
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
		"Admin List(%d)", meta.TotalCount)
	return renderView(c, admin_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		admin_views.List(titlePage, admins, meta, params),
	))
}

func (ah *AdminHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	admin, err := ah.AdminServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch admin with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/admin")
	}
	titlePage := fmt.Sprintf(
		"Admin | %s", admin.Name)
	return renderView(c, admin_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		admin_views.View(admin),
	))
}

func (ah *AdminHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Admin | Create"
	return renderView(c, admin_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		admin_views.Create(),
	))
}

func (ah *AdminHandler) CreateHandler(c echo.Context) error {
	var admin models.AdminBody
	if err := c.Bind(&admin); err != nil {
		setFlashmessages(c, "error", "cannot parse admin body")
		return ah.CreatePage(c)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 8)

	if err != nil {
		setFlashmessages(c, "error", "cannot set new password")
		return ah.CreatePage(c)
	}

	newAdmin := models.Admin{
		Email:    admin.Email,
		Phone:    admin.Phone,
		Address:  admin.Address,
		Password: string(hashedPassword),
		Status:   admin.Status,
		Role:     admin.Role,
		Image:    admin.Image,
	}

	newAdmin, err = ah.AdminServices.Create(newAdmin)
	if err != nil {
		setFlashmessages(c, "error", "cannot create a new admin")
		return ah.CreatePage(c)
	}
	imageURLs := UploadImage(c, ah.UploadServices, "", fmt.Sprintf("admin/%d", newAdmin.ID))

	if len(imageURLs) > 0 {
		newAdmin.Image = imageURLs[0]
		_, err = ah.AdminServices.Update(newAdmin)
		if err != nil {
			setFlashmessages(c, "error", "cannot upload admin image")
			return ah.CreatePage(c)
		}
	}
	setFlashmessages(c, "success", "admin created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/admin")
}

func (ah *AdminHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Admin | Update"
	id := c.Param("id")
	admin, err := ah.AdminServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("admin with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	return renderView(c, admin_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		admin_views.Update(admin),
	))
}

func (ah *AdminHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")

	admin, err := ah.AdminServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("admin with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ah.UpdatePage(c)
	}

	if err := c.Bind(&admin); err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return ah.UpdatePage(c)
	}

	admin, err = ah.AdminServices.Update(admin)
	if err != nil {
		errorMsg = fmt.Sprintf("cannot update admin with id %s ", id)
		setFlashmessages(c, "error", errorMsg)
		return ah.UpdatePage(c)
	}
	imageURLs := UploadImage(c, ah.UploadServices, "", fmt.Sprintf("admin/%d", admin.ID))

	if len(imageURLs) > 0 {
		admin.Image = imageURLs[0]
		_, err = ah.AdminServices.Update(admin)
		if err != nil {
			setFlashmessages(c, "error", "cannot upload admin image")
			return ah.CreatePage(c)
		}
	}
	setFlashmessages(c, "success", "admin updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/admin")
}

func (ah *AdminHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	admin, err := ah.AdminServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("admin with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/admin")
	}
	err = ah.AdminServices.Delete(admin)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete admin with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/adinm")
	}
	setFlashmessages(c, "success", "Admin successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/admin")
}
