package models

import "time"

type Client struct {
	ID        uint       `form:"id" gorm:"primaryKey"`
	Email     string     `form:"email" gorm:"unique"`
	Name      string     `form:"name"`
	Phone     string     `form:"phone"`
	Address   string     `form:"address"`
	Image     string     `form:"image"`
	Status    string     `form:"status" gorm:"default:ACTIVE"`
	Projects  []Project  `form:"-" gorm:"foreignKey:ClientID"`
	Invoices  []Invoice  `form:"-" gorm:"foreignKey:ClientID"`
	Deleted   bool       `form:"-" `
	CreatedAt time.Time  `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `form:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `form:"-" `
}
