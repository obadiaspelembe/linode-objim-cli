package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/obadiaspelembe/linode-objim-cli/internal/linode"
)

type SignedURLRequest struct {
	Method      string `json:"method"`
	Name        string `json:"name"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

type AttachmentRequest struct {
	ContentDisposition string `json:"content_disposition,omitempty"`
	Method             string `json:"method"`
	Name               string `json:"name"`
}

type SignedURLResponse struct {
	SignedURL string `json:"url"`
}

func GetLinodeSignedURL(token, region, bucket, key, contentType, method string) (string, error) {

	apiClient := linode.NewAPI(token, region)

	body, _ := json.Marshal(
		AttachmentRequest{
			ContentDisposition: "attachment",
			Method:             "GET",
			Name:               key,
		})

	if method == "PUT" {
		body, _ = json.Marshal(
			SignedURLRequest{
				Method:      "PUT",
				Name:        key,
				ContentType: contentType,
				ExpiresIn:   3600,
			})
	}
	req := apiClient.InitializePostRequest(
		fmt.Sprintf("https://api.linode.com/v4/object-storage/buckets/%s/%s/object-url", region, bucket),
		body)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("linode API error %s: %s", res.Status, string(b))
	}

	var out SignedURLResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.SignedURL, nil
}

func SaveToLocalFile(presignedURL, localPath string) error {

	resp, err := http.Get(presignedURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed %s: %s", resp.Status, string(b))
	}

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0o755); err != nil { // rwx for owner, rx for group/others
		return fmt.Errorf("create dir %q: %w", dir, err)
	}

	f, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	return err
}

func UploadToBucket(presignedURL, localPath, contentType string) error {

	fileBytes, err := os.ReadFile(localPath)
	if err != nil {
		panic(err)
	}

	putReq, err := http.NewRequest("PUT", presignedURL, bytes.NewReader(fileBytes))
	if err != nil {
		panic(err)
	}

	putReq.Header.Set("Content-Type", contentType)

	putRes, err := http.DefaultClient.Do(putReq)
	if err != nil {
		return err
	}

	defer putRes.Body.Close()

	if putRes.StatusCode == 200 {
		return nil
	} else {
		_, err := io.ReadAll(putRes.Body)

		return err
	}
}
