package models

import "encoding/json"

type Error struct {
	message string
	code    int64
}

func MakeErrorResponse(message string, code int64) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Code() int64 {
	return e.code
}

func (e *Error) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
