package services

import (
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"gorm.io/gorm"
)

func NewClientService(uStore *gorm.DB) *ClientServices {
	return &ClientServices{
		DB: uStore,
	}
}

type ClientServices struct {
	DB *gorm.DB
}

func (cs *ClientServices) GetALL(query queries.InvoiceQueryParams) ([]models.Client, models.Meta, error) {
	var clients []models.Client
	var totalCount int64

	dbQuery := cs.buildInvoicesURL(query)
	dbQuery.Session(&gorm.Session{}).Model(&models.Invoice{}).Count(&totalCount)
	offset := (query.Page - 1) * query.Limit
	dbQuery.Order(query.SortBy + " " + query.OrderBy).Offset(offset).Limit(query.Limit).Find(&clients)
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

	return clients, meta, nil
}

func (cs *ClientServices) GetID(id string) (*models.Client, error) {
	var client models.Client
	result := cs.DB.Where("id = ?", id).First(&client)
	if result.Error != nil {
		return &models.Client{}, result.Error
	}
	return &client, nil
}

func (cs *ClientServices) Create(client *models.Client) error {
	return cs.DB.Create(&client).Error
}
func (cs *ClientServices) Update(client *models.Client) error {
	return cs.DB.Model(&client).Updates(client).Error
}

func (cs *ClientServices) Delete(client *models.Client) error {
	client.Deleted = true
	now := time.Now()
	client.DeletedAt = &now
	return cs.DB.Model(&client).Updates(client).Error
}

func (cs *ClientServices) buildInvoicesURL(query queries.InvoiceQueryParams) *gorm.DB {
	dbQuery := cs.DB.
		Preload("Projects").
		Preload("Invoices")

	if query.SearchTerm != "" {
		term := "%" + query.SearchTerm + "%"
		dbQuery = dbQuery.Where("name LIKE ?", term)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.InvoiceType != "" {
		dbQuery = dbQuery.Where("invoice_type = ?", query.InvoiceType)
	}
	return dbQuery

}
