package http

import (
	"bytes"
	"encoding/json"
	"io"
)

func Get(url string, out any, headers map[string]string) error {
	return Request("GET", url, nil, out, headers)
}

func Post(url string, data any, out any, headers map[string]string) error {
	return Request("POST", url, data, out, headers)
}

func Put(url string, data any, out any, headers map[string]string) error {
	return Request("PUT", url, data, out, headers)
}

func Delete(url string, out any, headers map[string]string) error {
	return Request("DELETE", url, nil, out, headers)
}

func Request(method string, url string, data any, out any, headers map[string]string) error {
	var in io.Reader

	if data != nil {
		var raw []byte
		raw, err := json.Marshal(data)
		in = bytes.NewBuffer(raw)
		if err != nil {
			return err
		}
	}

	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/json"

	return httpRequest(method, url, in, out, headers)
}
