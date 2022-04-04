package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GetRequest executes a GET request to the desired `url`, adding any `params`
// to the raw query.
//
// It optionally can use a given http.Client, or will default to a standard
// client with a 10s timeout.
//
// Returns the parsed response or error if this function failed.
func GetRequest(
	address string,
	params *map[string]string,
	client *http.Client,
) (map[string]any, error) {
	// Init client if none given
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	// Construct request
	req, err := http.NewRequest(http.MethodGet, address, nil)
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
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var parsed map[string]any
	err = json.Unmarshal([]byte(body), &parsed)

	return parsed, err
}

func PostRequest(
	address string,
	contentType string,
	params *map[string]string,
	files *map[string][]byte, // TODO: accomadate direct file names
	client *http.Client,
) (map[string]any, error) {
	// Init client if none given
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	// Construct request
	var (
		req *http.Request
		err error
	)
	switch contentType {
	case "application/x-www-form-urlencoded":
		data := url.Values{}
		if params != nil {
			for k, v := range *params {
				data.Set(k, v)
			}
		}
		req, err = http.NewRequest(http.MethodPost, address, strings.NewReader(data.Encode()))
		req.Header.Set("Content-type", contentType)
		req.Header.Add("Content-length", strconv.Itoa(len(data.Encode())))
		if err != nil {
			return nil, err
		}

	case "multipart/form-data":
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Write additional params
		if params != nil {
			for k, v := range *params {
				fw, err := writer.CreateFormField(k)
				if err != nil {
					return nil, err
				}
				_, err = io.Copy(fw, strings.NewReader(v))
				if err != nil {
					return nil, err
				}
			}
		}

		// Write files
		if files != nil {
			for k, f := range *files {
				fw, err := writer.CreateFormFile(k, k)
				if err != nil {
					return nil, err
				}
				_, err = io.Copy(fw, bytes.NewReader(f))
				if err != nil {
					return nil, err
				}
			}
		}

		// Create the request
		writer.Close()
		req, err = http.NewRequest(http.MethodPost, address, bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-type", writer.FormDataContentType())

	case "application/json":
		// Marshal the params
		if params == nil || len(*params) == 0 {
			return nil, errors.New("no parameters provided")
		}
		jsonParams, err := json.Marshal(*params)
		if err != nil {
			return nil, err
		}

		// Create request
		req, err = http.NewRequest(http.MethodPost, address, bytes.NewBuffer(jsonParams))
		req.Header.Set("Content-type", contentType)

	default:
		return nil, errors.New("unsupported content type")
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var parsed map[string]any
	err = json.Unmarshal([]byte(body), &parsed)

	return parsed, err
}
