package routes

import (
	"github.com/adnan007d/aeom-image-compressor/internals/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetuRoutes(app *fiber.App) {
	app.Static("/css", "dist/css")
	app.Static("/js", "dist/js")
	app.Static("/images", "images")

	app.Get("/", handlers.IndexPage)

	app.Post("/upload", handlers.CompressImages)

	app.Get("/zip/:id", handlers.CreateZip)

}
