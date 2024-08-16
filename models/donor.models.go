package models

import (
	"time"

	"gorm.io/gorm"
)

type DonorBody struct {
	ID          uint   `gorm:"primaryKey" form:"id"`
	Name        string `gorm:"size:100;not null" form:"name"`
	Phone       string `gorm:"size:100;unique;not null" form:"phone"`
	Gender      string `form:"gender"`
	DateOfBirth string `form:"date_of_birth" `
	BloodTypeID uint   `form:"blood_type_id"`
	CityID      uint   `form:"city_id"`
	Address     string `form:"address"`
}

type Donor struct {
	ID          uint      `gorm:"primaryKey" form:"id"`
	Name        string    `gorm:"size:100;not null" form:"name"`
	Phone       string    `gorm:"size:100;unique;not null" form:"phone"`
	Gender      string    `form:"gender"`
	DateOfBirth time.Time `gorm:"not null" form:"date_of_birth" `
	BloodTypeID uint      `form:"blood_type_id"`
	BloodType   BloodType `form:"blood_type" gorm:"foreignKey:BloodTypeID"`
	CityID      uint      `form:"city_id"`
	City        City      `form:"city" gorm:"foreignKey:CityID"`
	Address     string    `form:"address"`
	CreatedAt   time.Time `form:"created_at"`
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type BloodType struct {
	ID   uint   `gorm:"primaryKey" form:"id"`
	Type string `gorm:"size:3;unique;not null" form:"type"`
}

type City struct {
	ID        uint      `gorm:"primaryKey" form:"id"`
	Name      string    `gorm:"size:100;not null;unique" form:"name"`
	CreatedAt time.Time `form:"created_at"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalParam(s string) error {
	t, err := time.Parse("2006-01-02", s) // Adjust format as needed
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}
