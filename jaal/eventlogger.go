package jaal

import (
	"io"

	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type EventLogger interface {
	Log(event *Event)
}

type EventLog struct {
	l  *logrus.Logger
	sl *SystemLogger
}

func NewEventLogger(out io.Writer, systemLogger *SystemLogger, indent string) *EventLog {
	l := logrus.New()
	l.Out = out
	l.Formatter = &eventLogFormatter{indent}
	l.Hooks.Add(&SlackNotifier{})
	return &EventLog{l, systemLogger}
}

func (el *EventLog) Log(event *Event) {
	el.l.WithField("data", event).Info("")
}

type eventLogFormatter struct {
	indent string
}

func (f *eventLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var serialized []byte
	var err error
	serialized, err = json.MarshalIndent(entry.Data["data"], "", f.indent)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
