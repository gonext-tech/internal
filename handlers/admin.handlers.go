package handlers

import (
	"fmt"
	"net/http"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"github.com/gonext-tech/internal/views/admin_views"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	GetALL(queries.InvoiceQueryParams) ([]models.Admin, models.Meta, error)
	GetID(id string) (*models.Admin, error)
	Create(*models.Admin) error
	Update(*models.Admin) error
	Delete(*models.Admin) error
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
	var query queries.InvoiceQueryParams
	if err := c.Bind(&query); err != nil {
		errorMsg = "can't read query params"
		setFlashmessages(c, "error", errorMsg)
	}
	query.SetDefaults()
	admins, meta, err := ah.AdminServices.GetALL(query)
	if err != nil {
		isError = false
		errorMsg = "can't fetch admins"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	if c.Request().Header.Get("X-Partial-Content") == "true" {
		// Return only the table content
		return renderView(c, admin_views.List(
			fmt.Sprintf("Admin List(%d)", meta.TotalCount),
			admins,
			meta,
			query,
		))
	}
	titlePage := fmt.Sprintf(
		"Admin List(%d)", meta.TotalCount)
	return renderView(c, admin_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		admin_views.List(titlePage, admins, meta, query),
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

	newAdmin := &models.Admin{
		Email:    admin.Email,
		Phone:    admin.Phone,
		Address:  admin.Address,
		Password: string(hashedPassword),
		Status:   admin.Status,
		Role:     admin.Role,
		Image:    admin.Image,
	}

	err = ah.AdminServices.Create(newAdmin)
	if err != nil {
		setFlashmessages(c, "error", "cannot create a new admin")
		return ah.CreatePage(c)
	}
	imageURLs := UploadImage(c, ah.UploadServices, "", fmt.Sprintf("admin/%d", newAdmin.ID))

	if len(imageURLs) > 0 {
		newAdmin.Image = imageURLs[0]
		err = ah.AdminServices.Update(newAdmin)
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

	if err := c.Bind(admin); err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return ah.UpdatePage(c)
	}

	err = ah.AdminServices.Update(admin)
	if err != nil {
		errorMsg = fmt.Sprintf("cannot update admin with id %s ", id)
		setFlashmessages(c, "error", errorMsg)
		return ah.UpdatePage(c)
	}
	imageURLs := UploadImage(c, ah.UploadServices, "", fmt.Sprintf("admin/%d", admin.ID))

	if len(imageURLs) > 0 {
		admin.Image = imageURLs[0]
		err = ah.AdminServices.Update(admin)
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
