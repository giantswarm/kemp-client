package kempclient

import (
	"encoding/xml"
	"fmt"
	"io"

	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
)

type ErrorResponse struct {
	Debug string `xml:",innerxml"`
	Error string
}

func (c *Client) parseResponse(reader io.Reader, result interface{}) error {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader
	return decoder.Decode(result)
}

func (c *Client) parseError(code int, reader io.Reader) error {
	errorResponse := ErrorResponse{}
	err := c.parseResponse(reader, &errorResponse)
	if err != nil {
		return err
	}
	if c.debug {
		fmt.Println("DEBUG:", errorResponse.Debug)
	}

	return newError(code, errorResponse.Error)
}

func (c *Client) parseSuccess(reader io.Reader, data interface{}) error {
	err := c.parseResponse(reader, data)
	if err != nil {
		return err
	}

	return nil
}
