package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dvo-dev/go-get-started/pkg/customerrors"
	"github.com/dvo-dev/go-get-started/pkg/datastorage"
)

// DataStorageHandler is a wrapper struct for `DataStorage` implementations,
// providing REST handlers to the underlying storage solution.
type DataStorageHandler struct {
	storage datastorage.DataStorage
}

// Initialize initializes and returns a pointer to a `DataStorageHandler`,
// which will wrap around any `DataStorage` implementation, `storage`.
func (h DataStorageHandler) Initialize(storage datastorage.DataStorage) *DataStorageHandler {
	return &DataStorageHandler{
		storage: storage,
	}
}

// HandleClientRequest will parse and execute on any client requests intended
// for access to a `DataStorage`.
//
// Supported request methods are GET, POST, DELETE.
//
// To use, assign to your `http.Handler` with the desired URL pattern.
func (h *DataStorageHandler) HandleClientRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.Method {

		case http.MethodGet:
			err = h.retrieveData(w, r)

		case http.MethodPost:

		case http.MethodDelete:

		default:
			cErr := customerrors.ClientErrorBadMethod{
				RequestMethod: r.Method,
			}
			w.WriteHeader(cErr.StatusCode())
			err = json.NewEncoder(w).Encode(cErr.ClientErrorMsg())
		}

		if err != nil {
			log.Printf("DataStorageHandler - error writing response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *DataStorageHandler) retrieveData(w http.ResponseWriter, r *http.Request) error {
	// Parse request param
	dataKey := r.URL.Query().Get("name")
	log.Printf(
		"DataStorageHandler - attempting retrieval of data associated with the key %s",
		dataKey,
	)

	// Attempt to retrieve the data from our storage
	data, err := h.storage.RetrieveData(dataKey)
	if err != nil {
		log.Printf(
			"DataStorageHandler - failed to retrieve data with the key: %s\n\t%v",
			dataKey, err,
		)

		// Tell client if they used a bad key
		cErr := customerrors.ClientErrorMessage{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(cErr)
		return err
	}

	// Attempt to send retrieved data to the client
	w.WriteHeader(http.StatusFound)
	_, err = w.Write(data)
	return err
}
