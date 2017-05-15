package infrastructure

import (
	"errors"
)

type (
	TplName string

	ApiJsonError struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	IntentError struct {
		error
		Status int
		Code   int
	}
)

func (self *IntentError) Error() string {
	return self.error.Error()
}

func MakeIntentError(err interface{}, code ...int) *IntentError {
	if err == nil {
		return nil
	}
	var er error
	if e, ok := err.(error); ok {
		er = e
	} else if e, ok := err.(string); ok {
		er = errors.New(e)
	} else if e, ok := err.(*IntentError); ok {
		return e
	}

	if er == nil {
		return nil
	}

	status := 500
	errCode := 0
	if len(code) > 0 {
		if code[0] > 0 {
			status = 406
			errCode = code[0]
		}
	}
	return &IntentError{er, status, errCode}
}

var (
	tempFilePrefix  string
	encodeSecretKey string
	htmlTimeFormat  string
)

func init() {
	tempFilePrefix = ""
	encodeSecretKey = ""
	htmlTimeFormat = "RFC1123"
}

func InitInfrastructure(tempPrefix, encodeSecKey, timeFormat string) {
	tempFilePrefix = tempPrefix
	encodeSecretKey = encodeSecKey
	htmlTimeFormat = timeFormat
}
