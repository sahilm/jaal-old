package secureshell

import (
	"fmt"

	"crypto/rand"
	"crypto/rsa"

	"net"
	"time"

	"io/ioutil"

	"github.com/sahilm/jaal/jaal"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	Address        string
	sshHostKeyFile string
	ioTimeout      time.Duration
	quit           chan bool
}

func NewServer(address string, sshHostKeyFile string) *Server {
	quit := make(chan bool)
	timeout := 5 * time.Minute

	return &Server{
		Address:        address,
		ioTimeout:      timeout,
		sshHostKeyFile: sshHostKeyFile,
		quit:           quit,
	}
}

func (s *Server) Stop() {
	s.quit <- true
}

func (s *Server) Listen(eventHandler func(*jaal.Event), systemLogHandler func(interface{})) {
	go systemLogHandler(fmt.Sprintf("starting ssh listener at %v", s.Address))

	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		systemLogHandler(jaal.FatalError{Err: err})
	}

	defer listener.Close()

loop:
	for {
		select {
		case <-s.quit:
			break loop
		default:
			s.accept(listener, eventHandler, systemLogHandler)
		}
	}
}

func (s *Server) accept(listener net.Listener, eventHandler func(*jaal.Event), systemLogHandler func(interface{})) {
	tcpConn, err := listener.Accept()
	if err != nil {
		systemLogHandler(err)
		return
	}
	tcpConn.SetDeadline(time.Now().Add(s.ioTimeout))

	config := config(s.sshHostKeyFile, systemLogHandler)
	sshServerConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
	if err != nil {
		systemLogHandler(err)
		return
	}

	defer sshServerConn.Close()

	sha, err := jaal.ToSHA256(sshServerConn.SessionID())
	if err != nil {
		systemLogHandler(err)
		sshServerConn.Close()
		return
	}

	remoteIP, _, err := net.SplitHostPort(sshServerConn.RemoteAddr().String())
	if err != nil {
		remoteIP = sshServerConn.RemoteAddr().String()
	}

	metadata := sshEventMetadata{
		RemoteIP:      remoteIP,
		ClientVersion: string(sshServerConn.ClientVersion()),
		CorrelationID: sha[0:7],
		Password:      sshServerConn.Permissions.Extensions[sha[0:7]], // See PasswordCallback
		Username:      sshServerConn.User(),
	}
	eventHandler(loginEvent(metadata))

	go sshRequestsHandler(reqs, eventHandler, systemLogHandler)

	for newChannel := range chans {
		sshChannelHandler(newChannel, metadata, eventHandler, systemLogHandler)
	}
}

func config(sshHostKeyFile string, systemLogHandler func(interface{})) *ssh.ServerConfig {
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

	hostKey, err := hostKey(sshHostKeyFile)
	if err != nil {
		systemLogHandler(jaal.FatalError{Err: err})
	}

	c.AddHostKey(hostKey)

	return c
}

func hostKey(sshHostKeyFile string) (ssh.Signer, error) {
	if sshHostKeyFile != "" {
		keyBytes, err := ioutil.ReadFile(sshHostKeyFile)
		if err != nil {
			return nil, jaal.FatalError{Err: err}
		}

		key, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			return nil, jaal.FatalError{Err: err}
		}
		return key, nil
	} else {
		return generateHostKey()
	}
}

func generateHostKey() (ssh.Signer, error) {
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

func loginEvent(metadata sshEventMetadata) *jaal.Event {
	event := &jaal.Event{
		Type:          "ssh login",
		Source:        metadata.RemoteIP,
		CorrelationID: metadata.CorrelationID,
	}
	now := time.Now()
	event.SourceHostName = jaal.LookupAddr(event.Source)
	event.UnixTime = now.Unix()
	event.Timestamp = now.UTC().Format(time.RFC3339)
	event.Summary = fmt.Sprintf("ssh login with Username: %v password: %v from %v(%v)",
		metadata.Username, metadata.Password, event.SourceHostName, event.Source)
	event.Data = metadata
	return event
}

type sshEventMetadata struct {
	CorrelationID string
	RemoteIP      string
	ClientVersion string
	Username      string
	Password      string
}
