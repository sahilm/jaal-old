package secureshell

import (
	"fmt"

	"github.com/sahilm/jaal/jaal"
)

type Server struct {
	Address string
}

func (s *Server) Listen(eventHandler func(*jaal.Event), systemLogHandler func(interface{})) {
	go systemLogHandler(fmt.Sprintf("starting ssh listener at %v", s.Address))
}
