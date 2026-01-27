package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/obadiaspelembe/linode-objim-cli/internal/linode"
)

type Object struct {
	Name         string `json:"name"`
	Size         int    `json:"size"`
	LastModified string `json:"last_modified"`
	ETag         string `json:"etag"`
}

type ObjectListResponse struct {
	Data        []Object `json:"data"`
	NextMarker  string   `json:"next_marker"`
	IsTruncated bool     `json:"is_truncated"`
}

func GetObjectList(bucketName string, region string, token string, recursive bool) ObjectListResponse {

	url := "https://api.linode.com/v4/object-storage/buckets/" + region + "/" + bucketName + "/object-list"


	if !recursive {
		url = url + "?delimiter=/&prefix="
	}

	apiClient := linode.NewAPI(token, region)

	req := apiClient.InitializeGetRequest(url)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var result ObjectListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	return result
}
