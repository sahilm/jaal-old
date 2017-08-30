package jaal

import (
	"encoding/json"

	"fmt"

	"github.com/sirupsen/logrus"
)

type Event struct {
	Timestamp      int64
	CorrelationID  string
	Source         string
	SourceHostName string
	Type           string
	Summary        string
	Data           interface{}
}

type Listener interface {
	Listen(eventHandler func(*Event), errHandler func(error))
}

func Listen(listeners []Listener, log *logrus.Logger) {
	for _, listener := range listeners {
		go listener.Listen(eventHandler(log), errHandler(log))
	}
}

func eventHandler(log *logrus.Logger) func(event *Event) {
	return func(e *Event) {
		b, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(b))
			//log.Info(string(b))
		}
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

type FatalError struct {
	Err error
}

func (f FatalError) Error() string {
	return f.Err.Error()
}
