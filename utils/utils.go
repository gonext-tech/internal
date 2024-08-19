package utils

import (
	"github.com/gonext-tech/internal/models"
	"gorm.io/gorm"
)

func GetCurrentDB(dbName string, pDB []models.ProjectsDB) *gorm.DB {
	for _, db := range pDB {
		if db.Name == dbName {
			return db.DB
		}
	}
	return nil
}

func ConvertCustomerToUser(customer models.Customer) models.User {
	user := models.User{
		ID:      customer.ID,
		Email:   customer.Email,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
		//Status:    customer.Status,
		Password:  customer.Password,
		Role:      customer.Role,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
	return user
}

func ConvertShopToCustomerShop(shop models.Shop) models.CustomerShop {
	// var workers []models.User
	// for _, worker := range shop.Workers {
	// 	shopWorker := ConvertCustomerToUser(worker)
	// 	workers = append(workers, shopWorker)
	// }

	customerShop := models.CustomerShop{
		ID:          shop.ID,
		Name:        shop.Name,
		Address:     shop.Address,
		OwnerID:     shop.OwnerID,
		Status:      shop.Status,
		CreatedAt:   shop.CreatedAt,
		ProjectName: shop.ProjectName,
		Workers:     shop.Workers,
		UpdatedAt:   shop.UpdatedAt,
	}
	return customerShop
}
