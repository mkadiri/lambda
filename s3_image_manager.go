package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
	"image/jpeg"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type S3ImageManager struct {
	S3Session *session.Session
	Bucket string
}

func (s3ImageManager *S3ImageManager) download(key string) image.Image {
	buff := &aws.WriteAtBuffer{}
	s3dl := s3manager.NewDownloader(s3ImageManager.S3Session)

	_, err := s3dl.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(s3ImageManager.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Printf("Could not download from S3: %v", err)
	}

	log.Printf("-- Decoding image: %v bytes", len(buff.Bytes()))

	imageBytes := buff.Bytes()
	reader := bytes.NewReader(imageBytes)

	formattedImage, err := jpeg.Decode(reader)

	if err != nil {
		log.Printf("bad response: %s", err)
	}

	return formattedImage
}

// encode to jpg, keep the original filename and upload to a folder in the same directory but a size prefix
// e.g. /cover-images/1100x250
func (s3ImageManager *S3ImageManager) upload(folder string, image image.Image, key string) {
	log.Printf("-- Encoding image for upload to S3")
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)

	if err != nil {
		log.Printf("-- JPEG encoding error: %v", err)
	}

	originalFilename := filepath.Base(key)
	fileName := strings.TrimSuffix(originalFilename, path.Ext(key)) + ".jpg"
	outputPath := folder + "/" + fileName

	log.Printf("-- Saving file to: %v", outputPath)

	uploader := s3manager.NewUploader(s3ImageManager.S3Session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(s3ImageManager.Bucket),
		Key:    aws.String(outputPath),
	})

	if err != nil {
		log.Printf("-- Failed to upload: %v", err)
	}

	log.Printf("-- Successfully uploaded to: %v", result.Location)
}