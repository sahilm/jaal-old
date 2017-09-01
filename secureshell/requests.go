package secureshell

import (
	"github.com/sahilm/jaal/jaal"
	"golang.org/x/crypto/ssh"
)

func sshRequestsHandler(reqs <-chan *ssh.Request,
	eventLogHandler func(event *jaal.Event), sysLogHandler func(interface{})) {

	ssh.DiscardRequests(reqs)
}
