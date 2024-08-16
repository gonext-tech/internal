package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/city_views"
)

type CityService interface {
	GetALL() ([]models.City, error)
	GetID(id string) (models.City, error)
	Create(models.City) (models.City, error)
	Update(models.City) (models.City, error)
	Delete(models.City) error
}

type CityHandler struct {
	CityServices CityService
}

func NewCityHandler(cs CityService) *CityHandler {
	return &CityHandler{
		CityServices: cs,
	}
}

func (ch *CityHandler) ListPage(c echo.Context) error {
	isError = false
	cities, err := ch.CityServices.GetALL()
	if err != nil {
		isError = true
		errorMsg = "can't fetch cities"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}
	titlePage := fmt.Sprintf(
		"City List(%d)", len(cities))
	return renderView(c, city_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		city_views.List(titlePage, cities),
	))
}

func (ch *CityHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	city, err := ch.CityServices.GetID(id)
	if err != nil {
		isError = true
		errorMsg = fmt.Sprintf("can't fetch city with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/city")
	}
	titlePage := fmt.Sprintf(
		"City | %s", city.Name)

	return renderView(c, city_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		city_views.View(city),
	))
}

func (ch *CityHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "City | Create"
	return renderView(c, city_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		city_views.Create(),
	))
}

func (ch *CityHandler) CreateHandler(c echo.Context) error {
	isError = false
	var city models.City
	if err := c.Bind(&city); err != nil {
		setFlashmessages(c, "error", err.Error())
		return ch.CreatePage(c)
	}
	_, err := ch.CityServices.Create(city)
	if err != nil {
		setFlashmessages(c, "error", err.Error())
		return ch.CreatePage(c)
	}
	setFlashmessages(c, "success", "City created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/city")
}

func (ch *CityHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "City | Update"
	id := c.Param("id")
	city, err := ch.CityServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("city with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}
	return renderView(c, city_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		city_views.Update(city),
	))
}

func (ch *CityHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	city, err := ch.CityServices.GetID(id)
	if err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}
	if err := c.Bind(&city); err != nil {
		isError = true
		errorMsg = "cannot parse the city body"
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}
	city, err = ch.CityServices.Update(city)
	if err != nil {
		errorMsg = err.Error()
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}
	setFlashmessages(c, "success", "city updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/city")
}

func (ch *CityHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	blood, err := ch.CityServices.GetID(id)
	if err != nil {
		isError = true
		errorMsg = fmt.Sprintf("city with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/city")
	}
	err = ch.CityServices.Delete(blood)
	if err != nil {
		isError = true
		errorMsg = fmt.Sprintf("couldnt delete city with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/city")
	}
	setFlashmessages(c, "success", "city successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/city")
}
