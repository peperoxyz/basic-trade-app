package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	err error
)

func StartDB() {
	//  using godotenv for keeping the credentials
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbport := os.Getenv("DB_PORT")

	// connect to db
	config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbport, dbname)
 	db, err = gorm.Open(mysql.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	// db.Debug().AutoMigrate(models.Admin{}, models.Product{}, models.Variant{})
}

// function untuk get connection DB
func GetDB() *gorm.DB {
	return db
}