package models

import (
	"time"

	"gorm.io/gorm"
)

type Stats struct {
	ID                 uint      `form:"id" gorm:"primaryKey"`
	ShopID             uint      `form:"shop_id"`
	Month              int       `form:"month"`
	Year               int       `form:"year"`
	TotalSubscriptions int       `form:"total_subscriptions"`
	TotalRevenue       float64   `form:"total_revenue"`
	TotalExpenses      float64   `form:"total_expenses"`
	NewSubscriptions   int       `form:"new_subscriptions"`
	NetProfit          float64   `form:"net_profit"`
	CreatedAt          time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeSave GORM hook to calculate NetProfit before saving
func (fr *Stats) BeforeSave(tx *gorm.DB) (err error) {
	fr.NetProfit = fr.TotalRevenue - fr.TotalExpenses
	return nil
}
