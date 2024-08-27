package services

import (
	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/utils"
	"sort"
)

type SubscriptionServices struct {
	STORES []models.ProjectsDB
}

func NewSubscriptionService(db []models.ProjectsDB) *SubscriptionServices {
	return &SubscriptionServices{
		STORES: db,
	}
}

func (ss *SubscriptionServices) GetALL(limit, page int, orderBy, sortBy, project, shop, status, searchTerm string) ([]models.Subscription, models.Meta, error) {
	allSubscription := []models.Subscription{}
	var totalRecords int64
	for _, store := range ss.STORES {
		var subscriptions []models.Subscription
		if project != "" && project != store.Name {
			continue
		}
		if store.Name == "Support" {
			continue
		}
		query := store.DB.Table("subscriptions").Order("created_at DESC")

		if searchTerm != "" {
			searchTermWithWildcard := "%" + searchTerm + "%"
			query = query.Where("name LIKE ?", searchTermWithWildcard)
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}

		if shop != "" {
			query = query.Where("shop_id", shop)
		}

		// Applying sorting parameters
		query = query.Order(sortBy + " " + orderBy)

		var count int64
		if err := query.Count(&count).Error; err != nil {
			continue
		}
		totalRecords += count

		// Execute the query
		if err := query.Preload("Shop").Preload("Shop.Owner").Preload("Membership").Find(&subscriptions).Error; err != nil {
			continue
		}
		// Add the project name to each customer
		for i := range subscriptions {
			subscriptions[i].ProjectName = store.Name
		}

		// Append the result to the allCustomers slice
		allSubscription = append(allSubscription, subscriptions...)
	}

	// Sorting all customers by date
	sort.SliceStable(allSubscription, func(i, j int) bool {
		if orderBy == "asc" {
			return allSubscription[i].CreatedAt.Before(allSubscription[j].CreatedAt)
		}
		return allSubscription[i].CreatedAt.After(allSubscription[j].CreatedAt)
	})

	// Pagination logic
	start := (page - 1) * limit
	end := start + limit
	if start > len(allSubscription) {
		start = len(allSubscription)
	}
	if end > len(allSubscription) {
		end = len(allSubscription)
	}
	paginatedSubscription := allSubscription[start:end]

	// Calculate total pages (lastPage)
	lastPage := (totalRecords + int64(limit) - 1) / int64(limit)

	// Construct pagination metadata
	meta := models.Meta{
		CurrentPage: page,
		TotalCount:  int(totalRecords),
		LastPage:    int(lastPage),
		Limit:       limit,
	}

	return paginatedSubscription, meta, nil
}

func (ss *SubscriptionServices) GetID(id, dbName string) (models.Subscription, error) {
	var subscription models.Subscription
	DB := utils.GetCurrentDB(dbName, ss.STORES)
	if result := DB.Table("subscriptions").Preload("Membership").Preload("Shop").First(&subscription, id); result.Error != nil {
		return models.Subscription{}, result.Error
	}
	subscription.ProjectName = dbName
	return subscription, nil
}

func (ss *SubscriptionServices) Create(subscription models.Subscription) (models.Subscription, error) {
	DB := utils.GetCurrentDB(subscription.ProjectName, ss.STORES)
	if result := DB.Table("subscriptions").Create(&subscription); result.Error != nil {
		return models.Subscription{}, result.Error
	}
	return subscription, nil
}
func (ss *SubscriptionServices) Update(subscription models.Subscription) (models.Subscription, error) {
	DB := utils.GetCurrentDB(subscription.ProjectName, ss.STORES)
	if result := DB.Table("subscriptions").Updates(&subscription); result.Error != nil {
		return models.Subscription{}, result.Error
	}
	return subscription, nil
}

func (ss *SubscriptionServices) Delete(subscription models.Subscription) (models.Subscription, error) {
	DB := utils.GetCurrentDB(subscription.ProjectName, ss.STORES)
	if result := DB.Table("memberships").Delete(&subscription); result.Error != nil {
		return models.Subscription{}, result.Error
	}
	return subscription, nil
}
