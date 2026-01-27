package commons

import (
	"strings"

	"github.com/obadiaspelembe/linode-objim-cli/internal/commands/api"
)

func ProcessBucketObject(config Config, bucket, objectKey string) bool {
	
	SignedURL, err := api.GetLinodeSignedURL(config.Token, config.Region, bucket, objectKey, "text/plain", "GET")

	ErrorCheck(err, "Error getting signed URL", false)
	
	err = api.SaveToLocalFile(SignedURL,  objectKey)

	if err != nil {
		return false
	}

	return true
}

func ProcessLocalObject(config Config, bucket, objectKey, contentType string) bool {

	fileName := objectKey

	if strings.Contains(objectKey, "../") {
		fileName = strings.ReplaceAll(objectKey, "../", "")
	}
	SignedURL, err := api.GetLinodeSignedURL(config.Token, config.Region, bucket, fileName, contentType, "PUT")

	ErrorCheck(err, "Error getting signed URL", false)

	api.UploadToBucket(SignedURL, objectKey, contentType)
	if err != nil {
		return false
	}
	return true
}