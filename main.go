package errors

import (
	"fmt"
	"os"
	"strings"
)

var (
	Separator = ": "
)

type Error struct {
	error
	next *Error
	len  int
}

func (err *Error) list() []error {
	if err.len == 0 {
		return []error{}
	}
	res := []error{}
	current := err
	if current.error != nil {
		res = append(res, current.error)
	}
	for current.next != nil {
		res = append(res, current.next.error)
		current = current.next
	}
	return res
}

func (err *Error) merge() error {
	if err.len == 0 {
		return nil
	}
	res := []string{}
	for _, e := range err.list() {
		res = append(res, e.Error())
	}
	return fmt.Errorf(strings.Join(res, Separator))
}

func (err *Error) Append(e error) *Error {

	if err.error == nil {
		err.error = e
		return err
	}

	ce := err
	for ce.next != nil {
		ce = ce.next
	}
	ce.next = New(e)
	err.len++
	return err
}

func (err *Error) Appendf(msg string, args ...any) *Error {
	return err.Append(fmt.Errorf(msg, args...))
}

func (err *Error) Len() int {
	return err.len
}

func (err *Error) Error() string {
	if err.len == 0 {
		return "<nil>"
	}
	return err.merge().Error()
}

func (err *Error) String() string {
	if err.len == 0 {
		return "<nil>"
	}
	return err.Error()
}

func (err *Error) List() []error {
	return err.list()
}

func Newf(msg string, args ...any) *Error {
	err := &Error{
		error: fmt.Errorf(msg, args...),
		len:   1,
	}
	return err
}

func New(errs ...error) *Error {
	l := []error{}
	for _, e := range errs {
		if e != nil {
			l = append(l, e)
		}
	}
	err := &Error{len: 0}
	if len(l) == 0 {
		return err
	}

	err.error = l[0]
	err.len++
	for _, e := range l[1:] {
		err.Append(e)
	}
	return err
}

func IsError(err error) bool {
	return err != nil
}

// Panic panics if `err` is not nil.
func Panic(err error, msg string, args ...any) {
	if err == nil {
		return
	}
	if e, ok := err.(*Error); ok && e.len == 0 {
		return
	}
	panic(fmt.Sprintf(msg+"\n", args...))
}

// Exit exists with the given `exitCode` if `err` is not nil.
func Exit(err error, exitCode int, msg string, args ...any) {
	if err == nil {
		return
	}
	if e, ok := err.(*Error); ok && e.len == 0 {
		return
	}
	fmt.Printf(msg+"\n", args...)
	os.Exit(exitCode)
}
