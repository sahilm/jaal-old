package jaal

import (
	"io"

	"github.com/sirupsen/logrus"
)

type SystemLogger struct {
	l *logrus.Logger
}

func NewSystemLogger(out io.Writer) *SystemLogger {
	l := logrus.New()
	l.Out = out
	return &SystemLogger{l}
}

func (sl *SystemLogger) Error(err error) {
	switch err := err.(type) {
	case FatalError:
		sl.l.Fatal(err)
	default:
		sl.l.Error(err)
	}
}

func (sl *SystemLogger) Info(a interface{}) {
	sl.l.Info(a)
}
