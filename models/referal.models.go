package models

import "time"

type Referal struct {
	ID               uint       `form:"id" gorm:"primaryKey"`
	Name             string     `form:"name"`
	Email            string     `form:"email" gorm:"unique"`
	Phone            string     `form:"phone"`
	Image            string     `form:"image"`
	Password         string     `form:"-"`
	TotalUsers       int        `form:"total_users"`
	TotalRevenue     float64    `form:"total_revenue"`
	Status           string     `form:"status"`
	RemainingRevenue float64    `form:"remaining_revenue"`
	LastWithDrawAt   *time.Time `form:"last_withdraw_at"`
	CreatedAt        time.Time  `form:"created_at"`
	UpdatedAt        time.Time  `form:"updated_at"`
}
