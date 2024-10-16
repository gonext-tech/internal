package models

import "time"

type Shop struct {
	ID                uint           `json:"id" form:"id" gorm:"primaryKey"`
	Name              string         `json:"name" form:"name"`
	OwnerID           uint           `json:"owner_id" form:"owner_id"`
	Owner             User           `json:"owner" form:"owner" gorm:"foreignKey:OwnerID"`
	Address           string         `json:"address" form:"address"`
	Image             string         `json:"image" form:"image"`
	Workers           []User         `json:"workers" form:"workers" gorm:"foreignKey:ShopID"`
	Subscriptions     []Subscription `json:"subscriptions" form:"subscriptions" gorm:"foreignKey:ShopID"`
	NextBillingDate   *time.Time     `json:"next_billing_date" form:"next_billing_date,omitempty"`
	SendWP            bool           `json:"send_wp" form:"send_wp"`
	WPMessage         string         `json:"wp_message" form:"wp_message"`
	Category          string         `json:"category"`
	Status            string         `json:"status" form:"status" gorm:"default:ACTIVE"`
	TotalIncome       float64        `json:"total_income" form:"total_income"`
	TotalExpenses     float64        `json:"total_expenses" form:"total_expenses"`
	TotalClients      int            `json:"total_clients" form:"total_clients"`
	TotalAppointments int            `json:"total_appointments" form:"total_appointments"`
	ProjectName       string         `json:"project_name" form:"project_name"`
	TopClients        []User         `json:"top_clients" form:"top_clients"`
	CreatedAt         time.Time      `json:"created_at" form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" form:"updated_at" gorm:"autoUpdateTime"`
}
