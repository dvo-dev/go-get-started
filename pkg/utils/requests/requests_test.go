package requests

import (
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Full test coverage for GetRequest, PostRequest
// CustomRequest can be covered using two different http methods (e.g. PUT, DELETE)
// Test can be housed in pkg/utils/requests/requests_test.go
// Separate test functions for each wrapper tested, use sub tests if needed

func TestRequests_GetRequest(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet, "Request is not for GET Method")
		assert.Equal(t, r.URL.Query().Get("param1"), "foobar", "Param values do not match")
	})

	testServer := httptest.NewServer(getTestHandler)
	params := map[string]string{"param1": "foobar"}
	_, err := GetRequest(testServer.URL, &params, nil)
	require.NoError(t, err)
}

func TestRequests_GetRequest_NilParams(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet, "Request is not for GET Method")
		r.ParseForm()
		assert.Empty(t, r.Form)
	})

	testServer := httptest.NewServer(getTestHandler)
	_, err := GetRequest(testServer.URL, nil, nil)
	require.NoError(t, err)
}

func TestRequests_PostRequest(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost, "Request is not for POST Method")
		r.ParseForm()
		assert.Empty(t, r.Form)
		// Check params, header. etc.
	})

	testServer := httptest.NewServer(getTestHandler)
	params := map[string]string{"param1": "foobar"}
	_, err := PostRequest(testServer.URL, "application/json", &params, nil, nil)
	require.NoError(t, err)
}

func TestRequest_CustomRequest(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "custom", "Request is not for custom method")
		r.ParseForm()
		assert.NotEmpty(t, r.Form)
		// Check params, header. etc.
	})

	testServer := httptest.NewServer(getTestHandler)
	params := map[string]string{"param1": "foobar"}
	_, err := CustomRequest(testServer.URL, "custom", &params, nil)
	require.NoError(t, err)
}
