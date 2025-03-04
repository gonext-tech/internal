package models

import "time"

type User struct {
	ID          uint      `form:"id" gorm:"primaryKey"`
	Email       string    `form:"email" gorm:"unique"`
	Name        string    `form:"name"`
	Phone       string    `form:"phone"`
	Address     string    `form:"address"`
	Image       string    `form:"image"`
	Status      string    `form:"status" gorm:"default:ACTIVE"`
	Password    string    `form:"password"`
	Role        string    `form:"role" gorm:"default:USER"`
	ShopID      uint      `form:"shop_id"`
	Shop        *Shop     `form:"shop" gorm:"foreignKey:ShopID"`
	ProjectName string    `form:"project_name"`
	CreatedAt   time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
