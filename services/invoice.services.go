package services

import (
	"errors"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type InvoiceService struct {
	Invoice models.Invoice
	DB      *gorm.DB
}

func NewInvoiceService(ex models.Invoice, db *gorm.DB) *InvoiceService {
	return &InvoiceService{
		Invoice: ex,
		DB:      db,
	}
}

func (is *InvoiceService) GetALL(limit, page int, orderBy, sortBy, invoiceType, searchTerm, status string) ([]models.Invoice, models.Meta, error) {
	var invoices []models.Invoice
	query := is.DB.Preload("User").Preload("Project")
	totalQuery := is.DB

	if searchTerm != "" {
		searchTermWithWildcard := "%" + searchTerm + "%"
		query = query.Where("name LIKE ?", searchTermWithWildcard)
		totalQuery = query
	}

	if status != "" {
		query = query.Where("status = ?", status)
		totalQuery = totalQuery.Where("status = ?", status)
	}

	if invoiceType != "" {
		query = query.Where("invoice_type = ?", invoiceType)
		totalQuery = totalQuery.Where("invoice_type = ?", invoiceType)
	}

	offset := (page - 1) * limit
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&invoices)
	totalRecords := int64(0)
	totalQuery.Model(&is.Invoice).Count(&totalRecords)
	lastPage := int64(0)
	if limit > 0 {
		lastPage = (totalRecords + int64(limit) - 1) / int64(limit)
	}
	meta := models.Meta{
		CurrentPage: page,
		TotalCount:  int(totalRecords),
		LastPage:    int(lastPage),
		Limit:       limit,
	}

	return invoices, meta, nil
}

func (is *InvoiceService) GetID(id string) (models.Invoice, error) {
	var invoice models.Invoice
	if id == "0" || id == "" {
		return is.Invoice, errors.New("no id provided")
	}
	if result := is.DB.Preload("User").Preload("Project").First(&invoice, id); result.Error != nil {
		return is.Invoice, result.Error
	}
	return invoice, nil
}

func (is *InvoiceService) Create(invoice models.Invoice) (models.Invoice, error) {
	if result := is.DB.Create(&invoice); result.Error != nil {
		return is.Invoice, result.Error
	}
	return invoice, nil
}

func (is *InvoiceService) Update(invoice models.Invoice) (models.Invoice, error) {
	if result := is.DB.Select("*").Save(&invoice); result.Error != nil {
		return is.Invoice, result.Error
	}
	return invoice, nil
}

func (is *InvoiceService) Delete(invoice models.Invoice) error {
	if result := is.DB.Delete(&invoice); result.Error != nil {
		return result.Error
	}
	return nil
}
