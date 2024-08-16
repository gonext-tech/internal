package models

import (
	"time"
)

type User struct {
	ID        uint          `json:"id" gorm:"primaryKey"`
	Email     string        `json:"username" gorm:"unique"`
	Name      string        `json:"name"`
	Phone     string        `json:"phone"`
	Address   string        `json:"address"`
	Active    string        `json:"active" gorm:"default:ACTIVE"`
	Password  string        `json:"-"`
	Portfolio string        `json:"portfolio"`
	Role      string        `json:"role" gorm:"default:'USER'"`
	Skills    *[]UserSkills `json:"skills" gorm:"foreignKey:UserID"`
	CreatedAt time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type UserSkills struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	TagID  uint `json:"tag_id"`
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
	Status  string
	SortBy  string
	OrderBy string
	Page    int
	Limit   int
}
