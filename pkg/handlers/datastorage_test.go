package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/dvo-dev/go-get-started/pkg/customerrors"
	"github.com/dvo-dev/go-get-started/pkg/datastorage"
	"github.com/dvo-dev/go-get-started/pkg/responses"
	"github.com/dvo-dev/go-get-started/pkg/utils/requests"
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
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}
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

	testName := "test name"
	testData := []byte("test data")

	params := map[string]string{"name": testName}
	uploadData := map[string][]byte{"data": testData}
	resp, err := requests.PostRequest(
		testURL,
		"multipart/form-data",
		&params,
		&uploadData,
		nil,
	)
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
	rcvMsg := responses.DataStored{}.GetResponse()
	err = json.Unmarshal(data, &rcvMsg)
	if err != nil {
		t.Error(
			"incorrect response format",
		)
	}
	expMsg := responses.DataStored{
		DataName: testName,
		Data:     []byte("test data"),
	}.GetResponse()
	// Have to do this back and forth because Data is undefined type
	JSON, _ := json.Marshal(expMsg)
	_ = json.Unmarshal(JSON, &expMsg)
	if !reflect.DeepEqual(expMsg, rcvMsg) {
		t.Errorf(
			"expected response: %+v but got: %+v",
			expMsg, rcvMsg,
		)
	}
}

func TestDataStorage_RetrieveData(t *testing.T) {
	// Init test server + client
	dsh := DataStorageHandler{}.Initialize(
		datastorage.MemStorage{}.Initialize(),
	)
	testServer := httptest.NewServer(dsh.HandleClientRequest())
	testURL := testServer.URL + "/datastorage"

	testName := "testname"
	testData := []byte("test data")
	err := dsh.storage.StoreData(testName, testData)
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Make GET request
	params := map[string]string{"name": testName}
	resp, err := requests.GetRequest(testURL, &params, nil)
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
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
	rcvMsg := responses.DataFound{}.GetResponse()
	err = json.Unmarshal(data, &rcvMsg)
	if err != nil {
		t.Error(
			"incorrect response format",
		)
	}
	expMsg := responses.DataFound{
		DataName: testName,
		Data:     []byte("test data"),
	}.GetResponse()
	// Have to do this back and forth because Data is undefined type
	JSON, _ := json.Marshal(expMsg)
	_ = json.Unmarshal(JSON, &expMsg)
	if !reflect.DeepEqual(expMsg, rcvMsg) {
		t.Errorf(
			"expected response: %+v but got: %+v",
			expMsg, rcvMsg,
		)
	}

	t.Run("nonexistent name", func(t *testing.T) {

		// Make GET request
		params["name"] = testName + "foo"
		resp, err := requests.GetRequest(testURL, &params, nil)
		if err != nil {
			t.Fatalf(
				"unexpected error occurred: %v",
				err,
			)
		}

		// Check status code
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf(
				"expected status code: %d but got: %d",
				http.StatusNotFound,
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
		rcvMsg := customerrors.ClientErrorMessage{}
		err = json.Unmarshal(data, &rcvMsg)
		if err != nil {
			t.Error(
				"incorrect response format",
			)
		}
		expMsg := customerrors.ClientErrorMessage{
			Error: customerrors.DataStorageNameNotFound{
				Name: testName + "foo",
			}.Error(),
		}
		if expMsg != rcvMsg {
			t.Errorf(
				"expected response: %+v but got: %+v",
				expMsg, rcvMsg,
			)
		}
	})
}

func TestDataStorage_DeleteData(t *testing.T) {
	// Init test server + client
	dsh := DataStorageHandler{}.Initialize(
		datastorage.MemStorage{}.Initialize(),
	)
	testServer := httptest.NewServer(dsh.HandleClientRequest())
	testURL := testServer.URL + "/datastorage"

	testName := "testname"
	testData := []byte("test data")
	err := dsh.storage.StoreData(testName, testData)
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Make DELETE request
	params := map[string]string{"name": testName}
	resp, err := requests.CustomRequest(testURL, http.MethodDelete, &params, nil)
	if err != nil {
		t.Fatalf(
			"unexpected error occurred: %v",
			err,
		)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
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
	rcvMsg := responses.DataDeleted{}.GetResponse()
	err = json.Unmarshal(data, &rcvMsg)
	if err != nil {
		t.Error(
			"incorrect response format",
		)
	}
	expMsg := responses.DataDeleted{
		DataName: testName,
	}.GetResponse()
	// Have to do this back and forth because Data is undefined type
	JSON, _ := json.Marshal(expMsg)
	_ = json.Unmarshal(JSON, &expMsg)
	if !reflect.DeepEqual(expMsg, rcvMsg) {
		t.Errorf(
			"expected response: %+v but got: %+v",
			expMsg, rcvMsg,
		)
	}

	t.Run("nonexistent name", func(t *testing.T) {

		// Make DELETE request
		resp, err := requests.CustomRequest(testURL, http.MethodDelete, &params, nil)
		if err != nil {
			t.Fatalf(
				"unexpected error occurred: %v",
				err,
			)
		}

		// Check status code
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf(
				"expected status code: %d but got: %d",
				http.StatusNotFound,
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
		rcvMsg := customerrors.ClientErrorMessage{}
		err = json.Unmarshal(data, &rcvMsg)
		if err != nil {
			t.Error(
				"incorrect response format",
			)
		}
		expMsg := customerrors.ClientErrorMessage{
			Error: customerrors.DataStorageNameNotFound{
				Name: testName,
			}.Error(),
		}
		if expMsg != rcvMsg {
			t.Errorf(
				"expected response: %+v but got: %+v",
				expMsg, rcvMsg,
			)
		}
	})
}
