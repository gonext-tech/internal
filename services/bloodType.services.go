package services

import (
	"errors"
	"strings"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type BloodService struct {
	Blood models.BloodType
	DB    *gorm.DB
}

func NewBloodService(b models.BloodType, db *gorm.DB) *BloodService {
	return &BloodService{
		Blood: b,
		DB:    db,
	}
}

func (bs *BloodService) GetALL() ([]models.BloodType, error) {
	var bloods []models.BloodType
	if result := bs.DB.Find(&bloods); result.Error != nil {
		return []models.BloodType{}, result.Error
	}
	return bloods, nil
}

func (bs *BloodService) GetID(id string) (models.BloodType, error) {
	var blood models.BloodType
	if id == "0" || id == "" {
		return bs.Blood, errors.New("no id provided")
	}
	if result := bs.DB.First(&blood, id); result.Error != nil {
		return bs.Blood, result.Error
	}
	return blood, nil
}

func (bs *BloodService) Create(blood models.BloodType) (models.BloodType, error) {
	err := bs.DB.Create(&blood).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return bs.Blood, errors.New("duplicate record: this blood type already exists")
		}
		return bs.Blood, err
	}
	return blood, nil
}

func (bs *BloodService) Update(blood models.BloodType) (models.BloodType, error) {
	if result := bs.DB.Model(&blood).Updates(blood); result.Error != nil {
		if strings.Contains(result.Error.Error(), "Error 1062") {
			return bs.Blood, errors.New("duplicate record: this blood type already exists")
		}
		return bs.Blood, result.Error
	}
	return blood, nil
}

func (bs *BloodService) Delete(blood models.BloodType) error {
	if result := bs.DB.Delete(&blood); result.Error != nil {
		return result.Error
	}
	return nil
}
