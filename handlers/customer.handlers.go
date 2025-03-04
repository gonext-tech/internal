package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"golang.org/x/crypto/bcrypt"

	//"github.com/gonext-tech/internal/views/components"
	"github.com/gonext-tech/internal/views/components"
	"github.com/gonext-tech/internal/views/components/user_components"
	"github.com/gonext-tech/internal/views/customer_views"
	"github.com/gonext-tech/internal/views/partials"
	"github.com/labstack/echo/v4"
)

type CustomerService interface {
	GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) ([]models.User, models.Meta, error)
	GetID(id, name string) (models.User, error)
	Create(models.User) (models.User, error)
	Update(models.User) (models.User, error)
	Delete(models.User) (models.User, error)
}

type CustomerHandler struct {
	CustomerServices CustomerService
	ProjectServices  ProjectService
	UploadServices   UploadService
}

func NewCustomerHandler(cs CustomerService, ps ProjectService, us UploadService) *CustomerHandler {
	return &CustomerHandler{
		CustomerServices: cs,
		ProjectServices:  ps,
		UploadServices:   us,
	}
}

func (ch *CustomerHandler) ListPage(c echo.Context) error {
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
	customers, meta, err := ch.CustomerServices.GetALL(limit, page, orderBy, sortBy, project, status, searchTerm)
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
		"Customers (%d)", meta.TotalCount)
	return renderView(c, customer_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		customer_views.List(titlePage, customers, meta, params),
	))
}

func (ch *CustomerHandler) SearchUser(c echo.Context) error {
	isError = false
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 5
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
	project := c.QueryParam("project_name")
	searchTerm := c.QueryParam("searchTerm")
	customers, _, err := ch.CustomerServices.GetALL(limit, page, orderBy, sortBy, project, status, searchTerm)
	if err != nil {
		isError = false
		errorMsg = "can't fetch customers"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	return renderView(c, components.CustomerResult(customers))
}

func (ch *CustomerHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	customer, err := ch.CustomerServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch hospital with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/membership")
	}
	titlePage := fmt.Sprintf(
		"Customer | %s", customer.Name)
	return renderView(c, customer_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		customer_views.View(customer),
	))
}

func (ch *CustomerHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Customer | Create"
	projects, _, _ := ch.ProjectServices.GetALL(50, 1, "desc", "id", "", "ACTIVE")
	return renderView(c, customer_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		customer_views.Create(projects),
	))
}

func (ch *CustomerHandler) CreateHandler(c echo.Context) error {
	var customer models.User
	if err := c.Bind(&customer); err != nil {
		setFlashmessages(c, "error", errorMsg)
		return ch.CreatePage(c)
	}
	source := c.Request().Header.Get("HX-Source")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	customer.Password = string(hashedPassword)
	customer, err = ch.CustomerServices.Create(customer)
	if err != nil {
		setFlashmessages(c, "error", "Can't create customer")
		return ch.CreatePage(c)
	}

	imageURLs := UploadImage(c, ch.UploadServices, customer.ProjectName, fmt.Sprintf("cutsomer/%d", customer.ID))

	if len(imageURLs) > 0 {
		customer.Image = imageURLs[0]
		_, err = ch.CustomerServices.Update(customer)
		if err != nil {
			setFlashmessages(c, "error", "Can't create customer")
			return ch.CreatePage(c)
		}
	}
	setFlashmessages(c, "success", "customer created successfully!!")
	if source != "" {

		successMsgs := getFlashmessages(c, "success")
		successAlert := partials.FlashMessages(nil, successMsgs)
		userCard := user_components.UserCard(customer)

		err := renderView(c, successAlert)
		if err != nil {
			return err
		}
		return renderView(c, userCard)

	}
	return c.Redirect(http.StatusSeeOther, "/customer")
}

func (ch *CustomerHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Customer | Update"
	id := c.Param("id")
	projectName := c.Param("name")
	customer, err := ch.CustomerServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("customer with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	projects, _, _ := ch.ProjectServices.GetALL(50, 1, "desc", "id", "", "ACTIVE")
	return renderView(c, customer_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		customer_views.Update(customer, projects),
	))
}

func (ch *CustomerHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	customer, err := ch.CustomerServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("customer with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}
	// Retrieve uploaded files
	imageURLs := UploadImage(c, ch.UploadServices, "internal", fmt.Sprintf("customer/%d", customer.ID))

	oldPassword := customer.Password
	if err := c.Bind(&customer); err != nil {
		errorMsg = "cannot parse the customer body"
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}
	if len(imageURLs) > 0 {
		customer.Image = imageURLs[0]
	}
	if oldPassword != customer.Password {
		err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(customer.Password))
		if err != nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
			if err != nil {
				errorMsg = "cannot parse the customer password , please contact support"
				setFlashmessages(c, "error", errorMsg)
				return ch.UpdatePage(c)
			}
			customer.Password = string(hashedPassword)
		}
	} else {
		customer.Password = oldPassword
	}

	_, err = ch.CustomerServices.Update(customer)
	if err != nil {
		errorMsg = fmt.Sprintf("customer with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ch.UpdatePage(c)
	}
	setFlashmessages(c, "success", "customer updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/customer")
}

func (ch *CustomerHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	source := c.Request().Header.Get("HX-Source")
	customer, err := ch.CustomerServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("customer with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		if source != "" {
			errMsgs := getFlashmessages(c, "error")
			return renderView(c, partials.FlashMessages(errMsgs, nil), 400)
		}
		return c.Redirect(http.StatusSeeOther, "/customer")
	}
	_, err = ch.CustomerServices.Delete(customer)

	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete customer with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		if source != "" {
			errMsgs := getFlashmessages(c, "error")
			return renderView(c, partials.FlashMessages(errMsgs, nil), 400)
		}
		return c.Redirect(http.StatusSeeOther, "/customer")
	}
	setFlashmessages(c, "success", "Customer successfully deleted!!")
	if source != "" {
		successMsgs := getFlashmessages(c, "success")
		return renderView(c, partials.FlashMessages(nil, successMsgs))
	}
	return c.Redirect(http.StatusSeeOther, "/customer")
}
