package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dvo-dev/go-get-started/pkg/customerrors"
	"github.com/dvo-dev/go-get-started/pkg/datastorage"
	"github.com/dvo-dev/go-get-started/pkg/responses"
	"github.com/dvo-dev/go-get-started/pkg/utils/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	resp, err := testClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	// Read response body
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Check error response
	rcvMsg := customerrors.ClientErrorMessage{}
	err = json.Unmarshal(data, &rcvMsg)
	require.NoError(t, err)

	expMsg := customerrors.ClientErrorBadMethod{
		RequestMethod: http.MethodPut,
	}.ClientErrorMsg()
	assert.Equal(t, expMsg.Error, rcvMsg.Error)
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

	// Generate and execute POST
	params := map[string]string{"name": testName}
	uploadData := map[string][]byte{"data": testData}
	resp, err := requests.PostRequest(
		testURL,
		"multipart/form-data",
		&params,
		&uploadData,
		nil,
	)
	require.NoError(t, err)

	// Check status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Read response body + assert
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	rcvMsg := responses.DataStored{}.GetResponse()
	err = json.Unmarshal(data, &rcvMsg)
	require.NoError(t, err)
	expMsg := responses.DataStored{
		DataName: testName,
		Data:     []byte("test data"),
	}.GetResponse()

	// Have to do this back and forth because Data is undefined type
	JSON, _ := json.Marshal(expMsg)
	_ = json.Unmarshal(JSON, &expMsg)
	assert.Equal(t, expMsg, rcvMsg)
}

func TestDataStorage_RetrieveData(t *testing.T) {
	// Init test server + client
	dsh := DataStorageHandler{}.Initialize(
		datastorage.MemStorage{}.Initialize(),
	)
	testServer := httptest.NewServer(dsh.HandleClientRequest())
	testURL := testServer.URL + "/datastorage"

	// Create GET request w/ params
	testName := "testname"
	testData := []byte("test data")
	err := dsh.storage.StoreData(testName, testData)
	require.NoError(t, err)

	// Make GET request
	params := map[string]string{"name": testName}
	resp, err := requests.GetRequest(testURL, &params, nil)
	require.NoError(t, err)

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read response body + assert
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	rcvMsg := responses.DataFound{}.GetResponse()
	err = json.Unmarshal(data, &rcvMsg)
	require.NoError(t, err)
	expMsg := responses.DataFound{
		DataName: testName,
		Data:     []byte("test data"),
	}.GetResponse()

	// Have to do this back and forth because Data is undefined type
	JSON, _ := json.Marshal(expMsg)
	_ = json.Unmarshal(JSON, &expMsg)
	assert.Equal(t, expMsg, rcvMsg)

	t.Run("nonexistent name", func(t *testing.T) {

		// Make GET request
		params["name"] = testName + "foo"
		resp, err := requests.GetRequest(testURL, &params, nil)
		require.NoError(t, err)

		// Check status code
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Read response body + assert
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		rcvMsg := customerrors.ClientErrorMessage{}
		err = json.Unmarshal(data, &rcvMsg)
		require.NoError(t, err)

		// Read response body + assert
		expMsg := customerrors.ClientErrorMessage{
			Error: customerrors.DataStorageNameNotFound{
				Name: testName + "foo",
			}.Error(),
		}
		assert.Equal(t, expMsg, rcvMsg)
	})
}

func TestDataStorage_DeleteData(t *testing.T) {
	// Init test server + client
	dsh := DataStorageHandler{}.Initialize(
		datastorage.MemStorage{}.Initialize(),
	)
	testServer := httptest.NewServer(dsh.HandleClientRequest())
	testURL := testServer.URL + "/datastorage"

	// Create DELETE request w/ params
	testName := "testname"
	testData := []byte("test data")
	err := dsh.storage.StoreData(testName, testData)
	require.NoError(t, err)

	// Make DELETE request
	params := map[string]string{"name": testName}
	resp, err := requests.CustomRequest(testURL, http.MethodDelete, &params, nil)
	require.NoError(t, err)

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read response body + assert
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	rcvMsg := responses.DataDeleted{}.GetResponse()
	err = json.Unmarshal(data, &rcvMsg)
	require.NoError(t, err)
	expMsg := responses.DataDeleted{
		DataName: testName,
	}.GetResponse()

	// Have to do this back and forth because Data is undefined type
	JSON, _ := json.Marshal(expMsg)
	_ = json.Unmarshal(JSON, &expMsg)
	assert.Equal(t, expMsg, rcvMsg)

	t.Run("nonexistent name", func(t *testing.T) {

		// Make DELETE request
		resp, err := requests.CustomRequest(testURL, http.MethodDelete, &params, nil)
		require.NoError(t, err)

		// Check status code
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Read response body + assert
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		rcvMsg := customerrors.ClientErrorMessage{}
		err = json.Unmarshal(data, &rcvMsg)
		require.NoError(t, err)

		expMsg := customerrors.ClientErrorMessage{
			Error: customerrors.DataStorageNameNotFound{
				Name: testName,
			}.Error(),
		}
		assert.Equal(t, expMsg, rcvMsg)
	})
}
