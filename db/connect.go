package db

import (
	"fmt"
	"log"
	"os"

	"github.com/prachi-satbhai0741/Taskboard/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" Failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&models.Task{}); err != nil {
		log.Fatalf(" Migration failed: %v", err)
	}

	log.Println(" Database connected and migrated")
}
