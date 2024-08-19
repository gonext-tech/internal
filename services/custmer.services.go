package services

import (
	"sort"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/utils"
)

type CustomerServices struct {
	STORES []models.ProjectsDB
}

func NewCustomerService(db []models.ProjectsDB) *CustomerServices {
	return &CustomerServices{
		STORES: db,
	}
}

func (cs *CustomerServices) GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) ([]models.Customer, models.Meta, error) {
	allCustomers := []models.Customer{}
	var totalRecords int64
	for _, store := range cs.STORES {
		var customers []models.Customer
		if project != "" && project != store.Name {
			continue
		}
		if store.Name == "Support" {
			continue
		}
		query := store.DB.Table("users").Order("created_at DESC")

		if searchTerm != "" {
			searchTermWithWildcard := "%" + searchTerm + "%"
			query = query.Where("name LIKE ?", searchTermWithWildcard)
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}

		// Applying sorting parameters
		query = query.Order(sortBy + " " + orderBy)

		var count int64
		if err := query.Count(&count).Error; err != nil {
			continue
		}
		totalRecords += count

		// Execute the query
		if err := query.Find(&customers).Error; err != nil {
			continue
		}

		// Append the result to the allCustomers slice
		allCustomers = append(allCustomers, customers...)
	}

	// Sorting all customers by date
	sort.SliceStable(allCustomers, func(i, j int) bool {
		if orderBy == "asc" {
			return allCustomers[i].CreatedAt.Before(allCustomers[j].CreatedAt)
		}
		return allCustomers[i].CreatedAt.After(allCustomers[j].CreatedAt)
	})

	// Pagination logic
	start := (page - 1) * limit
	end := start + limit
	if start > len(allCustomers) {
		start = len(allCustomers)
	}
	if end > len(allCustomers) {
		end = len(allCustomers)
	}
	paginatedCusotmer := allCustomers[start:end]

	// Calculate total pages (lastPage)
	lastPage := (totalRecords + int64(limit) - 1) / int64(limit)

	// Construct pagination metadata
	meta := models.Meta{
		CurrentPage: page,
		TotalCount:  int(totalRecords),
		LastPage:    int(lastPage),
		Limit:       limit,
	}

	return paginatedCusotmer, meta, nil
}

func (cs *CustomerServices) GetID(id, dbName string) (models.Customer, error) {
	var customer models.Customer
	DB := utils.GetCurrentDB(dbName, cs.STORES)
	if result := DB.Table("users").First(&customer, id); result.Error != nil {
		return models.Customer{}, result.Error
	}
	return customer, nil
}

func (cs *CustomerServices) Create(customer models.Customer) (models.Customer, error) {
	DB := utils.GetCurrentDB(customer.ProjectName, cs.STORES)
	if result := DB.Table("users").Create(&customer); result.Error != nil {
		return models.Customer{}, result.Error
	}
	return customer, nil
}
func (cs *CustomerServices) Update(customer models.Customer) (models.Customer, error) {
	DB := utils.GetCurrentDB(customer.ProjectName, cs.STORES)
	if result := DB.Table("users").Updates(&customer); result.Error != nil {
		return models.Customer{}, result.Error
	}
	return customer, nil
}

func (cs *CustomerServices) Delete(customer models.Customer) (models.Customer, error) {
	DB := utils.GetCurrentDB(customer.ProjectName, cs.STORES)
	if result := DB.Table("users").Delete(&customer); result.Error != nil {
		return models.Customer{}, result.Error
	}
	return customer, nil
}
