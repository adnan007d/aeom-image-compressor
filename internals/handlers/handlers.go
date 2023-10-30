package handlers

import (
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	_ "image/jpeg"
	_ "image/png"

	"github.com/a-h/templ"
	"github.com/adnan007d/aeom-image-compressor/internals/util"
	"github.com/adnan007d/aeom-image-compressor/internals/views"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
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

	formData := extractFormdata(form)

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
				Quality:   formData.quality,
				ImageType: formData.format,
				Width:     formData.width,
			},
		})

	}
	wg.Wait()

	close(outChannel)

	var images []views.CompressedImagesType
	for channelValue := range outChannel {
		images = append(images, channelValue)
	}

	component := views.CompressedImages(images, randomUUID)
	return renderTempl(c, component)
}

type FormData struct {
	quality float32
	width   int
	format  string
}

func extractFormdata(form *multipart.Form) FormData {
	formData := FormData{
		quality: 75,
		width:   0,
		format:  "webp",
	}
	if len(form.Value["quality"]) > 0 {
		quality, err := strconv.ParseFloat(form.Value["quality"][0], 32)
		// if there is error the default value will be used
		if err == nil {
			formData.quality = float32(quality)
		}
	}

	if len(form.Value["width"]) > 0 {
		width, err := strconv.ParseInt(form.Value["width"][0], 10, 64)
		if err == nil {
			formData.width = int(width)
		}
	}

	if len(form.Value["format"]) > 0 && (form.Value["format"][0] == "webp" || form.Value["format"][0] == "jpg" || form.Value["format"][0] == "png") {
		formData.format = form.Value["format"][0]
	}

	return formData
}

func CreateZip(c *fiber.Ctx) error {
	id := c.Params("id")

	targetPath := path.Join("images", id)
	if _, err := os.Stat(path.Join("images", id)); err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	targetFile := path.Join(targetPath, id+".zip")

	compressedDir := path.Join(targetPath, "compressed")

	err := util.ZipItAndShipIt(compressedDir, targetFile)

	if err != nil {
		log.Printf("Error while creating zip :%v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendFile(targetFile)
}
