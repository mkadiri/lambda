package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

type S3ObjectManager struct {
	S3Client *s3.S3
	Bucket string
}

func (s3ObjectManager *S3ObjectManager) getObjectsListAtCurrentLevel(folder string) (*s3.ListObjectsV2Output, error) {
	log.Printf("Retrieving list of objects at current level %q", folder)

	resp, err := s3ObjectManager.S3Client.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(s3ObjectManager.Bucket),
			Prefix: aws.String(folder),
			Delimiter: aws.String("/"),
		})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to list items in bucket %q, %v", s3ObjectManager.Bucket, err))
	}

	if len(resp.Contents) == 0 {
		return nil, errors.New(fmt.Sprintf("Folder path %q in bucket %q doesn't exist", folder, s3ObjectManager.Bucket))
	}

	return resp, nil
}