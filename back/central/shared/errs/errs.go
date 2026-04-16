package errs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

type Error struct {
	Code   string   `json:"code"`
	Msg    string   `json:"message"`
	Detail any      `json:"detail,omitempty"`
	Status int      `json:"-"`
	Meta   Meta     `json:"-"`
	cause  error    `json:"-"`
	stack  []string `json:"-"`
	base   bool     `json:"-"`
}

type Meta map[string]any

func New(msg string) *Error {
	return &Error{
		Msg: msg,
	}
}

func Newf(msg string, a ...any) *Error {
	return &Error{
		Msg: fmt.Sprintf(msg, a...),
	}
}

var allBaseErrors []*Error

func Base(code string, status int, msg ...string) *Error {
	// get the package name where the error was created
	ms := ""
	if len(msg) > 0 {
		ms = msg[0]
	}
	e := &Error{
		Code:   code,
		Status: status,
		Msg:    ms,
		base:   true,
	}

	allBaseErrors = append(allBaseErrors, e)
	return e
}

func (e *Error) Clone() *Error {
	return &Error{
		Code:   e.Code,
		Msg:    e.Msg,
		Status: e.Status,
		Meta:   e.Meta,
		cause:  e.cause,
	}
}

func (e *Error) WithStatus(status int) *Error {
	if e.base {
		return e.Clone().WithStatus(status)
	}
	e.Status = status
	return e
}

func (e *Error) WithCode(code string) *Error {
	if e.base {
		return e.Clone().WithCode(code)
	}
	e.Code = code
	return e
}

func (e *Error) WithMsg(msg string) *Error {
	if e.base {
		return e.Clone().WithMsg(msg)
	}
	e.Msg = msg
	return e
}

func (e *Error) WithMsgf(format string, a ...any) *Error {
	if e.base {
		return e.Clone().WithMsgf(format, a...)
	}
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

func (e *Error) WithDetail(detail any) *Error {
	if e.base {
		return e.Clone().WithDetail(detail)
	}
	e.Detail = detail
	return e
}

func (e *Error) WithMeta(meta Meta) *Error {
	if e.base {
		return e.Clone().WithMeta(meta)
	}
	e.Meta = meta
	return e
}

func (e *Error) Wrap(err error) *Error {
	if e.base {
		return e.Clone().Wrap(err)
	}
	e.cause = err
	return e
}

func (e *Error) Freeze() *Error {
	e.base = true
	return e
}

func (e *Error) Error() string {
	if e.Code == "" {
		return e.Msg
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Is(target error) bool {
	if e == target {
		return true
	}

	if v, ok := target.(*Error); ok {
		if e.Code == v.Code && e.Msg == v.Msg {
			return true
		}
	}

	return errors.Is(e, target)
}

func (er *Error) MarshalZerologObject(e *zerolog.Event) {
	e.Str("msg", er.Msg)
	e.Str("code", er.Code)
	if er.cause != nil {
		e.AnErr("cause", er.cause)
	}
	if er.Detail != nil {
		e.Interface("detail", er.Detail)
	}
	if er.Meta != nil {
		e.Interface("meta", er.Meta)
	}
	if er.stack != nil {
		e.Strs("stack", er.stack)
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

type MultiErr []error

func (m MultiErr) Error() string {
	e := make([]string, len(m))
	for i, err := range m {
		e[i] = err.Error()
	}
	return strings.Join(e, ", ")
}
