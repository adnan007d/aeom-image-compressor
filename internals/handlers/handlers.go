package handlers

import (
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/adnan007d/aeom-image-compressor/internals/views"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func renderTempl(c *fiber.Ctx, component templ.Component) error {
	c.Response().Header.Set("Content-Type", "text/html")

	return component.Render(c.Context(), c.Response().BodyWriter())
}


func IndexPage(c *fiber.Ctx) error {
	component := views.Index("wtf")
	return renderTempl(c, component)
}

func CompressImages(c *fiber.Ctx) error {
		files, err := c.MultipartForm()
		if err != nil {
			log.Printf("Error while decoding multiplat form: %v", err)
		}

		var dir = "images/" + uuid.New().String()

		os.MkdirAll(dir, fs.ModePerm)

		hell := make([]views.ImagesType, len(files.File))

		for _, file := range files.File["images"] {
			out, err := os.Create(dir + "/" + file.Filename)
			if err != nil {
				log.Printf("Error while creating %s: %v", file.Filename, err)
			}
			defer out.Close()

			infile, err := file.Open()
			if err != nil {
				log.Printf("Error while opening infile %s: %v", file.Filename, err)
			}
			defer infile.Close()

			io.Copy(out, infile)
			hell = append(hell, views.ImagesType{
				Src: dir + "/" + file.Filename,
			})
			log.Println(file.Filename)
		}

		component := views.CompressedImages(hell)

		return renderTempl(c, component)

	}
