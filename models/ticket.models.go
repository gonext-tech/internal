package models

import "time"

type Ticket struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `form:"email" json:"email"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Status      string `form:"status" json:"status" gorm:"default:OPEN"` // e.g., Open, In Progress, Closed
	Priority    string `form:"priority" json:"priority" gorm:"default:M"`
	Type        string `form:"type" json:"type"`
	// Add other fields as needed
	UserID    uint      `form:"user_id" json:"user_id"`
	ProjectID uint      `form:"project_id" json:"project_id"`
	User      User      `gorm:"foreignKey:UserID"`
	Project   Project   `json:"project" form:"project" gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Notification struct {
	ID        uint `gorm:"primaryKey"`
	TicketID  uint
	Ticket    Ticket `gorm:"foreignKey:TicketID"`
	Title     string
	Message   string
	New       bool
	CreatedAt time.Time
}
