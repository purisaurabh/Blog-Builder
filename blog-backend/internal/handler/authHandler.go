package handler

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/purisaurabh/blog-backend/internal/database"
	"github.com/purisaurabh/blog-backend/internal/helper"
	"github.com/purisaurabh/blog-backend/internal/models"
)

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User

	if err := c.BodyParser(&data); err != nil {
		log.Fatalf("Error parsing request body: %v", err)
	}

	// check the password is less than 8 characters
	if len(data["password"].(string)) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 8 characters",
		})
	}

	// check the email is valid
	if !helper.ValidEmail(strings.TrimSpace(data["email"].(string))) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.ID != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Email:     data["email"].(string),
		Phone:     data["phone"].(string),
	}

	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":    user,
		"message": "User registered successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Fatalf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to login",
		})
	}

	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if !user.ComparePassword(data["password"]) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	token, err := helper.GenerateToken(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}
