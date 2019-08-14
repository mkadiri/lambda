package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mkadiri/lambda-image-resizer/model"
	"log"
	"os"
	"strconv"
)

func main() {
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(event model.Event) (model.Response, error) {
	validate := event.Validate()

	if validate != nil {
		exitErrorf(validate.Error())
	}

	log.Printf("Resizing images in the bucket %q for the folder %q to the size %dx%d", event.Bucket, event.Folder, event.Width, event.Height)

	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION_NAME"))},
	)

	if err != nil {
		exitErrorf("Error initialising S3 client")
	}

	s3Client := s3.New(s3Session)

	s3ObjectManager := S3ObjectManager{s3Client, event.Bucket}
	resp, err := s3ObjectManager.getObjectsListAtCurrentLevel(event.Folder)

	if err != nil {
		exitErrorf(err.Error())
	}

	processS3Objects(resp.Contents, s3Session, event)

	return model.Response{Message: "Resizing complete"}, nil
}

func processS3Objects(objects []*s3.Object, s3Session *session.Session, event model.Event) {
	s3ImageValidator := S3ImageValidator{}

	for _, item := range objects {
		if !s3ImageValidator.isS3KeyAnImage(*item.Key) {
			log.Printf("Object %q not an image, skip", *item.Key)
			continue
		}

		log.Printf("Process image %q ", *item.Key)

		s3ImageManager := S3ImageManager{s3Session, event.Bucket}
		downloadedS3Image := s3ImageManager.download(*item.Key)
		bounds := downloadedS3Image.Bounds()

		if bounds.Max.X < event.Width {
			log.Printf("-- downloaded image width %q is smaller than the max width %q, skip resize", bounds.Max.X, event.Width)
			continue
		}

		if bounds.Max.Y < event.Height {
			log.Printf("-- downloaded image height %q is smaller than the max height %q, skip resize", bounds.Max.X, event.Height)
			continue
		}

		imageFormatter := ImageFormatter{}
		resizedImage := imageFormatter.resizeToRatioFromMaxDimensions(downloadedS3Image,event.Width, event.Height)
		resizedAndCroppedImage := imageFormatter.crop(resizedImage,event.Width, event.Height)

		outputFolder := event.Folder + strconv.Itoa(event.Width) + "x" + strconv.Itoa(event.Height)

		s3ImageManager.upload(outputFolder, resizedAndCroppedImage, *item.Key)
	}
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}