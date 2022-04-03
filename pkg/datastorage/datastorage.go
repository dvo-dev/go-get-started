package datastorage

// DataStorage is an interface to be satisified by any storage implementation,
// regardless of whether it is in memory, drive, etc.
type DataStorage interface {
	// RetrieveData allows user access to the DataStorage implementation,
	// retrieving data associated with the given `name`.
	//
	// Returns the data in byte forms if found, error elsewise.
	//
	// This method is thread safe - multiple reads to the same data is
	// allowed but writers will be blocked until their turn.
	RetrieveData(name string) ([]byte, error)

	// StoreData allows users to store given `data` (bytes) with an associated
	// `name`.
	//
	// This method will overwrite existing data if an already existing `name` is
	// given - returns an error if write fails.
	//
	// Callers may assume this method is thread safe.
	StoreData(name string, data []byte) error

	// DeleteData allows uers to delete data associated with a given `name`.
	//
	// This method returns an error if there is no data associated with `name`,
	// or if the deletion fails.
	//
	// Callers may assume this method is thread safe.
	DeleteData(name string) error
}
