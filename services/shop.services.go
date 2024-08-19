package services

import (
	"sort"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/utils"
)

type ShopServices struct {
	STORES []models.ProjectsDB
}

func NewShopService(db []models.ProjectsDB) *ShopServices {
	return &ShopServices{
		STORES: db,
	}
}

func (ss *ShopServices) GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) ([]models.Shop, models.Meta, error) {
	allShops := []models.Shop{}
	var totalRecords int64
	for _, store := range ss.STORES {
		var shops []models.Shop
		if project != "" && project != store.Name {
			continue
		}
		if store.Name == "Support" {
			continue
		}
		query := store.DB.Table("shops").Order("created_at DESC")
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
		if err := query.Preload("Owner").Find(&shops).Error; err != nil {
			continue
		}

		// Append the result to the allCustomers slice
		allShops = append(allShops, shops...)
	}

	// Sorting all customers by date
	sort.SliceStable(allShops, func(i, j int) bool {
		if orderBy == "asc" {
			return allShops[i].CreatedAt.Before(allShops[j].CreatedAt)
		}
		return allShops[i].CreatedAt.After(allShops[j].CreatedAt)
	})

	// Pagination logic
	start := (page - 1) * limit
	end := start + limit
	if start > len(allShops) {
		start = len(allShops)
	}
	if end > len(allShops) {
		end = len(allShops)
	}
	paginatedShops := allShops[start:end]

	// Calculate total pages (lastPage)
	lastPage := (totalRecords + int64(limit) - 1) / int64(limit)

	// Construct pagination metadata
	meta := models.Meta{
		CurrentPage: page,
		TotalCount:  int(totalRecords),
		LastPage:    int(lastPage),
		Limit:       limit,
	}

	return paginatedShops, meta, nil
}

func (ss *ShopServices) Fetch(project string) ([]models.Shop, error) {
	var shops []models.Shop
	DB := utils.GetCurrentDB(project, ss.STORES)
	DB.Table("shopss").Find(&shops)
	return shops, nil
}

func (ss *ShopServices) GetID(id, dbName string) (models.Shop, error) {
	var shop models.Shop
	DB := utils.GetCurrentDB(dbName, ss.STORES)
	if result := DB.Table("shops").Preload("Owner").First(&shop, id); result.Error != nil {
		return models.Shop{}, result.Error
	}
	return shop, nil
}

func (ss *ShopServices) Create(shop models.Shop) (models.Shop, error) {

	DB := utils.GetCurrentDB(shop.ProjectName, ss.STORES)
	if result := DB.Table("shops").Create(&shop); result.Error != nil {
		return models.Shop{}, result.Error
	}
	return shop, nil
}
func (ss *ShopServices) Update(shop models.Shop) (models.Shop, error) {
	DB := utils.GetCurrentDB(shop.ProjectName, ss.STORES)
	if result := DB.Table("shops").Save(&shop); result.Error != nil {
		return models.Shop{}, result.Error
	}
	return shop, nil
}

func (ss *ShopServices) Delete(shop models.Shop) (models.Shop, error) {
	DB := utils.GetCurrentDB(shop.ProjectName, ss.STORES)
	if result := DB.Table("shops").Delete(&shop); result.Error != nil {
		return models.Shop{}, result.Error
	}
	return shop, nil
}
