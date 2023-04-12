package utils

import (
	res "botgpt/internal/enum"
)

type KnownError struct {
	message string
	code    res.ResponseCode
}

func (e *KnownError) Error() string {
	return e.message
}

func (e *KnownError) Code() res.ResponseCode {
	return e.code
}

func NewKnownError(code res.ResponseCode, message string) error {
	return &KnownError{
		message: message,
		code:    code,
	}
}
