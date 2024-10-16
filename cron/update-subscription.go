package cron

import (
	"github.com/gonext-tech/internal/models"
	"time"
)

func CheckAndRenewSubscriptions(projectStores []models.ProjectsDB) error {
	// Get the current date
	currentDate := time.Now().Truncate(24 * time.Hour)

	var subscriptions []models.Subscription
	for _, store := range projectStores {
		// Fetch all subscriptions where the endDate is today and the status is "ACTIVE"
		result := store.DB.Preload("Membership").Preload("Shop").Where("end_date = ? AND status = ?", currentDate, "ACTIVE").Find(&subscriptions)
		if result.Error != nil {
			return result.Error
		}

		if len(subscriptions) == 0 {
			continue
		}

		// Loop through each subscription and renew it
		for _, sub := range subscriptions {
			if sub.Shop.Status != "ACTIVE" {
				continue
			}
			// Create a new subscription for the shop with updated start and end dates
			newSub := models.Subscription{
				ShopID:        sub.ShopID,
				MembershipID:  sub.MembershipID,
				StartDate:     currentDate,
				EndDate:       currentDate.AddDate(0, sub.Membership.Duration, 0), // Add the membership duration (in months)
				Currency:      sub.Currency,
				AutoRenewal:   sub.AutoRenewal,
				Amount:        sub.Amount,
				ProjectName:   sub.ProjectName,
				PaymentMethod: sub.PaymentMethod,
			}

			if err := store.DB.Create(&newSub).Error; err != nil {
				return err
			}

			// Update the old subscription's status to NOT_ACTIVE
			sub.Status = "NOT_ACTIVE"
			if err := store.DB.Save(&sub).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
