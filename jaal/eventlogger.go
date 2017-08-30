package jaal

import (
	"io"

	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type EventLogger interface {
	Log(event *Event)
}

type EventLog struct {
	l      *logrus.Logger
	el     *ErrLogger
	indent string
}

func NewEventLogger(out io.Writer, errLogger *ErrLogger, indent string) *EventLog {
	l := logrus.New()
	l.Formatter = &LogFormatter{indent}
	l.Out = out
	return &EventLog{l, errLogger, indent}
}

func (eventLog *EventLog) Log(event *Event) {
	enrichEvent(event, eventLog.el)

	eventLog.l.WithField("data", event).Info("")
}

func enrichEvent(event *Event, el *ErrLogger) {
	now := time.Now()
	event.SourceHostName = lookupAddr(event.Source, el)
	event.UnixTime = now.Unix()
	event.Timestamp = now.UTC().Format(time.RFC3339)
}

func lookupAddr(address string, el *ErrLogger) string {
	hosts, err := net.LookupAddr(address)
	if err != nil {
		el.Log(err)
		return "" // Don't care on err, just return nothing
	}
	return hosts[0]
}
