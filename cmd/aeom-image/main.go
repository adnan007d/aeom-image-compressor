package main

import (
	"log"

	"github.com/a-h/templ"
	"github.com/adnan007d/aeom-image-compressor/internals/views"
	"github.com/gofiber/fiber/v2"
)

func renderTempl(c *fiber.Ctx, component templ.Component) error {
	c.Response().Header.Set("Content-Type", "text/html")

	return component.Render(c.Context(), c.Response().BodyWriter())
}

func main() {
	app := fiber.New()


  app.Static("/css", "dist/css")
  app.Static("/js", "dist/js")

	app.Get("/", func(c *fiber.Ctx) error {
		component := views.Index("wtf")
		return renderTempl(c, component)
	})

	log.Fatal(app.Listen(":6969"))
}
