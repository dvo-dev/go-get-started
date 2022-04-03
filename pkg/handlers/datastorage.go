package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dvo-dev/go-get-started/pkg/customerrors"
	"github.com/dvo-dev/go-get-started/pkg/datastorage"
	"github.com/dvo-dev/go-get-started/pkg/responses"
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
			err = h.storeData(w, r)

		case http.MethodDelete:
			err = h.deleteData(w, r)
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
		"DataStorageHandler - attempting retrieval of data associated with the key: '%s'",
		dataKey,
	)

	// Attempt to retrieve the data from our storage
	data, err := h.storage.RetrieveData(dataKey)
	switch err {
	case nil:
		// Successfully retrieved associated data
		log.Printf(
			"DataStorageHandler - successfully retrieved data with key: '%s'",
			dataKey,
		)
		w.WriteHeader(http.StatusFound)

		// Attempt to write response message
		if rErr := responses.WriteJSON(w,
			responses.DataFound{
				DataName: dataKey,
				Data:     data,
			},
		); rErr != nil {
			log.Printf(
				"DataStorageHander - data retrieved but writing response failed: %v",
				err,
			)
		}
		return nil
	case customerrors.DataStorageNameNotFound{}:
		// Tell client if they used a bad key
		log.Printf(
			"DataStorageHandler - failed to retrieve data with the key: '%s'\n\t%v",
			dataKey, err,
		)

		cErr := customerrors.ClientErrorMessage{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(cErr)
		return err
	default:
		return err
	}
}

func (h *DataStorageHandler) storeData(w http.ResponseWriter, r *http.Request) error {
	// Parse request params
	err := r.ParseMultipartForm(2 << 20) // TODO: set constant somewhere
	if err != nil {
		log.Printf(
			"DataStorageHandler - failed to parse storage request: %v",
			err,
		)
		return err
	}

	// Get user defined name to associate with data
	name := strings.TrimSpace(r.PostFormValue("name"))
	if len(name) > 0 {
		log.Printf(
			"DataStorageHandler - attempting to write data to key: '%s'",
			name,
		)
	} else {
		log.Println(
			"DataStorageHandler - no name given by client, will use file name...",
		)
	}

	// Read user uploaded data
	var dataBuffer bytes.Buffer
	file, header, err := r.FormFile("data")
	if len(name) == 0 {
		name = header.Filename
	}
	if err != nil {
		log.Printf(
			"DataStorageHandler - failed to read request file with key: '%s'\n\t%v",
			name, err,
		)
		return err
	}
	defer file.Close()
	io.Copy(&dataBuffer, file)
	data := dataBuffer.Bytes()

	// Attempt to write the data to our storage
	err = h.storage.StoreData(name, data)
	switch err {
	// Other errors in the future
	case nil:
		log.Printf(
			"successfully wrote to storage data with key: '%s'",
			name,
		)
		w.WriteHeader(http.StatusCreated)

		// Attempt to write response message
		if rErr := responses.WriteJSON(w, responses.DataStored{
			DataName: name,
			Data:     data,
		}); rErr != nil {
			log.Printf(
				"DataStorageHander - data written but writing response failed: %v",
				err,
			)
		}
		return nil
	default:
		return err
	}
}

func (h *DataStorageHandler) deleteData(w http.ResponseWriter, r *http.Request) error {
	// Parse request param
	dataKey := r.URL.Query().Get("name")
	log.Printf(
		"DataStorageHandler - attempting deletion of data associated with the key: '%s'",
		dataKey,
	)

	// Attempt to delete the data fro our storage
	err := h.storage.DeleteData(dataKey)
	switch err {
	case nil:
		// Successfully deleted data associated with key
		log.Printf(
			"DataStorageHandler - successfully deleted data with the key: '%s'",
			dataKey,
		)
		w.WriteHeader(http.StatusOK)

		// Attempt to write response message
		if rErr := responses.WriteJSON(w,
			responses.DataDeleted{DataName: dataKey},
		); rErr != nil {
			log.Printf(
				"DataStorageHander - data deleted but writing response failed: %v",
				err,
			)
		}
		return nil
	case customerrors.DataStorageNameNotFound{}:
		// Tell client if they used a bad key
		log.Printf(
			"DataStorageHandler - failed to retrieve data with the key: '%s'\n\t%v",
			dataKey, err,
		)

		cErr := customerrors.ClientErrorMessage{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(cErr)
		return err
	default:
		return err
	}
}
