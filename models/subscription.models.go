package models

import "time"

type Subscription struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	MembershipID    uint       `json:"membership_id" gorm:"foreignKey:MembershipID"`
	Membership      Membership `json:"membership"`
	ProjectName     string     `json:"project_name"`
	ShopID          uint       `json:"shop_id"`
	Shop            Shop       `json:"shop" gorm:"foreignKey:ShopID"`
	Notes           string     `json:"notes"`
	Status          string     `json:"status" gorm:"default:TOPAY"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	AutoRenewal     bool       `json:"auto_renewal"`
	PaymentMethod   string     `json:"payment_method"`
	NextBillingDate time.Time  `json:"next_billing_date"`
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	Discounts       string     `json:"discounts"`
	UserAddress     string     `json:"user_address"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type Membership struct {
	ID          uint      `json:"id"`
	ShopID      uint      `json:"shop_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Duration    int       `json:"duration"`
	Notes       string    `json:"notes"`
	MaxUsers    int       `json:"max_users"`
	ProjectName string    `json:"project_name" gorm:"default:Qwik"`
	Status      string    `json:"status" gorm:"default:ACTIVE"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
