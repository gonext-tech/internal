package models

import "time"

type Expense struct {
	ID          uint `gorm:"primaryKey"`
	Description string
	Amount      float64
	ProjectID   uint    `form:"project_id"`
	Project     Project `form:"project" gorm:"foreignKey:projectID"`
	ExpenseType string
	ExpenseDate time.Time
	Month       int
	Year        int
	CreatedAt   time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
