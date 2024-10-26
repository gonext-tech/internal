package models

import "time"

type Invoice struct {
	ID          uint      `form:"id" gorm:"primaryKey"`
	ShopID      uint      `form:"shop_id"`
	Description string    `form:"description"`
	Amount      float64   `form:"amount"`
	ProjectID   uint      `form:"project_id"`
	Project     Project   `form:"project" gorm:"foreignKey:ProjectID"`
	UserID      uint      `form:"user_id"`
	User        User      `form:"user" gorm:"foreignKey:UserID"`
	InvoiceType string    `form:"invoice_type"`
	Recurring   bool      `form:"recurring"`
	Period      int       `form:"period"`
	InvoiceDate time.Time `form:"invoice_date"`
	CreatedAt   time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
