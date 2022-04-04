package responses

import (
	"encoding/json"
	"net/http"
)

// Response is a simple interface to ensure custom responses return a
// standardized response body.
type Response interface {
	GetResponse() responsePayload
}

// responsePayload is the streamlined non-error client response body.
type responsePayload struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// WriteJSON encodes a `Response` into JSON format, to send to the client via
// the `http.ResponseWriter`.
func WriteJSON(w http.ResponseWriter, resp Response) error {
	return json.NewEncoder(w).Encode(resp.GetResponse())
}
