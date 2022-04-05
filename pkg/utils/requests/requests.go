package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GetRequest executes a GET request to the desired `targetURL`, adding any
// `params` to the raw query.
//
// It optionally can use a given http.Client, or will default to a standard
// client with a 10s timeout.
//
// Returns the raw response or error if this function failed.
// TODO: support authentication methods
func GetRequest(
	targetURL string,
	params *map[string]string,
	client *http.Client,
) (*http.Response, error) {
	// Init client if none given
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	// Add query params if applicable
	var queryParams *url.Values
	if params != nil && len(*params) != 0 {
		queryParams = &url.Values{}
		for k, v := range *params {
			queryParams.Add(k, v)
		}
	}

	// Execute the request
	resp, err := client.Get(targetURL + func() string {
		if queryParams != nil {
			return "?" + (*queryParams).Encode()
		} else {
			return ""
		}
	}())

	return resp, err
}

// PostRequest executes a POST request to the desired `targetURL`, supporting
// 3 `contentType`:
// application/x-www-form-urlencoded
// multipart/form-data
// application/json
//
// Pass any form parameters as `params`, and file/data to POST as `uploadData`.
//
// It optionally can use a given http.Client, or will default to a standard
// client with a 10s timeout.
//
// Returns the raw response or error if this function failed.
// TODO: support authentication methods
func PostRequest(
	targetURL string,
	contentType string,
	params *map[string]string,
	uploadData *map[string][]byte, // TODO: accomadate direct file names
	client *http.Client,
) (*http.Response, error) {
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
		if params != nil && len(*params) > 0 {
			for k, v := range *params {
				data.Set(k, v)
			}
		}
		req, err = http.NewRequest(http.MethodPost, targetURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-type", contentType)
		req.Header.Add("Content-length", strconv.Itoa(len(data.Encode())))
		if err != nil {
			return nil, err
		}

	case "multipart/form-data":
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Write additional params
		if params != nil && len(*params) > 0 {
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

		// Write uploadData
		if uploadData != nil && len(*uploadData) > 0 {
			for k, f := range *uploadData {
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
		req, err = http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body.Bytes()))
		if err != nil {
			return nil, err
		}
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
		req, err = http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(jsonParams))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-type", contentType)

	default:
		return nil, errors.New("unsupported content type")
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
