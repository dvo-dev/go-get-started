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
			assert.Equal(t, http.MethodGet, r.Method, "Request is not for GET Method")
			err := r.ParseForm()
			require.NoError(t, err)
			assert.Equal(t, r.URL.Query().Get("param1"), "foobar", "Param values do not match")
		})

		testServer := httptest.NewServer(getTestHandler)
		params := map[string]string{"param1": "foobar"}
		_, err := GetRequest(testServer.URL, &params, nil)
		require.NoError(t, err)
	})
	t.Run("nil params", func(t *testing.T) {
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method, "Request is not for GET Method")
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
	t.Run("application/json", func(t *testing.T) {
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method, "Request is not for POST Method")
			err := r.ParseForm()
			require.NoError(t, err)
			assert.Empty(t, r.Form)
		})

		testServer := httptest.NewServer(getTestHandler)
		params := map[string]string{"param1": "foobar"}
		_, err := PostRequest(testServer.URL, "application/json", &params, nil, nil)
		require.NoError(t, err)
	})
	t.Run("no params application/json", func(t *testing.T) {
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method, "Request is not for POST Method")
			err := r.ParseForm()
			require.NoError(t, err)
			assert.Empty(t, r.Form)
		})

		testServer := httptest.NewServer(getTestHandler)
		params := map[string]string{}
		_, err := PostRequest(testServer.URL, "application/json", &params, nil, nil)
		if assert.Error(t, err) {
			assert.Equal(t, "no parameters provided", err.Error())
		}
	})
	t.Run("default case", func(t *testing.T) {
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})

		testServer := httptest.NewServer(getTestHandler)
		params := map[string]string{}
		_, err := PostRequest(testServer.URL, "", &params, nil, nil)
		if assert.Error(t, err) {
			assert.Equal(t, "unsupported content type", err.Error())
		}
	})
	t.Run("application/x-www-form-urlencoded", func(t *testing.T) {
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method, "Request is not for POST Method")
			err := r.ParseForm()
			require.NoError(t, err)
		})

		testServer := httptest.NewServer(getTestHandler)
		params := map[string]string{}
		_, err := PostRequest(testServer.URL, "application/x-www-form-urlencoded", &params, nil, nil)
		require.NoError(t, err)
	})
	t.Run("multipart/form-data", func(t *testing.T) {
		// TODO: Test for upload data
		getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method, "Request is not for POST Method")
			err := r.ParseForm()
			require.NoError(t, err)
			assert.Empty(t, r.Form)
		})

		testServer := httptest.NewServer(getTestHandler)
		params := map[string]string{"param1": "foobar"}
		value := []byte("value")
		uploadData := map[string][]byte{"file": value}
		_, err := PostRequest(testServer.URL, "multipart/form-data", &params, &uploadData, nil)
		require.NoError(t, err)
	})
}

func TestRequest_CustomRequest(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "custom", r.Method, "Request is not for custom method")
		err := r.ParseForm()
		require.NoError(t, err)
		assert.Equal(t, "foobar", r.URL.Query().Get("param1"), "Param values do not match")
	})

	testServer := httptest.NewServer(getTestHandler)
	params := map[string]string{"param1": "foobar"}
	_, err := CustomRequest(testServer.URL, "custom", &params, nil)
	require.NoError(t, err)
}
