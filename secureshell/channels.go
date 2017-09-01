package secureshell

import (
	"fmt"
	"time"

	"io"

	"github.com/sahilm/jaal/jaal"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type tcpip struct {
	DestinationAddress string
	SourceAddress      string
	DestinationPort    uint32
	SourcePort         uint32
}

type x11 struct {
	SourceAddress string
	SourcePort    uint32
}

func sshChannelHandler(newChannel ssh.NewChannel, metadata sshEventMetadata,
	eventLogHandler func(event *jaal.Event), syslogHandler func(interface{})) {

	var out struct{}

	switch newChannel.ChannelType() {
	case "direct-tcpip":
		out := tcpip{}
		err := ssh.Unmarshal(newChannel.ExtraData(), &out)
		if err != nil {
			syslogHandler(err)
		}
	case "x11":
		out := x11{}
		err := ssh.Unmarshal(newChannel.ExtraData(), &out)
		if err != nil {
			syslogHandler(err)
		}
	}

	eventLogHandler(channelEvent(metadata, newChannel.ChannelType(), out))

	channel, reqs, err := newChannel.Accept()
	if err != nil {
		syslogHandler(err)
		return
	}

	defer channel.Close()
	go sshRequestsHandler(reqs, metadata, eventLogHandler, syslogHandler)

	if newChannel.ChannelType() == "session" {
		term := terminal.NewTerminal(channel, "$ ")
		for {
			line, err := term.ReadLine()
			if err != nil {
				if err != io.EOF {
					syslogHandler(err)
				}
				break
			}
			eventLogHandler(termLineEvent(metadata, line))
		}
	} else {
		data := make([]byte, 1024)
		n, err := channel.Read(data)
		if err != nil && err != io.EOF {
			syslogHandler(err)
		}
		eventLogHandler(genericEvent(metadata, data, n))
	}
}

func channelEvent(metadata sshEventMetadata, channelType string, data interface{}) *jaal.Event {
	event := &jaal.Event{
		Type:          "new ssh channel",
		Source:        metadata.RemoteIP,
		CorrelationID: metadata.CorrelationID,
	}
	enrichEvent(event)
	event.Summary = fmt.Sprintf("ssh channel open: %v from %v(%v)",
		channelType, event.SourceHostName, event.Source)
	event.Data = data
	return event
}

func termLineEvent(metadata sshEventMetadata, line string) *jaal.Event {
	event := &jaal.Event{
		Type:          "ssh command",
		Source:        metadata.RemoteIP,
		CorrelationID: metadata.CorrelationID,
	}
	enrichEvent(event)
	event.Summary = fmt.Sprintf("ssh command: %v from %v(%v)",
		line, event.SourceHostName, event.Source)
	event.Data = line
	return event
}

func genericEvent(metadata sshEventMetadata, bytes []byte, n int) *jaal.Event {
	event := &jaal.Event{
		Type:          "ssh",
		Source:        metadata.RemoteIP,
		CorrelationID: metadata.CorrelationID,
	}
	enrichEvent(event)
	event.Summary = fmt.Sprintf("ssh data recv from %v(%v)", event.SourceHostName, event.Source)
	event.Data = string(bytes[:n])
	return event
}

func enrichEvent(event *jaal.Event) {
	now := time.Now()
	event.SourceHostName = jaal.LookupAddr(event.Source)
	event.UnixTime = now.Unix()
	event.Timestamp = now.UTC().Format(time.RFC3339)
}
