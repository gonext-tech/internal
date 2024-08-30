package models

import (
	"time"
)

type Appointment struct {
	ID                  uint                 `json:"id" form:"id" gorm:"primaryKey"`
	ClientID            uint                 `json:"client_id" form:"client_id"`
	Client              Client               `json:"client" form:"client" gorm:"foreignKey:ClientID"`
	UserID              uint                 `json:"user_id" form:"user_id"`
	User                User                 `json:"user" form:"user" gorm:"foreignKey:UserID"`
	Date                time.Time            `json:"date" form:"date"`
	ServiceID           uint                 `json:"service_id" form:"service_id"`
	Service             Service              `json:"service" form:"service" gorm:"foreignKey:ServiceID"`
	AppointmentServices []AppointmentService `json:"appointment_services" form:"appointment_services" gorm:"foreignKey:AppointmentID"`
	PaymentStatus       string               `json:"payment_status" form:"payment_status" gorm:"default:TOPAY"`
	Priority            string               `json:"priority" form:"priority" gorm:"default:NO"`
	Price               float64              `json:"price" form:"price"`
	Status              string               `json:"status" form:"status" gorm:"default:PENDING"`
	Duration            int                  `json:"duration" form:"duration"`
	Notes               string               `json:"notes" form:"notes"`
	ShopID              uint                 `json:"shop_id" form:"shop_id"`
	Shop                Shop                 `json:"shop" form:"shop" gorm:"foreignKey:ShopID"`
	PaidAt              *time.Time           `json:"paid_at" form:"paid_at,omitempty"`
	NotificationSendAt  *time.Time           `json:"notification_send_at" form:"notification_send_at,omitempty"`
	CreatedAt           time.Time            `json:"created_at" form:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at" form:"updated_at"`
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
