package me

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/krave-n/deeperror"
)

const stackFrames = 2

func Err(err error, msg string) error {
	return deeperror.NewS(rand.Int63(), msg, err, stackFrames)
}

func NewErr(msg string) error {
	return deeperror.NewS(rand.Int63(), msg, nil, stackFrames)
}

func GetErrorMessage(e interface{}) string {
	msg := ""

	if e == nil {
		return msg
	} else if err, ok := e.(error); ok {
		msg = err.Error()
	} else if str, ok := e.(string); ok {
		msg = str
	} else if _, ok := e.(runtime.Error); ok {
		msg = err.Error()
	} else {
		msg = ""
	}
	return msg
}

// StackTrace get the stack trace, if inside of a panic recover, will return the stack that called the panic
// r the data returned from the panic
func StackTrace(er interface{}) string {
	var panicMessage string
	buffer := make([]byte, 4096)
	runtime.Stack(buffer, true)

	panicMessage = GetErrorMessage(er)
	stackTrace := fmt.Sprintf("%s - callstack: %s", panicMessage, buffer)
	return stackTrace
}
