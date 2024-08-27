package models

import (
	"time"
)

type Appointment struct {
	ID                  uint                 `form:"id" gorm:"primaryKey"`
	ClientID            uint                 `form:"client_id"`
	Client              Client               `form:"client" gorm:"foreignKey:ClientID"`
	UserID              uint                 `form:"user_id"`
	User                User                 `form:"user" gorm:"foreignKey:UserID"`
	Date                time.Time            `form:"date"`
	ServiceID           uint                 `form:"service_id"`
	Service             Service              `form:"service" gorm:"foreignKey:ServiceID"`
	AppointmentServices []AppointmentService `form:"appointment_services" gorm:"foreignKey:AppointmentID"`
	PaymentStatus       string               `form:"payment_status" gorm:"default:TOPAY"`
	Priority            string               `form:"priority" gorm:"default:NO"`
	Price               float64              `form:"price"`
	Status              string               `form:"status" gorm:"default:PENDING"`
	Duration            int                  `form:"duration"`
	Notes               string               `form:"notes"`
	ShopID              uint                 `form:"shop_id"`
	Shop                Shop                 `form:"shop" gorm:"foreignKey:ShopID"`
	PaidAt              *time.Time           `form:"paid_at,omitempty"`
	NotificationSendAt  *time.Time           `form:"notification_send_at,omitempty"`
	CreatedAt           time.Time            `form:"created_at"`
	UpdatedAt           time.Time            `form:"update_at"`
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
