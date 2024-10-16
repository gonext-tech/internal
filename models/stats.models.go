package models

import "time"

type Stats struct {
	ID                 uint `gorm:"primaryKey"`
	Month              int
	Year               int
	TotalSubscriptions int
	TotalRevenue       float64
	TotalExpenses      float64
	NewSubscriptions   int
	NetProfit          float64
	CreatedAt          time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
