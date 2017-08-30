package web

import (
	"net/http"
	"time"

	"fmt"

	"net"

	"github.com/sahilm/jaal/jaal"
)

type Server struct {
	Address      string
	eventHandler func(*jaal.Event)
	errHandler   func(error)
}

type requestData struct {
	URI    string
	Method string
	Header http.Header
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := event(r, s.errHandler)
	s.eventHandler(event)

	setHeaders(w.Header())
}

func setHeaders(header http.Header) {
	header.Add("Server", "nginx")
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("X-Powered-by", "PHP/5.4.45")
}

func event(r *http.Request, errHandler func(error)) *jaal.Event {
	id, err := jaal.ToSHA256(r.RemoteAddr)
	if err != nil {
		errHandler(err)
	}
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		errHandler(err)
	}
	now := time.Now()
	event := &jaal.Event{
		UnixTime:       now.Unix(),
		Timestamp:      now.UTC().Format(time.RFC3339),
		CorrelationID:  id[0:7],
		Source:         remoteIP,
		SourceHostName: jaal.LookupAddr(r.RemoteAddr),
		Type:           "http",
		Summary:        fmt.Sprintf("received %v at %v from %v", r.Method, r.URL, remoteIP),
		Data: &requestData{
			URI:    r.RequestURI,
			Method: r.Method,
			Header: r.Header,
		},
	}
	return event
}

func (s *Server) Listen(eventHandler func(*jaal.Event), errHandler func(error)) {
	s.eventHandler = eventHandler
	s.errHandler = errHandler

	server := &http.Server{
		Addr:           s.Address,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		errHandler(err)
	}
}
