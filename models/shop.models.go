package models

import "time"

type Shop struct {
	ID                uint      `form:"id" gorm:"primaryKey"`
	Name              string    `form:"name"`
	OwnerID           uint      `form:"owner_id"`
	Owner             User      `form:"owner" gorm:"foreignKey:OwnerID"`
	Address           string    `form:"address"`
	Image             string    `form:"image"`
	Workers           []User    `form:"workers" gorm:"foreignKey:ShopID"`
	Status            string    `form:"status" gorm:"default:ACTIVE"`
	TotalIncome       float64   `form:"total_income"`
	TotalExpenses     float64   `form:"total_expenses"`
	TotalClients      int       `form:"total_clients"`
	TotalAppointments int       `form:"total_appointments"`
	ProjectName       string    `form:"project_name"`
	TopClients        []User    `form:"top_clients"`
	CreatedAt         time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
