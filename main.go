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

	validate := event.Validate()

	if validate != nil {
		exitErrorf(validate.Error())
	}

	log.Printf("Resizing images in the bucket %q for the folder %q to the size %dx%d", event.Bucket, event.Folder, event.Width, event.Height)

	resp := getImagesAtCurrentLevel(event.Bucket, event.Folder)

	for _, item := range resp.Contents {
		if !isObjectImage(*item.Key) {
			log.Println("File not an image, skip: ", *item.Key)
			continue
		}

		log.Println("Process image file: ", *item.Key)

		downloadedS3Image := downloadS3Image(event.Bucket, *item.Key)

		imageFormatter := ImageFormatter{}
		resizedImage := imageFormatter.resize(downloadedS3Image)
		resizedAndCroppedImage := imageFormatter.crop(resizedImage)

		uploadImage(event.Bucket, event.Folder, resizedAndCroppedImage, *item.Key)
	}

	return model.Response{Message: fmt.Sprintf("%s is %d years old!", event.Width, event.Height)}, nil
}

func isObjectImage(key string) bool {
	var validImageExtensions = [3]string{"jpg", "jpeg", "png"}

	for _, ext := range validImageExtensions {
		if strings.HasSuffix(key, "." + ext) {
			return true
		}
	}

	return false
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

// returns all the images in the current "folder" only and not the sub-folders
func getImagesAtCurrentLevel(bucket string, folder string) *s3.ListObjectsV2Output {
	log.Println("Retrieving images")

	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(folder),
			Delimiter: aws.String("/"),
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
func uploadImage(bucket string, folder string, image image.Image, key string) {
	log.Printf("-- Encoding image for upload to S3")
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)

	if err != nil {
		log.Printf("-- JPEG encoding error: %v", err) //todo: check what format we support
	}

	originalFilename := filepath.Base(key)
	fileName := strings.TrimSuffix(originalFilename, path.Ext(key)) + ".jpg"

	outputPath := folder + strconv.Itoa(maxWidth) + "x" + strconv.Itoa(maxHeight) + "/" + fileName

	log.Printf("-- Saving file to: %v", outputPath)

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(bucket),
		Key:    aws.String(outputPath),
	})

	if err != nil {
		log.Printf("-- Failed to upload: %v", err)
	}

	log.Printf("-- Successfully uploaded to: %v", result.Location)
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}