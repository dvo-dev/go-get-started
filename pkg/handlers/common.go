package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dvo-dev/go-get-started/pkg/customerrors"
)

// RecoveryWrapper is a function wrapper for the actual intended route handling
// functions.
//
// This acts as essentially middleware to catch panics by the wrapped handler,
// and recover gracefully while logging the error.
//
// This function takes advantage of Go's `recover()` function which can be
// utilized by executing the wrapper handler (`h`) within a Goroutine.
func RecoveryWrapper(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		defer func() {
			if r := recover(); nil != r {
				log.Printf("Error occurred: %v\n, recovered", r)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// HealthStatus is a simple struct to construct the response payload of
// `HandleHealth`.
type HealthStatus struct {
	Status string `json:"status"`
}

// HandleHealth is to essentially act as a "heartbeat" for its server.
//
// Other services can validate the associated server is live by receiving a
// simple http.StatusOK (`200`) with a status field of "healthy" from this
// function.
func HandleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {

		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(HealthStatus{
				Status: "healthy",
			})

		default:
			cErr := customerrors.ClientErrorBadMethod{
				RequestMethod: r.Method,
			}
			w.WriteHeader(cErr.StatusCode())
			err = json.NewEncoder(w).Encode(cErr.ClientErrorMsg())
		}

		if err != nil {
			log.Printf("HandleHealth - error writing response: %v", err)
		}
	}
}
