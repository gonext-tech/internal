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
