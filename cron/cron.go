package cron

import (
	"log"

	"github.com/gonext-tech/internal/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func StartCron(store *gorm.DB, projectStores []models.ProjectsDB) {
	c := cron.New(cron.WithSeconds())

	// Schedule the job to run every night at midnight
	c.AddFunc("0 0 2 * * *", func() {
		err := CheckAndRenewSubscriptions(projectStores)
		if err != nil {
			log.Printf("Error renewing subscriptions: %v", err)
		}

	})
	c.Start()
}
