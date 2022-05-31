package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

func IsJSON(contentType string) bool {
	return contentType == "application/json" || contentType == "text/json"
}

func IsXML(contentType string) bool {
	return contentType == "application/xml" || contentType == "text/xml"
}

func Decode(contentType string, reader io.Reader, msg interface{}) error {
	if IsJSON(contentType) {
		return json.NewDecoder(reader).Decode(msg)
	} else if IsXML(contentType) {
		return xml.NewDecoder(reader).Decode(msg)
	}

	return fmt.Errorf("unsupport content type: %s", contentType)
}
