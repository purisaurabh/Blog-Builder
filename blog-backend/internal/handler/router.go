package handler

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	app.Post("/api/register", Register)
	app.Post("/api/login", Login)
	app.Post("/api/post", CreatePost)
	app.Get("/api/posts", AllPost)
	app.Get("/api/posts/:id", GetBlogPost)
	app.Put("/api/post/:id", UpdateBlogPost)
	app.Get("/api/unique", UniquePost)
	app.Delete("/api/post/:id", DeletePost)
	app.Post("/api/upload", UploadImage)
}
