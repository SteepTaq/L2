package model

import (
	"encoding/json"
	"fmt"
)

type errorResp struct {
	Err string `json:"error"`
}

func ErrorResp(err string) errorResp {
	return errorResp{
		Err: err,
	}
}

func (e errorResp) ToJSON() []byte {
	return marshal("error response", e)
}

type resultResp struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ResultResp(message string) resultResp {
	return resultResp{
		Message: message,
	}
}

func ResultWithDataResp(message string, data any) resultResp {
	return resultResp{
		Message: message,
		Data:    data,
	}
}

func (e resultResp) ToJSON() []byte {
	return marshal("result response", e)
}

func marshal(name string, r any) []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic(fmt.Sprintf("Unable to unmarshal %s: %s", name, err))
	}

	return bytes
}
