package handlers

import (
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/labstack/echo/v4"
)

type AutomationHandler struct {
	DB []models.ProjectsDB
}

func NewAutomationHandler(db []models.ProjectsDB) *AutomationHandler {
	return &AutomationHandler{
		DB: db,
	}
}

func (ah *AutomationHandler) GetAppointments(c echo.Context) error {
	var appointments []models.Appointment
	now := time.Now()
	nextHour := now.Add(1 * time.Hour)
	for _, db := range ah.DB {
		if db.Name == "Qwik" {
			db.DB.Preload("User").Preload("Shop").Preload("Client").
				Joins("JOIN shops ON shops.id = appointments.shop_id").
				Where("appointments.date BETWEEN ? AND ?", now, nextHour).
				Where("appointments.status = ?", "PENDING").
				Where("appointments.notification_send_at IS NULL").
				Where("shops.send_wp = ? AND shops.status = ?", true, "ACTIVE").
				Find(&appointments)
		}
	}
	return c.JSON(200, &appointments)

}

func (ah *AutomationHandler) UpdateAppointment(c echo.Context) error {
	var appointment models.Appointment
	id := c.Param("id")
	for _, db := range ah.DB {
		if db.Name == "Qwik" {
			db.DB.First(&appointment, id)
			now := time.Now()
			appointment.NotificationSendAt = &now
			db.DB.Model(&appointment).Updates(appointment)
		}
	}
	return c.JSON(200, appointment)
}
