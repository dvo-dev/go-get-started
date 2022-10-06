package customerrors

import "fmt"

// DataStorageNameNotFound is an `error` that is associated with a failed key
// access to a `DataStorage` - there is no stored data associated with the given
// `Name`.
type DataStorageNameNotFound struct {
	Name string
}

func (e DataStorageNameNotFound) Error() string {
	return fmt.Sprintf(
		"attempted to access data associated with name: %s - not found",
		e.Name,
	)
}
