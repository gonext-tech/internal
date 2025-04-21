package services

import (
	"sort"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/utils"
)

type MembershipServices struct {
	STORES []models.ProjectsDB
}

func NewMembershipService(db []models.ProjectsDB) *MembershipServices {
	return &MembershipServices{
		STORES: db,
	}
}

func (ms *MembershipServices) GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) ([]models.Membership, models.Meta, error) {
	allMemberships := []models.Membership{}
	var totalRecords int64
	for _, store := range ms.STORES {
		var memberships []models.Membership
		if project != "" && project != store.Name {
			continue
		}
		if store.Name == "Support" {
			continue
		}
		query := store.DB.Table("memberships").Order("created_at DESC")

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
		if err := query.Find(&memberships).Error; err != nil {
			continue
		}
		// Add the project name to each customer
		for i := range memberships {
			memberships[i].ProjectName = store.Name
		}

		// Append the result to the allCustomers slice
		allMemberships = append(allMemberships, memberships...)
	}

	// Sorting all customers by date
	sort.SliceStable(allMemberships, func(i, j int) bool {
		if orderBy == "asc" {
			return allMemberships[i].CreatedAt.Before(allMemberships[j].CreatedAt)
		}
		return allMemberships[i].CreatedAt.After(allMemberships[j].CreatedAt)
	})

	// Pagination logic
	start := (page - 1) * limit
	end := start + limit
	if start > len(allMemberships) {
		start = len(allMemberships)
	}
	if end > len(allMemberships) {
		end = len(allMemberships)
	}
	paginatedMemberships := allMemberships[start:end]

	// Calculate total pages (lastPage)
	lastPage := (totalRecords + int64(limit) - 1) / int64(limit)

	// Construct pagination metadata
	meta := models.Meta{
		CurrentPage: page,
		TotalCount:  int(totalRecords),
		LastPage:    int(lastPage),
		Limit:       limit,
	}

	return paginatedMemberships, meta, nil
}

func (ms *MembershipServices) Fetch(project string) ([]models.Membership, error) {
	var memberships []models.Membership
	DB := utils.GetCurrentDB(project, ms.STORES)
	DB.Table("memberships").Find(&memberships)
	return memberships, nil
}

func (ms *MembershipServices) GetID(id, dbName string) (models.Membership, error) {
	var membership models.Membership
	DB := utils.GetCurrentDB(dbName, ms.STORES)
	if result := DB.Table("memberships").First(&membership, id); result.Error != nil {
		return models.Membership{}, result.Error
	}
	return membership, nil
}

func (ms *MembershipServices) Create(membership models.Membership) (models.Membership, error) {
	DB := utils.GetCurrentDB(membership.ProjectName, ms.STORES)
	if result := DB.Table("memberships").Create(&membership); result.Error != nil {
		return models.Membership{}, result.Error
	}
	return membership, nil
}
func (ms *MembershipServices) Update(membership models.Membership) (models.Membership, error) {
	DB := utils.GetCurrentDB(membership.ProjectName, ms.STORES)
	if result := DB.Table("memberships").Updates(&membership); result.Error != nil {
		return models.Membership{}, result.Error
	}
	return membership, nil
}

func (ms *MembershipServices) Delete(membership models.Membership) (models.Membership, error) {
	DB := utils.GetCurrentDB(membership.ProjectName, ms.STORES)
	if result := DB.Table("memberships").Delete(&membership); result.Error != nil {
		return models.Membership{}, result.Error
	}
	return membership, nil
}
