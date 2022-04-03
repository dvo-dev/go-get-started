package responses

import "fmt"

// DataFound is the client response generator when data is successfully
// retrieved from `DataStorage`.
type DataFound struct {
	DataName string
	Data     []byte // TODO: add metadata?
}

func (d DataFound) GetResponse() responsePayload {
	return responsePayload{
		Status: "success",
		Message: fmt.Sprintf(
			"data with name: '%s' found",
			d.DataName,
		),
		Data: struct {
			Content string `json:"content"`
			Size    int    `json:"size"`
		}{
			Content: string(d.Data),
			Size:    len(d.Data),
		},
	}
}

// DataStored is the client response generator when data is successfully
// written to `DataStorage`.
type DataStored struct {
	DataName string
	Data     []byte
}

func (d DataStored) GetResponse() responsePayload {
	return responsePayload{
		Status: "success",
		Message: fmt.Sprintf(
			"data written to storage with name: '%s'",
			d.DataName,
		),
		Data: struct {
			Size int `json:"size"`
		}{
			Size: len(d.Data),
		},
	}
}

// DataDeleted is the client response generator when data is successfully
// deleted from `DataStorage`.
type DataDeleted struct {
	DataName string
}

func (d DataDeleted) GetResponse() responsePayload {
	return responsePayload{
		Status: "success",
		Message: fmt.Sprintf(
			"data with name: '%s' deleted",
			d.DataName,
		),
	}
}
