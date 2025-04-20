package services

import (
	"errors"
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"gorm.io/gorm"
)

type ServerService struct {
	DB *gorm.DB
}

func NewServerService(p models.MonitoredServer, db *gorm.DB) *ServerService {
	return &ServerService{
		DB: db,
	}
}

func (ss *ServerService) GetALL(query queries.InvoiceQueryParams) ([]models.MonitoredServer, models.Meta, error) {
	var servers []models.MonitoredServer
	var totalCount int64

	dbQuery := ss.buildInvoicesURL(query)
	dbQuery.Session(&gorm.Session{}).Model(&models.Invoice{}).Count(&totalCount)
	offset := (query.Page - 1) * query.Limit
	dbQuery.Order(query.SortBy + " " + query.OrderBy).Offset(offset).Limit(query.Limit).Find(&servers)
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

	return servers, meta, nil
}

func (ss *ServerService) GetID(id string) (*models.MonitoredServer, error) {
	var server models.MonitoredServer
	if id == "0" || id == "" {
		return nil, errors.New("no id provided")
	}
	if result := ss.DB.First(&server, id); result.Error != nil {
		return nil, result.Error
	}
	return &server, nil
}

func (ss *ServerService) Create(server *models.MonitoredServer) error {
	return ss.DB.Create(&server).Error
}

func (ss *ServerService) Update(server *models.MonitoredServer) error {
	return ss.DB.Model(&server).Updates(server).Error
}

func (ss *ServerService) Delete(server *models.MonitoredServer) error {
	now := time.Now()
	server.DeletedAt = &now
	server.Deleted = true
	return ss.DB.Model(&server).Updates(server).Error
}

func (ss *ServerService) buildInvoicesURL(query queries.InvoiceQueryParams) *gorm.DB {
	dbQuery := ss.DB

	if query.SearchTerm != "" {
		term := "%" + query.SearchTerm + "%"
		dbQuery = dbQuery.Where("name LIKE ?", term)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	return dbQuery

}
