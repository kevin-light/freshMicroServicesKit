package util

import "fmt"

type MyError struct {
	Code    int
	Message string
}

func NewMyError(code int, msg string) error {
	return &MyError{Code: code, Message: msg}
}

func (e *MyError) Error() string {
	fmt.Println(e.Message, "---01")
	return e.Message
}
