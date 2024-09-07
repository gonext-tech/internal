package models

import (
	"time"
)

type Client struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name"`
	Phone        string        `json:"phone"`
	CountryCode        string        `json:"country_code"`
	Address      string        `json:"address"`
	Image        string        `json:"image"`
	ShopID       uint          `json:"shop_id"`
	Appointments []Appointment `json:"appointments" gorm:"foreignKey:ClientID"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"update_at"`
}
