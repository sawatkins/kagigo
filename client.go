package kagi

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/httpjamesm/kagi-ai-go/constants"
)

type ClientConfig struct {
	APIKey     string
	APIVersion constants.ApiVersion
}

type Client struct {
	Config *ClientConfig
}

func NewClient(config *ClientConfig) *Client {
	return &Client{Config: config}
}

func (c *Client) GetAPIKey() string {
	return c.Config.APIKey
}

func (c *Client) SetAPIKey(apiKey string) {
	c.Config.APIKey = apiKey
}

func (c *Client) GetAPIVersion() constants.ApiVersion {
	return c.Config.APIVersion
}

func (c *Client) SetAPIVersion(apiVersion constants.ApiVersion) {
	c.Config.APIVersion = apiVersion
}

func (c *Client) getBaseURL() string {
	return constants.BASE_URL + "/" + string(c.Config.APIVersion)
}

func (c *Client) SendRequest(method, path string, data map[string]interface{}) (res []byte, err error) {

	baseURL := c.getBaseURL()

	client := resty.New()

	reqBuild := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bot %s", c.Config.APIKey)).
		SetBody(data)

	var resp *resty.Response

	switch method {
	case "GET":
		resp, err = reqBuild.Get(baseURL + path)
	case "POST":
		resp, err = reqBuild.Post(baseURL + path)
	case "PUT":
		resp, err = reqBuild.Put(baseURL + path)
	case "DELETE":
		resp, err = reqBuild.Delete(baseURL + path)
	default:
		err = fmt.Errorf("invalid method: %s", method)
		return
	}

	if resp.StatusCode() != 200 {
		err = fmt.Errorf("received status code %d", resp.StatusCode())
		return
	}

	if err != nil {
		return
	}

	res = resp.Body()
	return
}
