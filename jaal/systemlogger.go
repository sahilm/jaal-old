package jaal

import (
	"io"

	"github.com/sirupsen/logrus"
)

type SystemLogger struct {
	l *logrus.Logger
}

func NewSystemLogger(out io.Writer, indent string) *SystemLogger {
	l := logrus.New()
	l.Out = out
	l.Formatter = &logrus.JSONFormatter{}
	return &SystemLogger{l}
}

func (sl *SystemLogger) Error(err error) {
	switch err := err.(type) {
	case FatalError:
		sl.l.WithField("error", err).Fatal("")
	default:
		sl.l.WithField("error", err).Error("")
	}
}

func (sl *SystemLogger) Info(a interface{}) {
	sl.l.Info(a)
}
