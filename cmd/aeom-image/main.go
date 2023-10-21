package main

import (
	"log"

	"github.com/adnan007d/aeom-image-compressor/internals/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})

	routes.SetuRoutes(app)
	log.Fatal(app.Listen(":6969"))
}
