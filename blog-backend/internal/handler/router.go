package handler

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	app.Post("/api/register", Register)
	app.Post("/api/login", Login)
}
