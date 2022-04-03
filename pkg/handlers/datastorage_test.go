package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/dvo-dev/go-get-started/pkg/customerrors"
	"github.com/dvo-dev/go-get-started/pkg/datastorage"
	"github.com/dvo-dev/go-get-started/pkg/responses"
)

// TODO: all these tests need to be rewritten once test helpers are in

func TestDataStorage_BadMethod(t *testing.T) {
	// Init test server + client
	dsh := DataStorageHandler{}.Initialize(
		datastorage.MemStorage{}.Initialize(),
	)
	testServer := httptest.NewServer(dsh.HandleClientRequest())
	testURL := testServer.URL + "/datastorage"
	testClient := http.DefaultClient

	// Create a PUT request
	req, err := http.NewRequest(http.MethodPut, testURL, strings.NewReader(""))
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}
	resp, err := testClient.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected error code: %d but got: %d",
			http.StatusMethodNotAllowed,
			resp.StatusCode,
		)
	}

	// Read response body
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Check error response
	rcvMsg := customerrors.ClientErrorMessage{}
	err = json.Unmarshal(data, &rcvMsg)
	if err != nil {
		t.Error(
			"incorrect response format",
		)
	}
	expMsg := customerrors.ClientErrorBadMethod{
		RequestMethod: http.MethodPut,
	}.ClientErrorMsg()
	if expMsg.Error != rcvMsg.Error {
		t.Errorf(
			"expected response error: %s but got: %s",
			expMsg.Error, rcvMsg.Error,
		)
	}
}

func TestDataStorage_StoreData(t *testing.T) {
	// Init test server + client
	dsh := DataStorageHandler{}.Initialize(
		datastorage.MemStorage{}.Initialize(),
	)
	testServer := httptest.NewServer(dsh.HandleClientRequest())
	testURL := testServer.URL + "/datastorage"
	testClient := http.DefaultClient

	testName := "test name"
	testData := []byte("test data")

	// Multipart form, assign name
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormField("name")
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}
	_, err = io.Copy(fw, strings.NewReader(testName))
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Assign data
	fw, err = writer.CreateFormFile("data", "test.log")
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}
	_, err = io.Copy(fw, bytes.NewReader(testData))
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Close multipart form + make POST request
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, testURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Check status code
	if resp.StatusCode != http.StatusCreated {
		t.Errorf(
			"expected status code: %d but got: %d",
			http.StatusCreated,
			resp.StatusCode,
		)
	}

	// Read response body + assert
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}
	recMsg := responses.ResponsePayload{}
	err = json.Unmarshal(data, &recMsg)
	if err != nil {
		t.Error(
			"incorrect response format",
		)
	}
	expMsg := responses.DataStored{
		DataName: "test",
		Data:     []byte("test data"),
	}.GetResponse()
	// Have to do this back and forth because Data is undefined type
	JSON, _ := json.Marshal(expMsg)
	_ = json.Unmarshal(JSON, &expMsg)
	if !reflect.DeepEqual(expMsg, recMsg) {
		t.Errorf(
			"expected response: %+v but got: %+v",
			expMsg, recMsg,
		)
	}
}