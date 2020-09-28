package hc_errors

import (
	"encoding/json"
	"google.golang.org/grpc/status"
)

type JsonError struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

func (e JsonError) Error() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

func NewGeneralError(err error) *JsonError {
	return &JsonError{
		Type:        GENERAL_ERROR_TYPE,
		Description: err.Error(),
	}
}

func UnwrapJsonErrorFromRPCError(err error) (*JsonError, bool) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}

	var jError JsonError
	e := json.Unmarshal([]byte(st.Message()), &jError)
	if e != nil {
		return nil, false
	}
	return &jError, true
}

func (e *JsonError) WithData(key string, val interface{}) *JsonError {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data[key] = val
	return e
}
