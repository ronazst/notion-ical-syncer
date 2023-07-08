package util

import "net/http"

type Error struct {
	code    int
	message string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Code() int {
	return e.code
}

func NewInternalError(message string) error {
	return &Error{
		code:    http.StatusInternalServerError,
		message: message,
	}
}
func NewUserInputError(message string) error {
	return &Error{
		code:    http.StatusBadRequest,
		message: message,
	}
}
