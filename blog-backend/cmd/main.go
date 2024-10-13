package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/purisaurabh/blog-backend/internal/database"
	"github.com/purisaurabh/blog-backend/internal/handler"
)

func main() {
	database.Connect()
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	port := os.Getenv("PORT")
	app := fiber.New()
	handler.Routes(app)
	app.Listen(":" + port)
}
