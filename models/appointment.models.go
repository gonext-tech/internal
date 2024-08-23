package models

import (
	"time"
)

type Appointment struct {
	ID                  uint                 `json:"id" gorm:"primaryKey"`
	ClientID            uint                 `json:"client_id"`
	Client              Client               `json:"client" gorm:"foreignKey:ClientID"`
	UserID              uint                 `json:"user_id"`
	User                User                 `json:"user" gorm:"foreignKey:UserID"`
	Date                time.Time            `json:"date"`
	ServiceID           uint                 `json:"service_id"`
	Service             Service              `json:"service" gorm:"foreignKey:ServiceID"`
	AppointmentServices []AppointmentService `json:"appointment_services" gorm:"foreignKey:AppointmentID"`
	PaymentStatus       string               `json:"payment_status" gorm:"default:TOPAY"`
	Priority            string               `json:"priority" gorm:"default:NO"`
	Price               float64              `json:"price"`
	Status              string               `json:"status" gorm:"default:PENDING"`
	Duration            int                  `json:"duration"`
	Notes               string               `json:"notes"`
	ShopID              uint                 `json:"shop_id"`
	Shop                Shop                 `json:"shop" gorm:"foreignKey:ShopID"`
	PaidAt              *time.Time           `json:"paid_at,omitempty"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"update_at"`
}

type AppointmentService struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	AppointmentID uint    `json:"appointment_id"`
	ServiceID     uint    `json:"service_id"`
	Service       Service `json:"service" gorm:"foreignKey:ServiceID"`
}

type Service struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Color    string  `json:"color"`
	ShopID   uint    `json:"shop_id"`
	Duration int     `json:"duration"`
}
