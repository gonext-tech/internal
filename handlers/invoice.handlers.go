package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"github.com/gonext-tech/internal/views/invoice_views"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type InvoiceService interface {
	GetALL(queries.InvoiceQueryParams) ([]models.Invoice, models.Meta, error)
	GetID(id string) (*models.Invoice, error)
	Create(*models.Invoice) error
	Update(*models.Invoice) error
	Delete(*models.Invoice) error
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
	var query queries.InvoiceQueryParams
	if err := c.Bind(&query); err != nil {
		errorMsg = "can't read query params"
		setFlashmessages(c, "error", errorMsg)
	}
	query.SetDefaults()
	expenses, meta, err := ih.InvoiceServices.GetALL(query)
	if err != nil {
		isError = false
		errorMsg = "can't fetch invoice"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	if c.Request().Header.Get("X-Partial-Content") == "true" {
		// Return only the table content
		return renderView(c, invoice_views.List(
			fmt.Sprintf("Invoice List(%d)", meta.TotalCount),
			expenses,
			meta,
			query,
		))
	}

	titlePage := fmt.Sprintf(
		"Invoice List(%d)", meta.TotalCount)
	return renderView(c, invoice_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		invoice_views.List(titlePage, expenses, meta, query),
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
	projects, _, _ := ih.ProjectServices.GetALL(
		queries.InvoiceQueryParams{
			Status:  "ACTIVE",
			OrderBy: "desc",
			SortBy:  "id",
			Page:    1,
			Limit:   50,
		})
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

	if !ok {
		setFlashmessages(c, "error", "user is not authenticated")
		return ih.CreatePage(c)
	}

	var invoice models.Invoice
	if err := c.Bind(&invoice); err != nil {
		return err
	}
	log.Println("userIDDD", userID)
	invoice.CreatedByID = userID
	err := ih.InvoiceServices.Create(&invoice)
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

	projects, _, _ := ih.ProjectServices.GetALL(queries.InvoiceQueryParams{
		Status:  "ACTIVE",
		OrderBy: "desc",
		SortBy:  "id",
		Page:    1,
		Limit:   50,
	})

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
	if err := c.Bind(invoice); err != nil {
		errorMsg = "cannot parse the project body"
		setFlashmessages(c, "error", errorMsg)
		return ih.UpdatePage(c)
	}
	log.Println("invoice-clientttt", invoice.ClientID)

	err = ih.InvoiceServices.Update(invoice)
	if err != nil {
		errorMsg = fmt.Sprintf("invoice with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return ih.UpdatePage(c)
	}
	setFlashmessages(c, "success", "invoice updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/invoice")
}

func (ih *InvoiceHandler) PaidHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	userID, ok := sess.Values[user_id_key].(uint)
	if !ok {
		setFlashmessages(c, "error", "Unauthorized access")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	id := c.Param("id")
	invoice, err := ih.InvoiceServices.GetID(id)
	if err != nil {
		errorMsg := fmt.Sprintf("Invoice with ID %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/invoice")
	}
	paidAmountStr := c.FormValue("amount_paid")
	notes := c.FormValue("notes")

	paidAmount, err := strconv.ParseFloat(paidAmountStr, 64)
	if err != nil || paidAmount <= 0 || paidAmount > invoice.Total || invoice.AmountPaid+paidAmount > invoice.Total {
		return renderView(c, invoice_views.PayError(invoice.ID, "Amount is incorrect"))
	}

	now := time.Now()
	// If full payment, just update the main invoice
	if paidAmount == invoice.Total {
		invoice.AmountPaid = paidAmount
		invoice.PaymentStatus = "PAID"
		invoice.PaidToID = userID
		invoice.PaidAt = &now
		invoice.Notes = notes
		err = ih.InvoiceServices.Update(invoice)
		if err != nil {
			return renderView(c, invoice_views.PayError(invoice.ID, "Failed to update invoice"))
		}
		setFlashmessages(c, "success", "Invoice fully paid")
		c.Response().Header().Set("HX-Refresh", "true")
		return c.NoContent(http.StatusOK)
	}

	invoice.AmountPaid += paidAmount
	if invoice.AmountPaid == invoice.Total {
		invoice.PaymentStatus = "PAID"
	} else {
		invoice.PaymentStatus = "P_PAID"
	}
	invoice.PaidToID = userID
	invoice.PaidAt = &now
	invoice.Notes = notes
	err = ih.InvoiceServices.Update(invoice)

	// Else: create a sub-invoice
	subInvoice := &models.Invoice{
		PaymentStatus:   "PAID",
		Amount:          paidAmount,
		AmountPaid:      paidAmount,
		Discount:        0,
		Total:           paidAmount,
		ClientID:        invoice.ClientID,
		ProjectID:       invoice.ProjectID,
		InvoiceRefID:    &invoice.ID,
		Recurring:       invoice.Recurring,
		RecurringPeriod: invoice.RecurringPeriod,
		Category:        invoice.Category,
		CreatedByID:     userID,
		PaidToID:        userID,
		InvoiceType:     invoice.InvoiceType,
		IssueDate:       invoice.IssueDate,
		DueDate:         invoice.DueDate,
		PaidAt:          &now,
		Notes:           notes,
	}

	err = ih.InvoiceServices.Create(subInvoice)
	if err != nil {
		return renderView(c, invoice_views.PayError(invoice.ID, "Failed to create sub-invoice"))
	}

	setFlashmessages(c, "success", "Partial payment recorded")

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (ih *InvoiceHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	invoice, err := ih.InvoiceServices.GetID(id)
	if err != nil {
		errorMsg = fmt.Sprintf("invoice with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		c.Response().Header().Set("HX-Refresh", "true")
		return c.NoContent(http.StatusOK)
	}
	err = ih.InvoiceServices.Delete(invoice)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete invoice with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		c.Response().Header().Set("HX-Refresh", "true")
		return c.NoContent(http.StatusOK)
	}
	setFlashmessages(c, "success", "Invoice successfully deleted!!")
	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (ih *InvoiceHandler) GenerateInvoicePDF(c echo.Context) error {
	id := c.Param("id")
	invoice, err := ih.InvoiceServices.GetID(id)
	if err != nil {
		return err
	}
	return renderView(c, invoice_views.PrintInvoicePage(invoice))

}
