package main

import (
	"github.com/disintegration/imaging"
	"image"
	"log"
)

type ImageFormatter struct {}

func (imageFormatter *ImageFormatter) resize(image image.Image) image.Image {
	bounds := image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	newHeight := (height / width) * maxWidth

	if (newHeight > maxHeight) {
		log.Printf("-- Resizing image by height")
		return imaging.Resize(image, 0, maxHeight, imaging.Lanczos)

	} else {
		log.Printf("-- Resizing image by width")
		return imaging.Resize(image, maxWidth, 0, imaging.Lanczos)
	}
}

func (imageFormatter *ImageFormatter) crop(image image.Image) image.Image {
	bounds := image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	if width == maxWidth && height == maxHeight {
		log.Printf("-- Image is in the correct dimension, no need to crop")
		return image
	}

	log.Printf("-- Cropping image to fit the max dimensions: %dx%d", maxWidth, maxHeight)

	return imaging.CropCenter(image, maxWidth, maxHeight)
}