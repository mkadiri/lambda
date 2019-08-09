package main

import (
	"bitbucket.org/quidco/lambda/model"
	"bytes"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var svc *s3.S3
var sess *session.Session

const maxWidth = 1100
const maxHeight = 250

func main() {
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(event model.Event) (model.Response, error) {
	initS3Client()

	bucket := os.Getenv("S3_BUCKET_NAME")
	folder := os.Getenv("S3_MERCHANT_COVER_PHOTOS_FOLDER_PATH")

	resp := getBannerImages(bucket, folder)

	for _, item := range resp.Contents {
		if !strings.HasSuffix(*item.Key, ".jpg") { //todo add other formats
			fmt.Println("Not an image, skip")
			continue
		}

		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")

		downloadedS3Image := downloadS3Image(bucket, *item.Key)
		resizedImage := resizeImage(downloadedS3Image)
		resizedAndCroppedImage := cropImage(resizedImage)

		uploadImage(bucket, resizedAndCroppedImage, *item.Key)
	}

	return model.Response{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
}

func initS3Client() {
	var err error

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION_NAME"))},
	)

	if err != nil {
		exitErrorf("Error initialising S3 client")
	}

	svc = s3.New(sess)
}

func getBannerImages(bucket string, folder string) *s3.ListObjectsV2Output {
	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(folder),
		})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	if len(resp.Contents) == 0 {
		exitErrorf("Folder path %q in bucket %q doesn't exist", folder, bucket)
	}

	return resp
}

func downloadS3Image(bucket string, key string) image.Image {
	buff := &aws.WriteAtBuffer{}
	s3dl := s3manager.NewDownloader(sess)
	_, err := s3dl.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Printf("Could not download from S3: %v", err)
	}

	log.Printf("Decoding image: %v bytes", len(buff.Bytes()))

	imageBytes := buff.Bytes()
	reader := bytes.NewReader(imageBytes)

	formattedImage, err := jpeg.Decode(reader)

	if err != nil {
		log.Printf("bad response: %s", err)
	}

	return formattedImage
}

func resizeImage(image image.Image) image.Image{
	bounds := image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	newHeight := (height / width) * maxWidth

	if (newHeight > maxHeight) {
		log.Printf("Resizing image by height")
		return imaging.Resize(image, 0, maxHeight, imaging.Lanczos)

	} else {
		log.Printf("Resizing image by width")

		return imaging.Resize(image, maxWidth, 0, imaging.Lanczos)
	}
}

func cropImage(image image.Image) image.Image {
	bounds := image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	if (width == maxWidth && height == maxHeight) {
		log.Printf("Image is in the correct dimension, no need to crop")
		return image
	}

	log.Printf("Cropping image to fit the max dimensions")

	return imaging.CropCenter(image, maxWidth, maxHeight)
}

// encode to jpeg, keep the original filename and upload to a folder in the same directory e.g. /cover-images/1100x250
func uploadImage(bucket string, image image.Image, key string) {
	log.Printf("Encoding image for upload to S3")
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)

	if err != nil {
		log.Printf("JPEG encoding error: %v", err) //todo: check what format we support
	}

	originalFilename := filepath.Base(key)
	fileName := strings.TrimSuffix(originalFilename, path.Ext(key)) + ".jpg"

	outputPath := os.Getenv("S3_MERCHANT_COVER_PHOTOS_FOLDER_PATH") +
		"/" + strconv.Itoa(maxWidth) + "x" + strconv.Itoa(maxHeight) +
		"/" + fileName

	log.Printf("Saving file to: %v", outputPath)

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(bucket),
		Key:    aws.String(outputPath),
	})

	if err != nil {
		log.Printf("Failed to upload: %v", err)
	}

	log.Printf("Successfully uploaded to: %v", result.Location)
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}