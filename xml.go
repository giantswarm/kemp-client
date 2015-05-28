package kempclient

import (
	"encoding/xml"
	"io"

	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
)

func ParseResponse(reader io.Reader, result interface{}) error {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader
	return decoder.Decode(result)
}
