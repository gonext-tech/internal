package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/invoice_views"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type InvoiceService interface {
	GetALL(limit, page int, orderBy, sortBy, invoiceType, searchTerm, status string) ([]models.Invoice, models.Meta, error)
	GetID(id string) (models.Invoice, error)
	Create(models.Invoice) (models.Invoice, error)
	Update(models.Invoice) (models.Invoice, error)
	Delete(models.Invoice) error
}

type InvoiceHandler struct {
	InvoiceServices InvoiceService
	ProjectServices ProjectService
}

func NewInvoiceHandler(is InvoiceService, ps ProjectService) *InvoiceHandler {
	return &InvoiceHandler{
		InvoiceServices: is,
		ProjectServices: ps,
	}
}

func (ih *InvoiceHandler) ListPage(c echo.Context) error {
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
	invoiceType := c.QueryParam("type")
	status := c.QueryParam("status")
	searchTerm := c.QueryParam("searchTerm")
	expenses, meta, err := ih.InvoiceServices.GetALL(limit, page, orderBy, sortBy, invoiceType, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch invoice"
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
		"Invoice List(%d)", meta.TotalCount)
	return renderView(c, invoice_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		invoice_views.List(titlePage, expenses, meta, params),
	))
}

func (ih *InvoiceHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	invoice, err := ih.InvoiceServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch invoice with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/invoice")
	}
	titlePage := fmt.Sprintf(
		"Invoice | #%d", invoice.ID)
	return renderView(c, invoice_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		invoice_views.View(invoice),
	))
}

func (ih *InvoiceHandler) CreatePage(c echo.Context) error {

	isError = false
	titlePage := "Invoice | Create"
	projects, _, _ := ih.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, invoice_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		invoice_views.Create(projects),
	))
}

func (ih *InvoiceHandler) CreateHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	userID, ok := sess.Values[user_id_key].(uint)
	log.Println("userrrId", userID)
	if !ok {
		setFlashmessages(c, "error", "user is not authenticated")
		return ih.CreatePage(c)
	}

	var invoice models.Invoice
	if err := c.Bind(&invoice); err != nil {
		return err
	}
	invoice.UserID = userID
	_, err := ih.InvoiceServices.Create(invoice)
	if err != nil {
		return err
	}
	setFlashmessages(c, "success", "invoice created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/invoice")
}

func (ih *InvoiceHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Invoice | Update"
	id := c.Param("id")
	invoice, err := ih.InvoiceServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("invoice with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	projects, _, _ := ih.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, invoice_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		invoice_views.Update(invoice, projects),
	))
}

func (ih *InvoiceHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")

	invoice, err := ih.InvoiceServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("invoice with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ih.UpdatePage(c)
	}

	if err := c.Bind(&invoice); err != nil {
		log.Println("err", err)
		errorMsg = "cannot parse the project body"
		setFlashmessages(c, "error", errorMsg)
		return ih.UpdatePage(c)
	}

	invoice, err = ih.InvoiceServices.Update(invoice)
	if err != nil {
		errorMsg = fmt.Sprintf("invoice with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ih.UpdatePage(c)
	}
	setFlashmessages(c, "success", "invoice updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/invoice")
}

func (ih *InvoiceHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	invoice, err := ih.InvoiceServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("invoice with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/invoice")
	}
	err = ih.InvoiceServices.Delete(invoice)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete invoice with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/invoice")
	}
	setFlashmessages(c, "success", "Invoice successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/invoice")
}
