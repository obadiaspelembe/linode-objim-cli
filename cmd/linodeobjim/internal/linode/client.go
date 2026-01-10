package linode

import "net/http"

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

func (api *API) InitializeRequest(url string) *http.Request {

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+api.Token)

	return req
}
