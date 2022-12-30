package govisa

import "log"

type ErrorCode int

const (
	ErrBadRequest ErrorCode = 400
	ErrAuth       ErrorCode = 401
	ErrNotFound   ErrorCode = 404
	ErrInternal   ErrorCode = 500
)

type Error struct {
	E error     // error
	M string    // message
	C ErrorCode // code
}

func NewError() *Error {
	return &Error{}
}

func (e *Error) Error(err error) *Error {
	e.E = err
	return e
}

func (e *Error) Code(c ErrorCode) *Error {
	e.C = c
	return e
}

func (e *Error) Message(msg string) *Error {
	e.M = msg
	return e
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
