package storage

import (
	"fmt"
	"ionixx/api/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

/**
 * Function to initialize the database
 */
func InitDB() {
	var err error
	// Open database connection using the DSN
	DB, err = gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	// If there is an error opening the database, panic and exit
	if err != nil {
		panic(err)
	}
	fmt.Println("DB connected successfully!")
	DB.AutoMigrate(&models.User{})
}

/**
 * Function to initialize the test database
 */
func InitTestDB() {
	var err error
	// Open database connection using the DSN
	DB, err = gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	// If there is an error opening the database, panic and exit
	if err != nil {
		panic(err)
	}
	fmt.Println("DB connected successfully!")
	DB.Migrator().DropTable(&models.User{})
	DB.AutoMigrate(&models.User{})
}
