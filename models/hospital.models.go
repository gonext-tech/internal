package models

import "time"

type Hospital struct {
	ID         uint      `gorm:"primaryKey" form:"id"`
	Name       string    `gorm:"type:varchar(100);not null" form:"name"`
	Address    string    `gorm:"type:varchar(255);not null" form:"address"`
	Status     string    `gorm:"type:varchar(100);not null" form:"status"`
	CityID     uint      `form:"city_id"`
	City       City      `form:"city" gorm:"foreignKey:CityID"`
	State      string    `gorm:"type:varchar(100);not null" form:"state"`
	PostalCode string    `gorm:"type:varchar(20);not null" form:"postal_code"`
	Phone      string    `gorm:"type:varchar(20);not null" form:"phone"`
	Email      string    `gorm:"type:varchar(100);not null;unique" form:"email"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
