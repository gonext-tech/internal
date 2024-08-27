package services

import (
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/utils"
)

type AppointmentServices struct {
	STORES      []models.ProjectsDB
	Appointment models.Appointment
}

func NewAppointmentService(db []models.ProjectsDB) *AppointmentServices {
	return &AppointmentServices{
		STORES: db,
	}
}

func (as *AppointmentServices) GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) ([]models.Appointment, models.Meta, error) {
	DB := utils.GetCurrentDB("Qwik", as.STORES)
	var appointments []models.Appointment
	query := DB
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
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&appointments)
	totalRecords := int64(0)
	totalQuery.Model(&as.Appointment).Count(&totalRecords)
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

	return appointments, meta, nil
}

func (as *AppointmentServices) GetID(id, dbName string) (models.Appointment, error) {
	var appointment models.Appointment
	DB := utils.GetCurrentDB("Qwik", as.STORES)
	if result := DB.Table("appointments").Preload("Client").First(&appointment, id); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (as *AppointmentServices) Create(appointment models.Appointment) (models.Appointment, error) {
	DB := utils.GetCurrentDB("Qwik", as.STORES)
	if result := DB.Table("appointments").Create(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}
func (as *AppointmentServices) Update(appointment models.Appointment) (models.Appointment, error) {
	DB := utils.GetCurrentDB("Qwik", as.STORES)
	if result := DB.Table("appointments").Updates(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (ss *SubscriptionServices) Delete(appointment models.Appointment) (models.Appointment, error) {
	DB := utils.GetCurrentDB("Qwik", ss.STORES)
	if result := DB.Table("appointments").Delete(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}
