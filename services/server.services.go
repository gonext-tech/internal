package services

import (
	"errors"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type ServerService struct {
	Server models.MonitoredServer
	DB     *gorm.DB
}

func NewServerService(p models.MonitoredServer, db *gorm.DB) *ServerService {
	return &ServerService{
		Server: p,
		DB:     db,
	}
}

func (ss *ServerService) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.MonitoredServer, models.Meta, error) {
	var servers []models.MonitoredServer
	query := ss.DB
	totalQuery := query
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
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&servers)
	totalRecords := int64(0)
	totalQuery.Model(&ss.Server).Count(&totalRecords)
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

	return servers, meta, nil
}

func (ss *ServerService) GetID(id string) (models.MonitoredServer, error) {
	var server models.MonitoredServer
	if id == "0" || id == "" {
		return ss.Server, errors.New("no id provided")
	}
	if result := ss.DB.First(&server, id); result.Error != nil {
		return ss.Server, result.Error
	}
	return server, nil
}

func (ss *ServerService) Create(server models.MonitoredServer) (models.MonitoredServer, error) {
	if result := ss.DB.Create(&server); result.Error != nil {
		return ss.Server, result.Error
	}
	return server, nil
}

func (ss *ServerService) Update(server models.MonitoredServer) (models.MonitoredServer, error) {
	if result := ss.DB.Model(&server).Updates(server); result.Error != nil {
		return ss.Server, result.Error
	}
	return server, nil
}

func (ss *ServerService) Delete(server models.MonitoredServer) error {
	if result := ss.DB.Delete(&server); result.Error != nil {
		return result.Error
	}
	return nil
}
