package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/referal_views"
	"github.com/labstack/echo/v4"
)

type ReferalService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Referal, models.Meta, error)
	GetID(id string) (models.Referal, error)
	Create(models.Referal) (models.Referal, error)
	Update(models.Referal) (models.Referal, error)
	Delete(models.Referal) error
}

type ReferalHandler struct {
	ReferalServices ReferalService
	UploadServices  UploadService
}

func NewReferalHandler(rs ReferalService, us UploadService) *ReferalHandler {
	return &ReferalHandler{
		ReferalServices: rs,
		UploadServices:  us,
	}
}

func (rh *ReferalHandler) ListPage(c echo.Context) error {
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
	referal, meta, err := rh.ReferalServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch referal"
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
		"Referal List(%d)", meta.TotalCount)
	return renderView(c, referal_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		referal_views.List(titlePage, referal, meta, params),
	))
}

func (rh *ReferalHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	referal, err := rh.ReferalServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch referal with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/referal")
	}
	titlePage := fmt.Sprintf(
		"Referal | %s", referal.Name)
	return renderView(c, referal_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		referal_views.View(referal),
	))
}

func (rh *ReferalHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Referal | Create"
	return renderView(c, referal_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		referal_views.Create(),
	))
}

func (rh *ReferalHandler) CreateHandler(c echo.Context) error {
	var referal models.Referal
	if err := c.Bind(&referal); err != nil {
		return err
	}
	_, err := rh.ReferalServices.Create(referal)
	if err != nil {
		return err
	}
	imageURLs := UploadImage(c, rh.UploadServices, "internal", fmt.Sprintf("referal/%d", referal.ID))

	if len(imageURLs) > 0 {
		referal.Image = imageURLs[0]
		_, err = rh.ReferalServices.Update(referal)
		if err != nil {
			setFlashmessages(c, "error", "Can't create project")
			return rh.CreatePage(c)
		}
	}
	setFlashmessages(c, "success", "referal created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/referal")
}

func (rh *ReferalHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Referal | Update"
	id := c.Param("id")
	referal, err := rh.ReferalServices.GetID(id)
	log.Println("referalll", referal)
	if err != nil {
		errorMsg = fmt.Sprintf("referal with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	return renderView(c, referal_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		referal_views.Update(referal),
	))
}

func (rh *ReferalHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")

	referal, err := rh.ReferalServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("referal with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return rh.UpdatePage(c)
	}

	if err := c.Bind(&referal); err != nil {
		log.Println("err", err)
		errorMsg = "cannot parse the referal body"
		setFlashmessages(c, "error", errorMsg)
		return rh.UpdatePage(c)
	}

	imageURLs := UploadImage(c, rh.UploadServices, "internal", fmt.Sprintf("referal/%d", referal.ID))

	if len(imageURLs) > 0 {
		referal.Image = imageURLs[0]
	}

	referal, err = rh.ReferalServices.Update(referal)
	if err != nil {
		errorMsg = fmt.Sprintf("referal with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return rh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "referal updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/referal")
}

func (rh *ReferalHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	referal, err := rh.ReferalServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("referal with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/referal")
	}
	err = rh.ReferalServices.Delete(referal)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete referal with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/referal")
	}
	setFlashmessages(c, "success", "Referal successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/referal")
}
