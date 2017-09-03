package secureshell

import (
	"fmt"

	"github.com/sahilm/jaal/jaal"
	"golang.org/x/crypto/ssh"
)

type env struct {
	Name, Value string
}

type exec struct {
	Command string
}

type tcpipForward struct {
	BindAddress string
	BindPort    uint32
}

func sshRequestsHandler(reqs <-chan *ssh.Request, metadata sshEventMetadata,
	eventLogHandler func(event *jaal.Event), sysLogHandler func(interface{})) {

	for r := range reqs {
		switch r.Type {
		case "env":
			data := env{}
			err := ssh.Unmarshal(r.Payload, &data)
			if err != nil {
				sysLogHandler(err)
			}
			eventLogHandler(requestEvent(metadata, r.Type, data))
		case "exec":
			data := exec{}
			err := ssh.Unmarshal(r.Payload, &data)
			if err != nil {
				sysLogHandler(err)
			}
			eventLogHandler(requestEvent(metadata, r.Type, data))
		case "tcpip-forward":
			data := tcpipForward{}
			err := ssh.Unmarshal(r.Payload, &data)
			if err != nil {
				sysLogHandler(err)
			}
			eventLogHandler(requestEvent(metadata, r.Type, data))
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
