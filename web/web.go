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
	id, err := jaal.ToSHA256(r.RemoteAddr)
	if err != nil {
		s.errHandler(err)
	}

	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		s.errHandler(err)
	}

	event := &jaal.Event{
		Timestamp:      time.Now().Unix(),
		CorrelationID:  id,
		Source:         remoteIP,
		SourceHostName: jaal.LookupAddr(r.RemoteAddr),
		Type:           "http",
		Summary:        eventSummary(r, remoteIP),
		Data:           processRequest(r),
	}
	s.eventHandler(event)
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

func eventSummary(r *http.Request, remoteIP string) string {
	return fmt.Sprintf("received %v at %v from %v", r.Method, r.URL, remoteIP)
}

func processRequest(r *http.Request) *requestData {
	return &requestData{
		URI:    r.RequestURI,
		Method: r.Method,
		Header: r.Header,
	}
}
