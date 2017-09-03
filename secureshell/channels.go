package secureshell

import (
	"github.com/sahilm/jaal/jaal"
	"golang.org/x/crypto/ssh"
)

func sshChannelHandler(newChannel ssh.NewChannel, metadata sshEventMetadata,
	eventLogHandler func(event *jaal.Event), syslogHandler func(interface{})) {

	channel, reqs, err := newChannel.Accept()
	if err != nil {
		syslogHandler(err)
		return
	}
	go sshRequestsHandler(channel, reqs, metadata, eventLogHandler, syslogHandler)
}
