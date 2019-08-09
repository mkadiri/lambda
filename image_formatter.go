package main

import (
	"github.com/disintegration/imaging"
	"image"
	"log"
)

type ImageFormatter struct {
	MaxWidth int
	MaxHeight int
}

func (imageFormatter *ImageFormatter) resize(image image.Image) image.Image {
	bounds := image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	newHeight := (height / width) * imageFormatter.MaxWidth

	if (newHeight > imageFormatter.MaxHeight) {
		log.Printf("-- Resizing image by height")
		return imaging.Resize(image, 0, imageFormatter.MaxHeight, imaging.Lanczos)

	} else {
		log.Printf("-- Resizing image by width")
		return imaging.Resize(image, imageFormatter.MaxWidth, 0, imaging.Lanczos)
	}
}

func (imageFormatter *ImageFormatter) crop(image image.Image) image.Image {
	bounds := image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	if width == imageFormatter.MaxWidth && height == imageFormatter.MaxHeight {
		log.Printf("-- Image is in the correct dimension, no need to crop")
		return image
	}

	log.Printf("-- Cropping image to fit the max dimensions: %dx%d", imageFormatter.MaxWidth, imageFormatter.MaxHeight)

	return imaging.CropCenter(image, imageFormatter.MaxWidth, imageFormatter.MaxHeight)
}