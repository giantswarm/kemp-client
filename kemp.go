package kempclient

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Config struct {
	User     string
	Password string
	Endpoint string
}

type Client struct {
	user     string
	password string
	endpoint string
}

type ErrorResponse struct {
	Error string
}

type SuccessResponse struct {
	XMLName xml.Name     `xml:"Response"`
	Data    DataResponse `xml:"Success>Data"`
}

type DataResponse struct {
	Parameters []ParameterResponse `xml:",any"`
}

type ParameterResponse struct {
	XMLName xml.Name `xml:""`
	Value   string   `xml:",chardata"`
}

func NewClient(config Config) *Client {
	c := &Client{
		user:     config.User,
		password: config.Password,
		endpoint: config.Endpoint,
	}

	return c
}

func (c *Client) Get(param string) (string, error) {
	parameters := make(map[string]string)
	parameters["param"] = param

	result, err := c.Request("get", parameters)
	if err != nil {
		return "", err
	}
	return result[param], nil
}

func (c *Client) Set(param, value string) (string, error) {
	parameters := make(map[string]string)
	parameters["param"] = param

	result, err := c.Request("get", parameters)
	if err != nil {
		return "", err
	}

	parameters["value"] = value
	_, err = c.Request("set", parameters)
	if err != nil {
		return "", err
	}

	return result[param], nil
}

func (c *Client) Request(cmd string, parameters map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	params := url.Values{}
	for key, val := range parameters {
		params.Set(key, val)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?%s", c.endpoint, cmd, params.Encode()), nil)
	if err != nil {
		return result, err
	}

	req.SetBasicAuth(c.user, c.password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if res.StatusCode >= 400 {
		errorResponse := ErrorResponse{}
		err := ParseResponse(res.Body, &errorResponse)
		if err != nil {
			return result, err
		}

		return result, errors.New(errorResponse.Error)
	}

	successResponse := SuccessResponse{}
	err = ParseResponse(res.Body, &successResponse)

	for _, param := range successResponse.Data.Parameters {
		result[param.XMLName.Local] = param.Value
	}
	return result, nil
}
