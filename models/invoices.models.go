package models

import "time"

type Invoice struct {
	ID              uint      `form:"id" gorm:"primaryKey"`
	Description     string    `form:"description"`
	Amount          float64   `form:"amount"`
	Discount        float64   `form:"discount"`
	Total           float64   `form:"total"`
	ProjectID       uint      `form:"project_id"`
	Project         Project   `form:"project" gorm:"foreignKey:ProjectID"`
	ClientID        uint      `form:"project_id"`
	Client          Client    `form:"-" gorm:"foreignKey:ClientID"`
	CreatedByID     uint      `form:"created_by_id"`
	CreatedBy       Admin     `form:"created_by" gorm:"foreignKey:CreatedByID"`
	PaidToID        uint      `form:"paid_to_id"`
	PaidTo          Admin     `form:"paid_to" gorm:"foreignKey:PaidToID"`
	CanceledByID    uint      `form:"canceled_by_id"`
	CanceledBy      Admin     `form:"canceled_by" gorm:"foreignKey:CanceledByID"`
	InvoiceType     string    `form:"invoice_type"`
	Recurring       bool      `form:"recurring"`
	RecurringPeriod int       `form:"recurring_period" gorm:"default:0"`
	PaymentStatus   string    `form:"payment_status" gorm:"default:TOPAY"`
	Notes           string    `form:"notes"`
	InvoiceDate     time.Time `form:"invoice_date"`
	CreatedAt       time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
