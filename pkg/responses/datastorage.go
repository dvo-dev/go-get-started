package responses

import "fmt"

type DataFound struct {
	DataName string
	Data     []byte // TODO: add metadata?
}

func (d DataFound) GetResponse() ResponsePayload {
	return ResponsePayload{
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

type DataStored struct {
	DataName string
	Data     []byte
}

func (d DataStored) GetResponse() ResponsePayload {
	return ResponsePayload{
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

type DataDeleted struct {
	DataName string
}

func (d DataDeleted) GetResponse() ResponsePayload {
	return ResponsePayload{
		Status: "success",
		Message: fmt.Sprintf(
			"data with name: '%s' deleted",
			d.DataName,
		),
	}
}
