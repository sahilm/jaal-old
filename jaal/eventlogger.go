package jaal

import (
	"io"

	"fmt"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type EventLogger interface {
	Log(event *Event)
}

type EventLog struct {
	l      *logrus.Logger
	indent string
}

func NewEventLogger(out io.Writer, indent string) *EventLog {
	l := logrus.New()
	l.Formatter = &LogFormatter{indent}
	l.Out = out
	return &EventLog{l, indent}
}

func (el *EventLog) Log(event *Event) {
	enrichEvent(event)

	el.l.WithField("data", event).Info("")
}

func enrichEvent(event *Event) {
	now := time.Now()
	event.SourceHostName = lookupAddr(event.Source)
	event.UnixTime = now.Unix()
	event.Timestamp = now.UTC().Format(time.RFC3339)
}

func lookupAddr(address string) string {
	ip, _, err := net.SplitHostPort(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to lookup %v. Error: %v", address, err)
		return "" // Don't care on err, just return nothing
	}
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to lookup %v. Error: %v", address, err)
		return "" // Don't care on err, just return nothing
	}
	return hosts[0]
}
