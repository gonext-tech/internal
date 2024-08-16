package services

import (
	"errors"
	"strings"

	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

type CityService struct {
	City models.City
	DB   *gorm.DB
}

func NewCityService(c models.City, db *gorm.DB) *CityService {
	return &CityService{
		City: c,
		DB:   db,
	}
}

func (cs *CityService) GetALL() ([]models.City, error) {
	var cities []models.City
	if result := cs.DB.Find(&cities); result.Error != nil {
		return []models.City{}, result.Error
	}
	return cities, nil
}

func (cs *CityService) GetID(id string) (models.City, error) {
	var city models.City
	if id == "0" || id == "" {
		return cs.City, errors.New("no id provided")
	}
	if result := cs.DB.First(&city, id); result.Error != nil {
		return cs.City, result.Error
	}
	return city, nil
}

func (cs *CityService) Create(city models.City) (models.City, error) {
	err := cs.DB.Create(&city).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return cs.City, errors.New("duplicate record: this city already exists")
		}
		return cs.City, err
	}
	return city, nil
}

func (cs *CityService) Update(city models.City) (models.City, error) {
	if result := cs.DB.Model(&city).Updates(city); result.Error != nil {
		if strings.Contains(result.Error.Error(), "Error 1062") {
			return cs.City, errors.New("duplicate record: this city already exists")
		}
		return cs.City, result.Error
	}
	return city, nil
}

func (cs *CityService) Delete(city models.City) error {
	if result := cs.DB.Delete(&city); result.Error != nil {
		return result.Error
	}
	return nil
}
