package main

import (
	"github.com/disintegration/imaging"
	"image"
	"log"
)

type ImageFormatter struct {}

func (*ImageFormatter) resizeToRatioFromMaxDimensions(image image.Image, maxWidth int, maxHeight int) image.Image {
	bounds := image.Bounds()
	imageWidth := bounds.Max.X
	imageHeight := bounds.Max.Y

	newHeight := (imageHeight / imageWidth) * maxWidth

	if newHeight > maxHeight {
		log.Printf("-- Resizing image by maxHeight")
		return imaging.Resize(image, 0, maxHeight, imaging.Lanczos)

	} else {
		log.Printf("-- Resizing image by maxWidth")
		return imaging.Resize(image, maxWidth, 0, imaging.Lanczos)
	}
}

func (*ImageFormatter) crop(image image.Image, width int, height int) image.Image {
	bounds := image.Bounds()
	imageWidth := bounds.Max.X
	imageHeight := bounds.Max.Y

	if imageWidth == width && imageHeight == height {
		log.Printf("-- Image is in the correct dimension, no need to crop")
		return image
	}

	log.Printf("-- Cropping image to fit the max dimensions: %dx%d", width, height)

	return imaging.CropCenter(image, width, height)
}