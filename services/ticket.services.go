package services

import (
	"errors"

	"github.com/ramyjaber1/internal/models"
	"gorm.io/gorm"
)

type TickerServices struct {
	Ticket models.Ticket
	DB     *gorm.DB
}

func NewTicketService(t models.Ticket, db *gorm.DB) *TickerServices {
	return &TickerServices{
		Ticket: t,
		DB:     db,
	}
}

func (ts *TickerServices) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.Ticket, models.Meta, error) {
	var tickets []models.Ticket
	query := ts.DB.Preload("Project").Preload("User")
	totalQuery := query

	if searchTerm != "" {
		searchTermWithWildcard := "%" + searchTerm + "%"
		query = query.Where("title LIKE ?", searchTermWithWildcard)
		totalQuery = query
	}

	if status != "" {
		query = query.Where("status = ?", status)
		totalQuery = totalQuery.Where("status = ?", status)
	}
	offset := (page - 1) * limit
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&tickets)
	totalRecords := int64(0)
	totalQuery.Model(&ts.Ticket).Count(&totalRecords)
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

	return tickets, meta, nil
}

func (ts *TickerServices) GetID(id string) (models.Ticket, error) {
	var ticket models.Ticket
	if id == "0" || id == "" {
		return models.Ticket{}, errors.New("no id provided")
	}
	if result := ts.DB.Preload("User").Preload("Project").First(&ticket, id); result.Error != nil {
		return models.Ticket{}, result.Error
	}
	return ticket, nil
}

func (ts *TickerServices) Create(ticket models.Ticket) (models.Ticket, error) {
	if result := ts.DB.Create(&ticket); result.Error != nil {
		return models.Ticket{}, result.Error
	}
	if result := ts.DB.Preload("Project").Find(&ticket); result.Error != nil {
		return models.Ticket{}, result.Error
	}
	return ticket, nil
}

func (ts *TickerServices) Update(ticket models.Ticket) (models.Ticket, error) {
	if result := ts.DB.Model(&ticket).Updates(ticket); result.Error != nil {
		return models.Ticket{}, result.Error
	}
	return ticket, nil
}

func (ts *TickerServices) Delete(ticket models.Ticket) error {
	if result := ts.DB.Delete(&ticket); result.Error != nil {
		return result.Error
	}
	return nil
}
