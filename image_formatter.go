package main

import (
	"github.com/disintegration/imaging"
	"image"
	"log"
)

type ImageFormatter struct {}

// imaging.Resize() preserves the image aspect ratio if the height/width is 0
//
// if the new height calculated (with ratio taken in to consideration) is less than the max height we want to resize by
// width e.g.
//
//
// image 1:
//
// max dimensions 1100 x 250
// image = 1500 x 280
//
// height = (h / w) * max w
// height = 280/1500 * 1100 = (205) -> 1100 x 205
//
// width = (w / h) * max h
// width = 1500/280 * 250 = (1339) -> 1339 x 250
//
// resize by HEIGHT in order to keep min width of 1100 intact
//
//
// image 2:
//
// max dimensions 1100 x 250
// image = 1115 x 280
//
// height = (h / w) * max w
// height = 280/1115 * 1100 = (276) -> 1100 x 276
//
// width = (w / h) * max h
// width = 1115/280 * 250 = (996) -> 996 x 250
//
// resize by WIDTH in order to keep min height of 250 intact
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