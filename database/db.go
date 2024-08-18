package database

import (
	"log"
	"os"
	"time"

	"github.com/gonext-tech/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var counts int64

func DBInit() (d *gorm.DB, dd []models.ProjectsDB, err error) {
	database := connectToDB()
	if database == nil {
		log.Panic("Can't connect to Postgres!")
	}
	database.AutoMigrate(models.User{}, models.Project{}, models.CommitStats{}, models.Subscription{}, models.Membership{})
	var projects []models.Project
	database.Where("status = ?", "ACTIVE").Find(&projects)
	var projectsDB []models.ProjectsDB
	for _, project := range projects {
		conn := projectDB(project.DBName)
		if conn != nil {
			db := models.ProjectsDB{
				Name: project.Name,
				DB:   conn,
			}
			projectsDB = append(projectsDB, db)
		}
	}
	//tableNames := []string{"appointment_services", "appointments", "clients", "expenses", "services", "shops", "statistics"}

	// Loop through the table names and drop each table
	//for _, tableName := range tableNames {
	//if err := database.Migrator().DropTable(tableName); err != nil {
	// Handle error
	//	log.Fatalf("Failed to drop table %s: %v", tableName, err)
	//}
	//}
	return database, projectsDB, nil
}

func connectToDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := dbUser + ":" + dbPassword + "@(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
	for {
		connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if err != nil {
			log.Println("DB not yet ready ...")
			counts++
		} else {
			log.Println("Connected to DB!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds ...")
		time.Sleep(1 * time.Second)
		continue
	}
}

func projectDB(dbName string) *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dsn := dbUser + ":" + dbPassword + "@(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
	for {
		connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if err != nil {
			log.Println("DB not yet ready ...")
			counts++
		} else {
			log.Println("Connected to DB!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds ...")
		time.Sleep(1 * time.Second)
		continue
	}
}
