package main

import "strings"

type S3ImageValidator struct {}

func (S3ImageValidator) isS3KeyAnImage(key string) bool {
	var validImageExtensions = [3]string{"jpg", "jpeg", "png"}

	for _, ext := range validImageExtensions {
		if strings.HasSuffix(key, "." + ext) {
			return true
		}
	}

	return false
}