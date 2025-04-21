package services

import (
	"errors"
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"gorm.io/gorm"
)

type ProjectService struct {
	DB *gorm.DB
}

func NewProjectService(p models.Project, db *gorm.DB) *ProjectService {
	return &ProjectService{
		DB: db,
	}
}

func (ps *ProjectService) GetALL(query queries.InvoiceQueryParams) ([]models.Project, models.Meta, error) {
	var projects []models.Project
	var totalCount int64

	dbQuery := ps.buildInvoicesURL(query)
	dbQuery.Session(&gorm.Session{}).Model(&models.Invoice{}).Count(&totalCount)
	offset := (query.Page - 1) * query.Limit
	dbQuery.Order(query.SortBy + " " + query.OrderBy).Offset(offset).Limit(query.Limit).Find(&projects)
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
	return projects, meta, nil
}

func (ps *ProjectService) GetID(id string) (*models.Project, error) {
	var project models.Project
	if id == "0" || id == "" {
		return &models.Project{}, errors.New("no id provided")
	}
	if result := ps.DB.Preload("Server").Preload("Lead").Preload("Client").First(&project, id); result.Error != nil {
		return &models.Project{}, result.Error
	}
	return &project, nil
}

func (ps *ProjectService) Create(project *models.Project) error {
	return ps.DB.Create(project).Error

}

func (ps *ProjectService) Update(project *models.Project) error {
	return ps.DB.Model(&models.Project{}).
		Where("id = ?", project.ID).
		Select("*").
		Updates(project).Error
}

func (ps *ProjectService) Delete(project *models.Project) error {
	now := time.Now()
	project.DeletedAt = &now
	project.Deleted = true
	return ps.DB.Model(&project).Updates(project).Error

}

func (ps *ProjectService) buildInvoicesURL(query queries.InvoiceQueryParams) *gorm.DB {
	dbQuery := ps.DB.Preload("Server").Preload("Lead").Preload("Client")

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
