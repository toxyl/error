package errors

import (
	"fmt"
	"os"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
	}{
		{"no error", New(nil)},
		{"single error", Newf("hello")},
		{"double error", Newf("hello").Appendf("world")},
		{"triple error", Newf("hello").Appendf("world").Append(fmt.Errorf("I'm an error"))},
		{"mixed errors", New(Newf("hello"), nil, fmt.Errorf("world"), nil, nil, fmt.Errorf("I'm an error")).Append(os.ErrDeadlineExceeded)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(IsError(tt.err.error))
			t.Log(tt.err)
			t.Logf("%v", tt.err.List())
		})
	}
}

func TestChain(t *testing.T) {
	tests := []struct {
		name string
		time int
	}{}
	for i := 0; i < 10; i++ {
		tests = append(tests, struct {
			name string
			time int
		}{
			name: "Test " + fmt.Sprint(i),
			time: i,
		})
	}

	e := New()
	Exit(e, 1, "there's no error, so we shouldn't exit")
	Exit(nil, 2, "shouldn't exit either")
	Panic(e, "we shouldn't panic here")
	Warn(e, "nothing to warn about")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.time%7 == 0 {
				e.Appendf("bar")
			} else if tt.time%5 == 0 {
				e.Appendf("foo")
			} else if tt.time%2 == 0 {
				e.Appendf("even")
			} else {
				e.Appendf("%d", tt.time)
			}
		})
	}
	t.Log(IsError(e.error))
	t.Log(e)
	t.Logf("%v", e.List())
	Warn(e, "this is a warning! you better watch out!")
	Panic(e, "we've reached the end of the test")
}
