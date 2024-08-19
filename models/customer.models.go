package models

import "time"

type Customer struct {
	ID          uint         `form:"id" gorm:"primaryKey"`
	Email       string       `form:"email" gorm:"unique"`
	Name        string       `form:"name"`
	Phone       string       `form:"phone"`
	Address     string       `form:"address"`
	Image       string       `form:"image"`
	Status      string       `form:"status" gorm:"default:ACTIVE"`
	Password    string       `form:"password"`
	Role        string       `form:"role" gorm:"default:USER"`
	ShopID      uint         `form:"shop_id"`
	Shop        CustomerShop `form:"shop" gorm:"foreignKey:ShopID"`
	ProjectName string       `form:"project_name"`
	CreatedAt   time.Time    `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `form:"updated_at" gorm:"autoUpdateTime"`
}

type CustomerShop struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	OwnerID     uint      `json:"owner_id"`
	Owner       User      `json:"owner" gorm:"foreignKey:OwnerID"`
	Address     string    `json:"address"`
	Image       string    `json:"image"`
	Workers     []User    `json:"workers" gorm:"foreignKey:ShopID"`
	Status      string    `json:"status" gorm:"default:ACTIVE"`
	ProjectName string    `form:"project_name"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
