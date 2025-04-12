package models

import (
	"time"
)

type Admin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `form:"email" gorm:"unique"`
	Name      string    `form:"name"`
	Image     string    `form:"image"`
	Phone     string    `form:"phone"`
	Address   string    `form:"address"`
	Status    string    `form:"status" gorm:"default:ACTIVE"`
	Password  string    `form:"-"`
	Role      string    `form:"role" gorm:"default:USER"`
	CreatedAt time.Time ` gorm:"autoCreateTime"`
	UpdatedAt time.Time ` gorm:"autoUpdateTime"`
}

type AdminBody struct {
	Email    string `form:"email" gorm:"unique"`
	Name     string `form:"name"`
	Image    string `form:"image"`
	Phone    string `form:"phone"`
	Address  string `form:"address"`
	Status   string `form:"status" gorm:"default:ACTIVE"`
	Password string `form:"password"`
	Role     string `form:"role" gorm:"default:USER"`
}

type PaginationResponse struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	LastPage    int `json:"last_page"`
	TotalCount  int `json:"total_count"`
}

type ParamResponse struct {
	Search  string
	Month   string
	Year    string
	Status  string
	SortBy  string
	OrderBy string
	Page    int
	Limit   int
}
