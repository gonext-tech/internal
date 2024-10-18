package models

import "time"

type Invoice struct {
	ID          uint      `form:"id" gorm:"primaryKey"`
	Description string    `form:"description"`
	Amount      float64   `form:"amount"`
	ProjectID   uint      `form:"project_id"`
	Project     Project   `form:"project" gorm:"foreignKey:projectID"`
	InvoiceType string    `form:"invoice_type"`
	InvoiceDate time.Time `form:"invoice_date"`
	Month       int       `form:"month"`
	Year        int       `form:"year"`
	CreatedAt   time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
