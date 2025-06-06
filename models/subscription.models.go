package models

import "time"

type Subscription struct {
	ID            uint       `form:"id" gorm:"primaryKey"`
	MembershipID  uint       `form:"membership_id" gorm:"foreignKey:MembershipID"`
	Membership    Membership `form:"membership"`
	ProjectName   string     `form:"project_name"`
	ShopID        uint       `form:"shop_id"`
	Shop          Shop       `form:"shop" gorm:"foreignKey:ShopID"`
	Notes         string     `form:"notes"`
	PaymentStatus string     `form:"payment_status" gorm:"default:TOPAY"`
	Status        string     `form:"status" gorm:"default:ACTIVE"`
	StartDate     time.Time  `form:"start_date"`
	EndDate       time.Time  `form:"end_date"`
	AutoRenewal   bool       `form:"auto_renewal"`
	PaymentMethod string     `form:"payment_method"`
	Tag           string     `form:"tag"`
	Amount        float64    `form:"amount"`
	Currency      string     `form:"currency"`
	Discounts     string     `form:"discounts"`
	UserAddress   string     `form:"user_address"`
	PaidAt        *time.Time `form:"paid_at,omitempty"`
	CreatedAt     time.Time  `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `form:"updated_at" gorm:"autoUpdateTime"`
}

type Membership struct {
	ID          uint      `form:"id"`
	Name        string    `form:"name"`
	Price       float64   `form:"price"`
	Duration    int       `form:"duration"`
	Notes       string    `form:"notes"`
	MaxUsers    int       `form:"max_users"`
	ProjectName string    `form:"project_name"`
	Status      string    `form:"status" gorm:"default:ACTIVE"`
	CreatedAt   time.Time `form:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `form:"updated_at" gorm:"autoUpdateTime"`
}
