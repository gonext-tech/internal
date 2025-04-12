package services

import (
	"errors"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type ProjectService struct {
	Project models.Project
	DB      *gorm.DB
}

func NewProjectService(p models.Project, db *gorm.DB) *ProjectService {
	return &ProjectService{
		Project: p,
		DB:      db,
	}
}

func (ps *ProjectService) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Project, models.Meta, error) {
	var projects []models.Project
	query := ps.DB.Preload("Server").Preload("Lead")
	totalQuery := ps.DB.Preload("Server").Preload("Lead")

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
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&projects)
	totalRecords := int64(0)
	totalQuery.Model(&ps.Project).Count(&totalRecords)
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

	return projects, meta, nil
}

func (ps *ProjectService) GetID(id string) (models.Project, error) {
	var project models.Project
	if id == "0" || id == "" {
		return models.Project{}, errors.New("no id provided")
	}
	if result := ps.DB.Preload("Server").Preload("Lead").First(&project, id); result.Error != nil {
		return models.Project{}, result.Error
	}
	return project, nil
}

func (ps *ProjectService) Create(project *models.Project) (*models.Project, error) {
	if result := ps.DB.Create(project); result.Error != nil {
		return &ps.Project, result.Error
	}
	return project, nil
}

func (ps *ProjectService) Update(project models.Project) (models.Project, error) {
	if result := ps.DB.Model(&project).Updates(project); result.Error != nil {
		return models.Project{}, result.Error
	}
	return project, nil
}

func (ps *ProjectService) Delete(project models.Project) error {
	if result := ps.DB.Delete(&project); result.Error != nil {
		return result.Error
	}
	return nil
}
