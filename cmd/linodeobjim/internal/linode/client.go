package linode

import (
	"bytes"
	"net/http"
)

type API struct {
	Token  string
	Region string
}

func NewAPI(token, region string) *API {
	return &API{
		Token:  token,
		Region: region,
	}
}

func (api *API) SetToken(token string) {
	api.Token = token
}

func (api *API) SetRegion(region string) {
	api.Region = region
}

func (api *API) InitializeGetRequest(url string) *http.Request {

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api.Token)

	return req
}


func (api *API) InitializePostRequest(url string, body []byte) *http.Request {

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api.Token)
	req.Header.Add("Content-Type", "application/json")

	return req
}
