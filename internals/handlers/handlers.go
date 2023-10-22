package handlers

import (
	"bytes"
	"image"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"sync"

	"github.com/a-h/templ"
	"github.com/adnan007d/aeom-image-compressor/internals/views"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tidbyt/go-libwebp/webp"
	_ "golang.org/x/image/webp"

	"image/draw"
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
		go CompressImage(&wg, outChannel, CompressImageParams{
			File: file,
			Dir:  dir,
			CompressedImageOptions: CompressedImageOptions{
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
	// component := views.CompressedImages([]views.CompressedImagesType{})

	return renderTempl(c, component)

}

type CompressImageParams struct {
	File                   *multipart.FileHeader
	Dir                    string
	CompressedImageOptions CompressedImageOptions
}

type CompressedImageOptions struct {
	Quality   float32
	ImageType string
}

func CompressImage(wg *sync.WaitGroup, images chan views.CompressedImagesType, params CompressImageParams) {
	defer wg.Done()
	originalFileName := params.Dir + "/original/" + params.File.Filename
	compressedFileName := params.Dir + "/compressed/" + params.File.Filename + "." + params.CompressedImageOptions.ImageType

	out, err := os.Create(originalFileName)
	if err != nil {
		log.Printf("Error while creating %s: %v", params.File.Filename, err)
		return
	}
	defer out.Close()

	openedFile, err := params.File.Open()
	if err != nil {
		log.Printf("Error while opening infile %s: %v", params.File.Filename, err)
		return
	}
	defer openedFile.Close()

	var inMemoryImage bytes.Buffer

	multiWriter := io.MultiWriter(&inMemoryImage, out)

	io.Copy(multiWriter, openedFile)

	// Decoding Image to standard format
	decodedImage, _, err := image.Decode(&inMemoryImage)
	if err != nil {
		log.Printf("Error decoding image %s: %v", originalFileName, err)
		return
	}

	// log.Printf("%v is %v", originalFileName, iType)

	decodedOutFile, err := os.Create(compressedFileName)
	if err != nil {
		log.Printf("Error Opening outfile %s: %v", compressedFileName, err)
		return
	}
	defer decodedOutFile.Close()

	var convertError error = nil

	switch params.CompressedImageOptions.ImageType {
	case "webp":
		convertError = ConvertToWebp(decodedOutFile, decodedImage, params.CompressedImageOptions)

	default:
		log.Println("File Not supported")
	}

	if convertError != nil {
		log.Printf("Error while encoding image %s: %v", compressedFileName, convertError)
	}

	images <- views.CompressedImagesType{
		Source: views.ImageType{Src: originalFileName},
		Dest:   views.ImageType{Src: compressedFileName},
	}
}

func ConvertToWebp(w io.Writer, srcImage image.Image, options CompressedImageOptions) error {
	newImage := image.NewRGBA(srcImage.Bounds())

	draw.Draw(newImage, srcImage.Bounds(), srcImage, srcImage.Bounds().Min, draw.Src)

	config, err := webp.ConfigPreset(webp.PresetDefault, options.Quality)
	if err != nil {
		return err
	}

	err = webp.EncodeRGBA(w, newImage, config)
	if err != nil {
		return err
	}
	return nil
}
