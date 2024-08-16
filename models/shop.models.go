package models

import "time"

type Shop struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	OwnerID           uint      `json:"owner_id"`
	Owner             User      `json:"owner" gorm:"foreignKey:OwnerID"`
	Address           string    `json:"address"`
	Image             string    `json:"image"`
	Workers           []User    `json:"workers" gorm:"foreignKey:ShopID"`
	Status            string    `json:"status" gorm:"default:ACTIVE"`
	TotalIncome       float64   `json:"total_income"`
	TotalExpenses     float64   `json:"total_expenses"`
	TotalClients      int       `json:"total_clients"`
	TotalAppointments int       `json:"total_appointments"`
	ProjectName       string    `json:"project_name"`
	TopClients        []User    `json:"top_clients"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
