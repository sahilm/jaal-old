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
	l.Formatter = &LogFormatter{indent}
	l.Out = out
	return &ErrLogger{l, indent}
}

func (el *ErrLogger) Log(err error) {
	switch e := err.(type) {
	case FatalError:
		el.l.Fatal(e)
	default:
		el.l.Error(e)
	}
}
