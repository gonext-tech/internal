package services

import (
	"time"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

func NewClientService(u models.Client, uStore *gorm.DB) *ClientServices {
	return &ClientServices{
		Client: u,
		DB:     uStore,
	}
}

type ClientServices struct {
	Client models.Client
	DB     *gorm.DB
}

func (cs *ClientServices) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Client, models.Meta, error) {
	var clients []models.Client
	query := cs.DB.Preload("Invoices").Preload("Projects")
	totalQuery := cs.DB.Preload("Invoices").Preload("Projects")

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
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&clients)
	totalRecords := int64(0)
	totalQuery.Model(&cs.Client).Count(&totalRecords)
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

	return clients, meta, nil
}

func (cs *ClientServices) GetID(id string) (models.Client, error) {
	var client models.Client
	result := cs.DB.Where("id = ?", id).First(&client)
	if result.Error != nil {
		return cs.Client, result.Error
	}
	return client, nil
}

func (cs *ClientServices) Create(client models.Client) (models.Client, error) {
	if result := cs.DB.Create(&client); result.Error != nil {
		return cs.Client, result.Error
	}
	return client, nil
}

func (cs *ClientServices) Update(client models.Client) (models.Client, error) {

	if result := cs.DB.Model(&client).Updates(client); result.Error != nil {
		return cs.Client, result.Error
	}
	return client, nil
}

func (cs *ClientServices) Delete(client models.Client) error {
	client.Deleted = true
	now := time.Now()
	client.DeletedAt = &now
	if result := cs.DB.Model(&client).Updates(client); result.Error != nil {
		return result.Error
	}
	return nil
}
