package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/donor_views"
)

type DonorService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Donor, models.Meta, error)
	GetID(id string) (models.Donor, error)
	Create(models.Donor) (models.Donor, error)
	Update(models.Donor) (models.Donor, error)
	Delete(models.Donor) error
}

type DonorHandler struct {
	DonorServices DonorService
	BloodServices BloodService
	CityServices  CityService
}

func NewDonorHandler(ds DonorService, bs BloodService, cs CityService) *DonorHandler {
	return &DonorHandler{
		DonorServices: ds,
		BloodServices: bs,
		CityServices:  cs,
	}
}

func (dh *DonorHandler) ListPage(c echo.Context) error {
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
	donors, meta, err := dh.DonorServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch donors"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}
	titlePage := fmt.Sprintf(
		"Project List(%d)", meta.TotalCount)
	return renderView(c, donor_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		donor_views.List(titlePage, donors, meta),
	))
}

func (dh *DonorHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	donor, err := dh.DonorServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch donor with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/donor")
	}
	titlePage := fmt.Sprintf(
		"Donor | %s", donor.Name)
	return renderView(c, donor_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		donor_views.View(donor),
	))
}

func (dh *DonorHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Donor | Create"
	cities, _ := dh.CityServices.GetALL()
	bloods, _ := dh.BloodServices.GetALL()
	return renderView(c, donor_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		donor_views.Create(cities, bloods),
	))
}

func (dh *DonorHandler) CreateHandler(c echo.Context) error {
	var donorBody models.DonorBody
	if err := c.Bind(&donorBody); err != nil {
		setFlashmessages(c, "error", errorMsg)
		return dh.CreatePage(c)
	}
	parsedDate, err := time.Parse("2006-01-02", donorBody.DateOfBirth) // Adjust format as needed
	if err != nil {
		setFlashmessages(c, "error", "Invalid date format")
		return dh.CreatePage(c)
	}
	donor := models.Donor{
		Name:        donorBody.Name,
		Gender:      donorBody.Gender,
		Phone:       donorBody.Phone,
		CityID:      donorBody.CityID,
		BloodTypeID: donorBody.BloodTypeID,
		Address:     donorBody.Address,
		DateOfBirth: parsedDate,
	}
	_, err = dh.DonorServices.Create(donor)
	if err != nil {
		setFlashmessages(c, "error", "Can't create donor")
		return dh.CreatePage(c)
	}
	setFlashmessages(c, "success", "donor created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/donor")
}

func (dh *DonorHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Donor | Update"
	id := c.Param("id")
	donor, err := dh.DonorServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("donor with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	cities, _ := dh.CityServices.GetALL()
	bloods, _ := dh.BloodServices.GetALL()
	return renderView(c, donor_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		donor_views.Update(donor, bloods, cities),
	))
}

func (dh *DonorHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	donor, err := dh.DonorServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("project with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return dh.UpdatePage(c)
	}

	var donorBody models.DonorBody
	if err := c.Bind(&donorBody); err != nil {
		errorMsg = "cannot parse the donor body"
		setFlashmessages(c, "error", errorMsg)
		return dh.UpdatePage(c)
	}
	parsedDate, err := time.Parse("2006-01-02", donorBody.DateOfBirth) // Adjust format as needed
	if err != nil {
		setFlashmessages(c, "error", "Invalid date format")
		return dh.CreatePage(c)
	}

	updatedDonor := models.Donor{
		ID:          donor.ID,
		Name:        donorBody.Name,
		Gender:      donorBody.Gender,
		Phone:       donorBody.Phone,
		CityID:      donorBody.CityID,
		BloodTypeID: donorBody.BloodTypeID,
		Address:     donorBody.Address,
		DateOfBirth: parsedDate,
	}
	_, err = dh.DonorServices.Update(updatedDonor)
	if err != nil {
		errorMsg = fmt.Sprintf("donor with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return dh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "donor updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/donor")
}

func (dh *DonorHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	donor, err := dh.DonorServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("donor with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/donor")
	}
	err = dh.DonorServices.Delete(donor)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete donor with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/donor")
	}
	setFlashmessages(c, "success", "Donor successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/donor")
}
