package api

import (
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

func GetLinodeSignedURL(token, region, bucket, key, contentType string) (string, error) {

	apiClient := linode.NewAPI(token, region)

	body, _ := json.Marshal(
		AttachmentRequest{
			ContentDisposition: "attachment",
			Method:             "GET",
			Name:               key,
		})

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
