package response

import (
	"encoding/json"
)

type Responsable interface {
	GetCode() int
	GetResponse() string
}

// Success
type Success struct {
	Code int `json:"status_code"`
}

func (s Success) GetCode() int {
	return s.Code
}

func (s Success) GetResponse() string {

	b, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return string(b)
}

// Error
type Error struct {
	Code    int    `json:"status_code"`
	Message string `json:"error_message"`
}

func (e Error) GetResponse() string {

	b, err := json.Marshal(e)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func (e Error) GetCode() int {
	return e.Code
}
