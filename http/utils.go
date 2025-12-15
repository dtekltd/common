package http

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	netHttp "net/http"
)

var client *netHttp.Client

func httpClient() *netHttp.Client {
	if client == nil {
		client = &netHttp.Client{
			Timeout: 10 * time.Second, // Adjust as needed
		}
	}
	return client
}

func httpRequest(method string, url string, in io.Reader, out any, headers map[string]string) error {
	req, err := netHttp.NewRequest(method, url, in)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if len(headers) > 0 {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}

	// Send request
	resp, err := httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and decode response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("failed to decode JSON response: %w (body: %s)", err, string(body))
		}
	}

	return nil
}
