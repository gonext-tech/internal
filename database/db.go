package database

import (
	"log"
	"os"
	"time"

	"github.com/gonext-tech/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var counts int64

func DBInit() (d *gorm.DB, dd []models.ProjectsDB, err error) {
	database := connectToDB()
	if database == nil {
		log.Panic("Can't connect to Mysql!")
	}
	database.AutoMigrate(models.Admin{}, models.Client{}, models.Project{}, models.Invoice{}, models.CommitStats{}, models.Subscription{}, models.Membership{}, models.Referal{}, models.Stats{}, models.MonitoredServer{}, models.Domain{})
	var projects []models.Project
	database.Where("status = ?", "ACTIVE").Find(&projects)
	var projectsDB []models.ProjectsDB
	var admin models.Admin
	database.Where("email = ?", "admin@gmail.com").First(&admin)
	if admin.ID == 0 {
		createAdmin(database)
	}

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

func createAdmin(db *gorm.DB) {
	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), 10)
	if err != nil {
		log.Println(err)
	}
	admin := models.Admin{
		Email:    "admin@gmail.com",
		Password: string(hash),
	}
	db.Create(&admin)
}
