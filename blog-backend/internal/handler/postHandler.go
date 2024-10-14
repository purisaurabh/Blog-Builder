package handler

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/purisaurabh/blog-backend/internal/database"
	"github.com/purisaurabh/blog-backend/internal/helper"
	"github.com/purisaurabh/blog-backend/internal/models"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var blogPost models.Blog
	if err := c.BodyParser(&blogPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
		})
	}

	if err := database.DB.Create(&blogPost).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create post",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Congrutulations! your oost is live",
	})

}

func AllPost(c *fiber.Ctx) error {
	pages, err := strconv.Atoi(c.Query("page", "1"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid page number",
		})
	}

	limit := 5
	offset := (pages - 1) * limit
	var total int64

	var getBlogs []models.Blog

	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getBlogs)
	database.DB.Model(&models.Blog{}).Count(&total)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": getBlogs,
		"meta": map[string]interface{}{
			"total":     total,
			"page":      pages,
			"last_page": math.Ceil(float64(total) / float64(limit)),
		},
	})

}

func GetBlogPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	var blogPost models.Blog

	database.DB.Where("user_id = ?", id).Preload("User").First(&blogPost)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": blogPost,
	})
}

func UpdateBlogPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	blog := models.Blog{
		ID: id,
	}
	if err := c.BodyParser(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
		})
	}

	database.DB.Model(&blog).Updates(blog)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
	})
}

func UniquePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	fmt.Println("cookie is : ", cookie)
	id, err := helper.VerifyToken(cookie)
	if err != nil {
		fmt.Println("error while verifying the token  : ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	fmt.Print("id is : ", id)

	var blog []models.Blog
	database.DB.Model(&blog).Where("user_id = ?", id).Preload("User").Find(&blog)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": blog,
	})
}

func DeletePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	blog := models.Blog{
		ID: id,
	}

	deleteQuery := database.DB.Delete(&blog)
	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}
