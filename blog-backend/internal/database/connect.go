package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/purisaurabh/blog-backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file : ", err)
	}

	dns := os.Getenv("DB_DNS")

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	DB = db
	fmt.Println("Database connected")

	db.AutoMigrate(
		models.User{},
	)

}
