package services

import (
	"errors"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type DonorService struct {
	Donor models.Donor
	DB    *gorm.DB
}

func NewDonorService(d models.Donor, db *gorm.DB) *DonorService {
	return &DonorService{
		Donor: d,
		DB:    db,
	}
}

func (ds *DonorService) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Donor, models.Meta, error) {
	var donors []models.Donor
	query := ds.DB.Preload("BloodType").Preload("City")
	totalQuery := ds.DB.Preload("BloodType").Preload("City")
	if searchTerm != "" {
		searchTermWithWildcard := "%" + searchTerm + "%"
		query = query.Where("name LIKE ? OR phone LIKE ?", searchTermWithWildcard, searchTermWithWildcard)
		totalQuery = query
	}

	if status != "" {
		query = query.Where("status = ?", status)
		totalQuery = totalQuery.Where("status = ?", status)
	}
	offset := (page - 1) * limit
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&donors)
	totalRecords := int64(0)
	totalQuery.Model(&ds.Donor).Count(&totalRecords)
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

	return donors, meta, nil
}

func (ds *DonorService) GetID(id string) (models.Donor, error) {
	var donor models.Donor
	if id == "0" || id == "" {
		return ds.Donor, errors.New("no id provided")
	}
	if result := ds.DB.Preload("City").Preload("BloodType").First(&donor, id); result.Error != nil {
		return ds.Donor, result.Error
	}
	return donor, nil
}

func (ds *DonorService) Create(donor models.Donor) (models.Donor, error) {
	if result := ds.DB.Create(&donor); result.Error != nil {
		return ds.Donor, result.Error
	}
	return donor, nil
}

func (ds *DonorService) Update(donor models.Donor) (models.Donor, error) {
	if result := ds.DB.Model(&donor).Updates(donor); result.Error != nil {
		return ds.Donor, result.Error
	}
	return donor, nil
}

func (ds *DonorService) Delete(donor models.Donor) error {
	if result := ds.DB.Delete(&donor); result.Error != nil {
		return result.Error
	}
	return nil
}
