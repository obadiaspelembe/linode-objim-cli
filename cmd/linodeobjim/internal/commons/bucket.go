package commons

import ( 
	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal/commands/api"
)

func ProcessBucketObject(config Config, bucket, objectKey string) bool {
	
	SignedURL, err := api.GetLinodeSignedURL(config.Token, config.Region, bucket, objectKey, "text/plain")

	ErrorCheck(err, "Error getting signed URL")

	err = api.SaveToLocalFile(SignedURL, "./"+objectKey)

	ErrorCheck(err, "Error saving file locally")

	return true
}
