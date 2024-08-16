package models

import "time"

type Customer struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Email       string       `json:"email" gorm:"unique"`
	Name        string       `json:"name"`
	Phone       string       `json:"phone"`
	Address     string       `json:"address"`
	Image       string       `json:"image"`
	Status      string       `json:"status" gorm:"default:ACTIVE"`
	Password    string       `json:"password"`
	Role        string       `json:"role" gorm:"default:USER"`
	ShopID      uint         `json:"shop_id"`
	Shop        CustomerShop `json:"shop" gorm:"foreignKey:ShopID"`
	ProjectName string       `json:"project_name"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

type CustomerShop struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	OwnerID   uint      `json:"owner_id"`
	Owner     User      `json:"owner" gorm:"foreignKey:OwnerID"`
	Address   string    `json:"address"`
	Image     string    `json:"image"`
	Workers   []User    `json:"workers" gorm:"foreignKey:ShopID"`
	Status    string    `json:"status" gorm:"default:ACTIVE"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
