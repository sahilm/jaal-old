package jaal

import (
	"io"

	"net"
	"time"

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
	return &EventLog{l, systemLogger}
}

func (el *EventLog) Log(event *Event) {
	enrichEvent(event, el.sl)
	el.l.WithField("data", event).Info("")
}

func enrichEvent(event *Event, sl *SystemLogger) {
	now := time.Now()
	event.SourceHostName = lookupAddr(event.Source, sl)
	event.UnixTime = now.Unix()
	event.Timestamp = now.UTC().Format(time.RFC3339)
}

func lookupAddr(address string, sl *SystemLogger) string {
	hosts, err := net.LookupAddr(address)
	if err != nil {
		sl.Error(err)
		return "" // Don't care on err, just return nothing
	}
	return hosts[0]
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
	return serialized, nil
}
