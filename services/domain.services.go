package services

import (
	"errors"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type DomainService struct {
	Domain models.Domain
	DB     *gorm.DB
}

func NewDomainService(p models.Domain, db *gorm.DB) *DomainService {
	return &DomainService{
		Domain: p,
		DB:     db,
	}
}

func (ss *DomainService) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Domain, models.Meta, error) {
	var domains []models.Domain
	query := ss.DB.Preload("Server")
	totalQuery := ss.DB.Preload("Server")
	if searchTerm != "" {
		searchTermWithWildcard := "%" + searchTerm + "%"
		query = query.Where("name LIKE ?", searchTermWithWildcard)
		totalQuery = query
	}

	if status != "" {
		query = query.Where("status = ?", status)
		totalQuery = totalQuery.Where("status = ?", status)
	}
	offset := (page - 1) * limit
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&domains)
	totalRecords := int64(0)
	totalQuery.Model(&ss.Domain).Count(&totalRecords)
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

	return domains, meta, nil
}

func (ds *DomainService) GetID(id string) (models.Domain, error) {
	var domain models.Domain
	if id == "0" || id == "" {
		return ds.Domain, errors.New("no id provided")
	}
	if result := ds.DB.First(&domain, id); result.Error != nil {
		return ds.Domain, result.Error
	}
	return domain, nil
}

func (ds *DomainService) Create(domain models.Domain) (models.Domain, error) {
	if result := ds.DB.Create(&domain); result.Error != nil {
		return ds.Domain, result.Error
	}
	return domain, nil
}

func (ds *DomainService) Update(domain models.Domain) (models.Domain, error) {
	if result := ds.DB.Model(&domain).Updates(domain); result.Error != nil {
		return ds.Domain, result.Error
	}
	return domain, nil
}

func (ds *DomainService) Delete(domain models.Domain) error {
	if result := ds.DB.Delete(&domain); result.Error != nil {
		return result.Error
	}
	return nil
}
