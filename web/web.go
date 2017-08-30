package web

import (
	"net/http"
	"time"

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
	event := &jaal.Event{
		UnixTime:      time.Now().Unix(),
		CorrelationID: id,
		Source:        "http",
		Data:          processRequest(r),
	}
	s.eventHandler(event)
}

func processRequest(r *http.Request) *requestData {
	return &requestData{
		URI:    r.RequestURI,
		Method: r.Method,
		Header: r.Header,
	}
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
