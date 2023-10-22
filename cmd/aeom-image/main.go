package main

import (
	"log"
	"os"
	"time"

	"github.com/adnan007d/aeom-image-compressor/internals/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
  go DeleteImagesDirectoryPeriodically()

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})

	routes.SetuRoutes(app)

	log.Fatal(app.Listen(":6969"))
}

func DeleteImagesDirectoryPeriodically() {
	for {
		os.RemoveAll("images")
		time.Sleep(time.Minute * 5)
	}
}
