package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// GetRequest executes a GET request to the desired `url`, adding any `params`
// to the raw query.
//
// It optionally can use a given http.Client, or will default to a standard
// client with a 10s timeout.
//
// Returns the parsed response or error if this function failed.
func GetRequest(url string, params *map[string]string, client *http.Client) (map[string]any, error) {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	// Construct request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Add query params
	if params != nil {
		query := req.URL.Query()
		for k, v := range *params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var parsed map[string]any
	json.Unmarshal([]byte(body), &parsed)

	return parsed, nil
}
