package services

import (
	"errors"
	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type ReferalService struct {
	Referal models.Referal
	DB      *gorm.DB
}

func NewReferalService(r models.Referal, db *gorm.DB) *ReferalService {
	return &ReferalService{
		Referal: r,
		DB:      db,
	}
}

func (rs *ReferalService) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Referal, models.Meta, error) {
	var referals []models.Referal
	query := rs.DB
	totalQuery := rs.DB
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
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&referals)
	totalRecords := int64(0)
	totalQuery.Model(&rs.Referal).Count(&totalRecords)
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

	return referals, meta, nil
}

func (rs *ReferalService) GetID(id string) (models.Referal, error) {
	var referal models.Referal
	if id == "0" || id == "" {
		return models.Referal{}, errors.New("no id provided")
	}
	if result := rs.DB.First(&referal, id); result.Error != nil {
		return models.Referal{}, result.Error
	}
	return referal, nil
}

func (rs *ReferalService) Create(referal models.Referal) (models.Referal, error) {
	if result := rs.DB.Create(&referal); result.Error != nil {
		return models.Referal{}, result.Error
	}
	return referal, nil
}

func (rs *ReferalService) Update(referal models.Referal) (models.Referal, error) {
	if result := rs.DB.Model(&referal).Updates(referal); result.Error != nil {
		return rs.Referal, result.Error
	}
	return referal, nil
}

func (rs *ReferalService) Delete(referal models.Referal) error {
	if result := rs.DB.Delete(&referal); result.Error != nil {
		return result.Error
	}
	return nil
}
