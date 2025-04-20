package models

import "time"

type Invoice struct {
	ID              uint       `form:"id" gorm:"primaryKey"`
	Description     string     `form:"description"`
	Amount          float64    `form:"amount"`
	AmountPaid      float64    `form:"amount_paid"`
	InvoiceRefID    *uint      `gorm:"index"`
	SubInvoices     []Invoice  `gorm:"foreignKey:InvoiceRefID"`
	Discount        float64    `form:"discount"`
	Total           float64    `form:"total"`
	ProjectID       uint       `form:"project_id"`
	Project         Project    `form:"project" gorm:"foreignKey:ProjectID"`
	ClientID        uint       `form:"client_id"`
	Client          Client     `form:"-" gorm:"foreignKey:ClientID"`
	CreatedByID     uint       `form:"created_by_id"`
	CreatedBy       Admin      `form:"created_by" gorm:"foreignKey:CreatedByID"`
	PaidToID        uint       `form:"paid_to_id"`
	PaidTo          Admin      `form:"paid_to" gorm:"foreignKey:PaidToID"`
	PaidAt          *time.Time `form:"paid_at"`
	CanceledByID    uint       `form:"canceled_by_id"`
	CanceledBy      Admin      `form:"canceled_by" gorm:"foreignKey:CanceledByID"`
	CanceledAt      *time.Time `form:"canceled_at" `
	InvoiceType     string     `form:"invoice_type"`
	Recurring       bool       `form:"recurring"`
	RecurringPeriod int        `form:"recurring_period" gorm:"default:0"`
	Category        string     `form:"category"`
	PaidFor         string     `form:"paid_for"`
	PaymentStatus   string     `form:"payment_status" gorm:"default:TOPAY"`
	Notes           string     `form:"notes"`
	IssueDate       time.Time  `form:"issue_date"`
	DueDate         time.Time  `form:"due_date"`
	Deleted         bool       `form:"deleted"`
	DeletedAt       *time.Time `form:"deleted_at" `
	CreatedAt       time.Time  `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `form:"updated_at" gorm:"autoUpdateTime"`
}
