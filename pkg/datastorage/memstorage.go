package datastorage

import (
	"sync"

	"github.com/dvo-dev/go-get-started/pkg/customerrors"
)

// MemStorage is an in-memory storage solution that implements the `DataStorage`
// interface.
//
// MemStorage stores data in the form of `[]byte` in map, with user defined
// `string` keys.
type MemStorage struct {
	data map[string][]byte
	rwMu *sync.RWMutex
}

// InitializeMemStorage initializes and returns a pointer to a clean
// `MemStorage`.
func (ms MemStorage) Initialize() *MemStorage {
	return &MemStorage{
		data: make(map[string][]byte),
		rwMu: &sync.RWMutex{},
	}
}

// RetrieveData checks the `MemStorage` for data associated with a given `name`.
//
// If the `name` is found, returns the data (`[]byte`), elsewise returns error.
//
// This method is thread safe.
func (ms *MemStorage) RetrieveData(name string) ([]byte, error) {
	// Acquire a read lock
	ms.rwMu.RLock()
	defer ms.rwMu.RUnlock()

	// Check if data exists associated with the given name
	data, found := ms.data[name]
	if !found {
		return []byte{}, customerrors.DataStorageNameNotFound{
			Name: name,
		}
	}

	// Return data with no error elsewise
	return data, nil
}

// StoreData writes data `[]byte` to the `MemStorage`, mapping it to the given
// `name`.
//
// If the `name` already exists, it will overwrite the previous data values.
//
// Returns an error if writing fails.
//
// This method is thread safe.
func (ms *MemStorage) StoreData(name string, data []byte) error {
	// Acquire a complete mutual exclusion lock
	ms.rwMu.Lock()
	defer ms.rwMu.Unlock()

	// Write the data and map to the assigned name
	ms.data[name] = data

	// TODO: add error handling for any future features that may require it,
	// e.g. memory capacity, etc.
	return nil
}

// DeleteData removes data in the `MemStorage` associated with the given `name`.
//
// If the `name` does not exist, an error will be returned.
//
// This method is thread safe.
func (ms *MemStorage) DeleteData(name string) error {
	// Acquire a complete mutual exclusion lock
	ms.rwMu.Lock()
	defer ms.rwMu.Unlock()

	// Check if data exists associated with the given name
	_, found := ms.data[name]
	if !found {
		return customerrors.DataStorageNameNotFound{
			Name: name,
		}
	}

	// Delete the data and return no error
	delete(ms.data, name)
	return nil
}
