package jaal

import (
	"io"

	"github.com/sirupsen/logrus"
)

type ErrLogger struct {
	l      *logrus.Logger
	indent string
}

func NewErrLogger(out io.Writer, indent string) *ErrLogger {
	l := logrus.New()
	l.Out = out
	return &ErrLogger{l, indent}
}

func (el *ErrLogger) Log(err error) {
	switch err := err.(type) {
	case FatalError:
		el.l.Fatal(err)
	default:
		el.l.Error(err)
	}
}
