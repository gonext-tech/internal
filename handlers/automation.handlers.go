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
			db.DB.Preload("User").Preload("Shop").Preload("Client").Where("date BETWEEN ? AND ?", now, nextHour).Find(&appointments)
		}
	}
	return c.JSON(200, &appointments)

}

func (ah *AutomationHandler) UpdateAppointment(c echo.Context) error {
	var appointment models.Appointment
	id := c.Param("id")
	for _, db := range ah.DB {
		if db.Name == "Appointment" {
			db.DB.Preload("").First(&appointment, id)
			appointment.UpdatedAt = time.Now()
		}
	}
	return c.JSON(200, appointment)
}
