package requests

import (
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequests_GetRequest(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Method, http.MethodGet, "Request is not for GET Method")
			assert.Equal(t, r.URL.Query().Get("param1"), "foobar", "Param values do not match")
		})

		testServer := httptest.NewServer(getTestHandler)
		params := map[string]string{"param1": "foobar"}
		_, err := GetRequest(testServer.URL, &params, nil)
		require.NoError(t, err)
	})
	t.Run("nil params", func(t *testing.T) {
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Method, http.MethodGet, "Request is not for GET Method")
			err := r.ParseForm()
			require.NoError(t, err)
			assert.Empty(t, r.Form)
		})

		testServer := httptest.NewServer(getTestHandler)
		_, err := GetRequest(testServer.URL, nil, nil)
		require.NoError(t, err)
	})
}

func TestRequests_PostRequest(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost, "Request is not for POST Method")
		err := r.ParseForm()
		require.NoError(t, err)
		assert.Empty(t, r.Form)
	})

	testServer := httptest.NewServer(getTestHandler)
	params := map[string]string{"param1": "foobar"}
	_, err := PostRequest(testServer.URL, "application/json", &params, nil, nil)
	require.NoError(t, err)
}

func TestRequest_CustomRequest(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "custom", "Request is not for custom method")
		err := r.ParseForm()
		require.NoError(t, err)
		assert.NotEmpty(t, r.Form)
	})

	testServer := httptest.NewServer(getTestHandler)
	params := map[string]string{"param1": "foobar"}
	_, err := CustomRequest(testServer.URL, "custom", &params, nil)
	require.NoError(t, err)
}
