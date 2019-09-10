package errors

/**
*    Copyright (C) 2017 Ethan Frey
**/

import (
	"fmt"

	"github.com/pkg/errors"

	wrsp "github.com/tepleton/wrsp/types"
)

const defaultErrCode = wrsp.CodeType_InternalError

type stackTracer interface {
	error
	StackTrace() errors.StackTrace
}

type TMError interface {
	stackTracer
	ErrorCode() wrsp.CodeType
	Message() string
}

type tmerror struct {
	stackTracer
	code wrsp.CodeType
	msg  string
}

func (t tmerror) ErrorCode() wrsp.CodeType {
	return t.code
}

func (t tmerror) Message() string {
	return t.msg
}

// Format handles "%+v" to expose the full stack trace
// concept from pkg/errors
func (t tmerror) Format(s fmt.State, verb rune) {
	// special case also show all info
	if verb == 'v' && s.Flag('+') {
		fmt.Fprintf(s, "%+v\n", t.stackTracer)
	}
	// always print the normal error
	fmt.Fprintf(s, "(%d) %s\n", t.code, t.msg)
}

// Result converts any error into a wrsp.Result, preserving as much info
// as possible if it was already a TMError
func Result(err error) wrsp.Result {
	tm := Wrap(err)
	return wrsp.Result{
		Code: tm.ErrorCode(),
		Log:  tm.Message(),
	}
}

// Wrap safely takes any error and promotes it to a TMError
func Wrap(err error) TMError {
	// nil or TMError are no-ops
	if err == nil {
		return nil
	}
	// and check for noop
	tm, ok := err.(TMError)
	if ok {
		return tm
	}

	return WithCode(err, defaultErrCode)
}

// WithCode adds a stacktrace if necessary and sets the code and msg,
// overriding the state if err was already TMError
func WithCode(err error, code wrsp.CodeType) TMError {
	// add a stack only if not present
	st, ok := err.(stackTracer)
	if !ok {
		st = errors.WithStack(err).(stackTracer)
	}
	// and then wrap it with TMError info
	return tmerror{
		stackTracer: st,
		code:        code,
		msg:         err.Error(),
	}
}

// New adds a stacktrace if necessary and sets the code and msg,
// overriding the state if err was already TMError
func New(msg string, code wrsp.CodeType) TMError {
	// create a new error with stack trace and attach a code
	st := errors.New(msg).(stackTracer)
	return tmerror{
		stackTracer: st,
		code:        code,
		msg:         msg,
	}
}
