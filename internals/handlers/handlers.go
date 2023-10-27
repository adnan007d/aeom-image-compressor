package handlers

import (
	"io/fs"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/a-h/templ"
	"github.com/adnan007d/aeom-image-compressor/internals/util"
	"github.com/adnan007d/aeom-image-compressor/internals/views"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
	_ "image/png"
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
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("Error while decoding multiplat form: %v", err)
	}

	parsedQuality, err := strconv.ParseFloat(form.Value["quality"][0], 32)

	if err != nil {
		log.Printf("Error while decoding quality %v: %v", form.Value["quality"], err)
	}
	quality := float32(parsedQuality)

	randomUUID := uuid.New().String()
	var dir = "images/" + randomUUID

	os.MkdirAll(dir+"/original", fs.ModePerm)
	os.MkdirAll(dir+"/compressed", fs.ModePerm)

	outChannel := make(chan views.CompressedImagesType, len(form.File["images"]))
	var wg sync.WaitGroup

	for _, file := range form.File["images"] {
		wg.Add(1)
		go util.CompressImage(&wg, outChannel, util.CompressImageParams{
			File: file,
			Dir:  dir,
			CompressedImageOptions: util.CompressedImageOptions{
				Quality:   quality,
				ImageType: "webp",
			},
		})

	}
	wg.Wait()

	close(outChannel)

	var images []views.CompressedImagesType
	for channelValue := range outChannel {
		images = append(images, channelValue)
	}

	component := views.CompressedImages(images)
	return renderTempl(c, component)
}
