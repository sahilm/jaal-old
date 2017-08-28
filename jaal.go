package jaal

import (
	"github.com/sirupsen/logrus"
)

type Event struct {
	UnixTime      int64
	CorrelationID string
	Source        string
	Data          interface{}
}

type Listener interface {
	Listen(eventHandler func(*Event), errHandler func(error))
}

type FatalError struct {
	Err error
}

func (f FatalError) Error() string {
	return f.Err.Error()
}

func Listen(listeners []Listener, log *logrus.Logger) {
	for _, listener := range listeners {
		go listener.Listen(eventHandler(log), errHandler(log))
	}
}

func eventHandler(log *logrus.Logger) func(event *Event) {
	return func(e *Event) {
		log.Info(e)
	}
}

func errHandler(log *logrus.Logger) func(error) {
	return func(err error) {
		switch e := err.(type) {
		case FatalError:
			log.Fatal(e)
		default:
			log.Info(e)
		}
	}
}
