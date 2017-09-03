package secureshell

import (
	"fmt"

	"io"

	"github.com/sahilm/jaal/jaal"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type exec struct {
	Command string
}

type tcpipForward struct {
	BindAddress string
	BindPort    uint32
}

func sshRequestsHandler(channel ssh.Channel, reqs <-chan *ssh.Request, metadata sshEventMetadata,
	eventLogHandler func(event *jaal.Event), sysLogHandler func(interface{})) {

	for r := range reqs {
		switch r.Type {
		case "exec":
			data := exec{}
			err := ssh.Unmarshal(r.Payload, &data)
			if err != nil {
				sysLogHandler(err)
			}
			eventLogHandler(requestEvent(metadata, r.Type, data))
			channel.Close()
		case "tcpip-forward":
			data := tcpipForward{}
			err := ssh.Unmarshal(r.Payload, &data)
			if err != nil {
				sysLogHandler(err)
			}
			eventLogHandler(requestEvent(metadata, r.Type, data))
		case "shell":
			term := terminal.NewTerminal(channel, "$ ")
			for {
				line, err := term.ReadLine()
				if err != nil {
					if err != io.EOF {
						sysLogHandler(err)
					}
					channel.Close()
					break
				}
				eventLogHandler(termLineEvent(metadata, line))
			}
		}

		if r.WantReply {
			err := r.Reply(true, nil)
			if err != nil {
				sysLogHandler(err)
			}
		}
	}
}

func requestEvent(metadata sshEventMetadata, reqType string, data interface{}) *jaal.Event {
	event := &jaal.Event{
		Type:          fmt.Sprintf("ssh %v", reqType),
		Source:        metadata.RemoteIP,
		CorrelationID: metadata.CorrelationID,
	}
	jaal.AddEventMetadata(event)
	event.Summary = fmt.Sprintf("ssh request: %v from %v(%v)", reqType, event.SourceHostName, event.Source)
	event.Data = data
	return event
}

func termLineEvent(metadata sshEventMetadata, line string) *jaal.Event {
	event := &jaal.Event{
		Type:          "ssh command",
		Source:        metadata.RemoteIP,
		CorrelationID: metadata.CorrelationID,
	}
	jaal.AddEventMetadata(event)
	event.Summary = fmt.Sprintf("ssh command: %v from %v(%v)",
		line, event.SourceHostName, event.Source)
	event.Data = line
	return event
}
