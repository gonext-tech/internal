package services

import (
	"errors"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"gorm.io/gorm"
)

type DomainService struct {
	DB *gorm.DB
}

func NewDomainService(db *gorm.DB) *DomainService {
	return &DomainService{
		DB: db,
	}
}

func (ss *DomainService) GetALL(query queries.InvoiceQueryParams) ([]models.Domain, models.Meta, error) {
	var domains []models.Domain
	var totalCount int64

	dbQuery := ss.buildInvoicesURL(query)
	dbQuery.Session(&gorm.Session{}).Model(&models.Invoice{}).Count(&totalCount)
	offset := (query.Page - 1) * query.Limit
	dbQuery.Order(query.SortBy + " " + query.OrderBy).Offset(offset).Limit(query.Limit).Find(&domains)
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

	return domains, meta, nil
}

func (ds *DomainService) GetID(id string) (*models.Domain, error) {
	var domain models.Domain
	if id == "0" || id == "" {
		return nil, errors.New("no id provided")
	}
	if result := ds.DB.First(&domain, id); result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}

func (ds *DomainService) Create(domain *models.Domain) error {
	return ds.DB.Create(&domain).Error

}

func (ds *DomainService) Update(domain *models.Domain) error {
	return ds.DB.Model(&domain).Updates(domain).Error

}

func (ds *DomainService) Delete(domain *models.Domain) error {
	return ds.DB.Delete(&domain).Error

}

func (ds *DomainService) buildInvoicesURL(query queries.InvoiceQueryParams) *gorm.DB {
	dbQuery := ds.DB.Preload("Server")

	if query.SearchTerm != "" {
		term := "%" + query.SearchTerm + "%"
		dbQuery = dbQuery.Where("name LIKE ?", term)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	return dbQuery

}
