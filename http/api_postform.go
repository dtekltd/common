package http

import (
	"bytes"
	"net/url"
)

// PostForm sends a URL-encoded POST request and decodes the JSON response
// Parameters:
//   - url: The target URL
//   - data: Form data as key-value pairs
//   - out: Pointer to a struct that will receive the decoded JSON response
//   - headers: headers as key-value pairs
//
// Returns:
//   - error if any occurred (including non-2xx status codes)
func PostForm(url string, data url.Values, out any, headers map[string]string) error {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	in := bytes.NewBufferString(data.Encode())
	return httpRequest("POST", url, in, out, headers)
}
