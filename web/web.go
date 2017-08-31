package web

import (
	"net/http"
	"time"

	"fmt"

	"net"

	"github.com/sahilm/jaal/jaal"
)

type Server struct {
	Address          string
	eventHandler     func(*jaal.Event)
	systemLogHandler func(interface{})
}

type requestData struct {
	URI    string
	Method string
	Header http.Header
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := event(r, s.systemLogHandler)
	go s.eventHandler(event)

	setHeaders(w.Header())
}

func setHeaders(header http.Header) {
	header.Add("Server", "nginx")
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("X-Powered-by", "PHP/5.4.45")
}

func event(r *http.Request, sysLogHandler func(interface{})) *jaal.Event {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		sysLogHandler(err)
	}

	id, err := jaal.ToSHA256(remoteIP)

	if err != nil {
		sysLogHandler(err)
	}

	event := &jaal.Event{
		CorrelationID: id[0:7],
		Source:        remoteIP,
		Type:          "http",
		Summary:       fmt.Sprintf("received %v at %v from %v", r.Method, r.URL, remoteIP),
		Data: &requestData{
			URI:    r.RequestURI,
			Method: r.Method,
			Header: r.Header,
		},
	}
	return event
}

func (s *Server) Listen(eventHandler func(*jaal.Event), systemLogHandler func(interface{})) {
	go systemLogHandler(fmt.Sprintf("starting web listener at %v", s.Address))

	s.eventHandler = eventHandler
	s.systemLogHandler = systemLogHandler

	server := &http.Server{
		Addr:           s.Address,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		systemLogHandler(jaal.FatalError{Err: err})
	}
}
