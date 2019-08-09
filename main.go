package main

import (
	"bitbucket.org/quidco/lambda/model"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
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

	resp := getObjectsAtCurrentLevel(event.Bucket, event.Folder)

	for _, item := range resp.Contents {
		if !isObjectImage(*item.Key) {
			log.Println("File not an image, skip: ", *item.Key)
			continue
		}

		log.Println("Process image file: ", *item.Key)

		s3ImageManager := S3ImageManager{event.Bucket}

		downloadedS3Image := s3ImageManager.download(*item.Key)

		bounds := downloadedS3Image.Bounds()

		if bounds.Max.X < maxWidth {
			log.Printf("-- downloaded image width %q is larger than the max width %q, skip resize", bounds.Max.X, maxWidth)
			continue
		}

		if bounds.Max.Y < maxHeight {
			log.Printf("-- downloaded image height %q is larger than the max height %q, skip resize", bounds.Max.X, maxWidth)
			continue
		}

		imageFormatter := ImageFormatter{}
		resizedImage := imageFormatter.resize(downloadedS3Image)
		resizedAndCroppedImage := imageFormatter.crop(resizedImage)

		s3ImageManager.upload(event.Folder, resizedAndCroppedImage, *item.Key)
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

// returns all objects in the current "folder" only and not the sub-folders
func getObjectsAtCurrentLevel(bucket string, folder string) *s3.ListObjectsV2Output {
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

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}