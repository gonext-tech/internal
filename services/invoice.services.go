package services

import (
	"errors"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"gorm.io/gorm"
)

type InvoiceService struct {
	DB *gorm.DB
}

func NewInvoiceService(db *gorm.DB) *InvoiceService {
	return &InvoiceService{
		DB: db,
	}
}

func (is *InvoiceService) GetALL(query queries.InvoiceQueryParams) ([]models.Invoice, models.Meta, error) {
	var invoices []models.Invoice
	var totalCount int64

	dbQuery := is.buildInvoicesURL(query)
	dbQuery.Session(&gorm.Session{}).Model(&models.Invoice{}).Count(&totalCount)
	offset := (query.Page - 1) * query.Limit
	dbQuery.Order(query.SortBy + " " + query.OrderBy).Offset(offset).Limit(query.Limit).Find(&invoices)
	lastPage := int64(0)
	if query.Limit > 0 {
		lastPage = (totalCount + int64(query.Limit) - 1) / int64(query.Limit)
	}
	meta := models.Meta{
		CurrentPage: query.Page,
		TotalCount:  int(totalCount),
		LastPage:    int(lastPage),
		Limit:       query.Limit,
	}

	return invoices, meta, nil
}

func (is *InvoiceService) GetID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	if id == "0" || id == "" {
		return nil, errors.New("no id provided")
	}
	query := is.DB.Preload("Project").Preload("Client").Preload("CreatedBy").Preload("PaidTo").Preload("CanceledBy").Preload("SubInvoices").Preload("SubInvoices.Project").Preload("SubInvoices.Client").Preload("SubInvoices.CreatedBy").Preload("SubInvoices.PaidTo").Preload("SubInvoices.CanceledBy")
	if result := query.First(&invoice, id); result.Error != nil {
		return nil, result.Error
	}
	return &invoice, nil
}

func (is *InvoiceService) Create(invoice *models.Invoice) error {
	return is.DB.Create(invoice).Error
}

func (is *InvoiceService) Update(invoice *models.Invoice) error {
	return is.DB.Model(&models.Invoice{}).
		Where("id = ?", invoice.ID).
		Select("*").
		Updates(invoice).Error
}

func (is *InvoiceService) Delete(invoice *models.Invoice) error {
	return is.DB.Model(&invoice).Updates(invoice).Error
}

func (is *InvoiceService) buildInvoicesURL(query queries.InvoiceQueryParams) *gorm.DB {
	dbQuery := is.DB.
		Preload("Project").
		Preload("Client").
		Preload("SubInvoices").
		Preload("CreatedBy").
		Preload("PaidTo").
		Preload("CanceledBy")

	dbQuery = dbQuery.Where("invoice_ref_id IS NULL")
	if query.SearchTerm != "" {
		term := "%" + query.SearchTerm + "%"
		dbQuery = dbQuery.Where("name LIKE ?", term)
	}

	if query.InvoiceType != "" {
		dbQuery = dbQuery.Where("invoice_type = ?", query.InvoiceType)
	}

	return dbQuery

}
