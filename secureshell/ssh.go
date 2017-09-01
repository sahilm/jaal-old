package secureshell

import (
	"fmt"

	"crypto/rand"
	"crypto/rsa"

	"net"
	"time"

	"github.com/sahilm/jaal/jaal"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Server struct {
	Address   string
	ioTimeout time.Duration
	quit      chan bool
}

func NewServer(address string) *Server {
	quit := make(chan bool)
	timeout := 5 * time.Minute

	return &Server{address, timeout, quit}
}

func (s *Server) Stop() {
	s.quit <- true
}

func (s *Server) Listen(eventHandler func(*jaal.Event), systemLogHandler func(interface{})) {
	go systemLogHandler(fmt.Sprintf("starting ssh listener at %v", s.Address))

	config := config(systemLogHandler)

	listener, err := net.Listen("tcp", s.Address)

	if err != nil {
		systemLogHandler(jaal.FatalError{Err: err})
	}

	defer listener.Close()

	go func() {
		for {
			tcpConn, err := listener.Accept()
			if err != nil {
				systemLogHandler(err)
				continue
			}

			tcpConn.SetDeadline(time.Now().Add(s.ioTimeout))

			sshServerConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
			if err != nil {
				systemLogHandler(err)
				continue
			}

			sha, err := jaal.ToSHA256(sshServerConn.SessionID())
			if err != nil {
				systemLogHandler(err)
				sshServerConn.Close()
				continue
			}

			metadata := sshEventMetadata{
				RemoteAddr:    sshServerConn.RemoteAddr(),
				ClientVersion: string(sshServerConn.ClientVersion()),
				CorrelationID: sha[0:7],
				Password:      sshServerConn.Permissions.Extensions[sha[0:7]], // See PasswordCallback
				Username:      sshServerConn.User(),
			}

			go ssh.DiscardRequests(reqs)

			go eventHandler(loginEvent(metadata))

			// Service the incoming Channel channel.
			for newChannel := range chans {
				// Channels have a type, depending on the application level
				// protocol intended. In the case of a shell, the type is
				// "session" and ServerShell may be used to present a simple
				// terminal interface.
				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					continue
				}
				channel, requests, err := newChannel.Accept()
				if err != nil {
					systemLogHandler(err)
				}

				// Sessions have out-of-band requests such as "shell",
				// "pty-req" and "env".  Here we handle only the
				// "shell" request.
				go func(in <-chan *ssh.Request) {
					for req := range in {
						req.Reply(req.Type == "shell", nil)
					}
				}(requests)

				term := terminal.NewTerminal(channel, "> ")

				go func() {
					defer channel.Close()
					for {
						line, err := term.ReadLine()
						if err != nil {
							break
						}
						fmt.Println(line)
					}
				}()
			}
		}
	}()

	<-s.quit
}

func config(systemLogHandler func(interface{})) *ssh.ServerConfig {
	c := &ssh.ServerConfig{
		// Allow everyone to login. This is a honeypot 😀
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			m := make(map[string]string)
			sessionID, err := jaal.ToSHA256(c.SessionID())
			if err != nil {
				systemLogHandler(err)
			}
			m[sessionID[0:7]] = string(pass)
			perms := &ssh.Permissions{
				Extensions: m, // Use extensions as a mechanism to save the session ID -> password mapping
			}
			return perms, nil
		},
	}

	hostKey, err := hostKey()
	if err != nil {
		systemLogHandler(jaal.FatalError{Err: err})
	}

	c.AddHostKey(hostKey)

	return c
}

func hostKey() (ssh.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func loginEvent(params sshEventMetadata) *jaal.Event {
	host, _, err := net.SplitHostPort(params.RemoteAddr.String())
	if err != nil {
		host = params.RemoteAddr.String()
	}

	event := &jaal.Event{
		Type:          "ssh login",
		Source:        host,
		CorrelationID: params.CorrelationID,
	}
	now := time.Now()
	event.SourceHostName = jaal.LookupAddr(event.Source)
	event.UnixTime = now.Unix()
	event.Timestamp = now.UTC().Format(time.RFC3339)
	event.Summary = fmt.Sprintf("ssh login with Username: %v password: %v from %v(%v)",
		params.Username, params.Password, event.SourceHostName, event.Source)
	event.Data = params
	return event
}

type sshEventMetadata struct {
	CorrelationID string
	RemoteAddr    net.Addr
	ClientVersion string
	Username      string
	Password      string
}
