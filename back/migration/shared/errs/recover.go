package errs

import (
	"fmt"
	"runtime"

	"github.com/secamc93/probability/back/migration/shared/log"
)

// Recover recovers from a panic if any.
// can be passed an error pointer to assign the error to.
// if no error is passed, the error is logged.
func Recover(log log.ILogger, err ...*error) {

	if r := recover(); r != nil {
		stack := GetErrorStack()
		msg := ""
		if e, ok := r.(error); ok {
			msg = e.Error()
		} else {
			msg = fmt.Sprintf("%v", r)
		}

		e := &Error{
			Msg:   msg,
			Code:  "unknown",
			stack: stack,
		}

		if len(err) > 0 {
			*err[0] = e
		} else {
			log.Error().Err(e).Msg("recovered from panic")
		}
	}
}

func GetErrorStack() []string {
	stack := make([]string, 0)
	// Get the size of the call stack
	const size = 32
	var pcs [size]uintptr
	n := runtime.Callers(3, pcs[:])

	// Get information about each function in the call stack
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		stack = append(stack, fmt.Sprintf("%s:%d", frame.File, frame.Line))
	}

	return stack
}
