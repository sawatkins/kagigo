package kagi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/httpjamesm/kagigo/constants"
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

func (c *Client) SendRequest(method, path string, data interface{}, v any) (err error) {

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
		var res UniversalSummarizerResponse
		err := json.Unmarshal(resp.Body(), &res)
		if err != nil || len(res.Errors) == 0 {
			return fmt.Errorf("received status code %d with unparseable error response", resp.StatusCode())
		}
		errObj := res.Errors[0]
		return fmt.Errorf("received status code %d. error object: %v", resp.StatusCode(),
			fmt.Sprintf("[code: %d, msg: %s, ref: %v]", errObj.Code, errObj.Msg, errObj.Ref))
	}

	if err != nil {
		return
	}

	// get reader from body
	body := resp.Body()
	reader := bytes.NewReader(body)
	ioReader := io.Reader(reader)

	return decodeResponse(ioReader, v)
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	if result, ok := v.(*string); ok {
		return decodeString(body, result)
	}
	return json.NewDecoder(body).Decode(v)
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
