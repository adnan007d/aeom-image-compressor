package util

import (
	"bytes"
	"image"
	// "image/draw"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/adnan007d/aeom-image-compressor/internals/views"
	"github.com/tidbyt/go-libwebp/webp"
	"golang.org/x/image/draw"
)

type CompressImageParams struct {
	File                   *multipart.FileHeader
	Dir                    string
	CompressedImageOptions CompressedImageOptions
}

type CompressedImageOptions struct {
	Quality   float32
	ImageType string
	Width     int
}

func CompressImage(wg *sync.WaitGroup, images chan views.CompressedImagesType, params CompressImageParams) {
	defer wg.Done()
	baseFileName := params.File.Filename[:len(params.File.Filename)-len(filepath.Ext(params.File.Filename))]
	originalFileName := params.Dir + "/original/" + params.File.Filename
	compressedFileName := params.Dir + "/compressed/" + baseFileName + "." + params.CompressedImageOptions.ImageType

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

	// NOTE: Maybe omit in memory and use the freshly written file?
	// Creating a writer to write to make a copy in memory of the image and write to disk
	var inMemoryImage bytes.Buffer
	multiWriter := io.MultiWriter(&inMemoryImage, out)
	io.Copy(multiWriter, openedFile)

	// Decoding Image to standard format
	decodedImage, _, err := image.Decode(&inMemoryImage)
	if err != nil {
		log.Printf("Error decoding image %s: %v", originalFileName, err)
		return
	}

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

	outFileStat, err := decodedOutFile.Stat()
	if err != nil {
		log.Printf("Error while reading file stat for %v: %v", decodedOutFile.Name(), err)
	}

	width := decodedImage.Bounds().Max.X

	if params.CompressedImageOptions.Width > 0 {
		width = params.CompressedImageOptions.Width
	}

	images <- views.CompressedImagesType{
		Source: views.ImageType{Src: originalFileName, Size: int(params.File.Size), Name: params.File.Filename},
		Dest:   views.ImageType{Src: compressedFileName, Size: int(outFileStat.Size()), Name: outFileStat.Name()},
		Width:  width,
	}
}

func ConvertToWebp(w io.Writer, srcImage image.Image, options CompressedImageOptions) error {
	bounds := srcImage.Bounds()

	if options.Width > 0 {
		// normal math to calculate the new height using to keep the same aspect ratio w1/h1 = w2/h2
		newHeight := (options.Width * srcImage.Bounds().Max.Y) / srcImage.Bounds().Max.X
		bounds = image.Rectangle{

			Min: image.Point{0, 0},
			Max: image.Point{X: options.Width, Y: newHeight},
		}
	}

	newImage := image.NewRGBA(bounds)
	draw.ApproxBiLinear.Scale(newImage, bounds, srcImage, srcImage.Bounds(), draw.Src, nil)

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
