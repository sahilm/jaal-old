package main

import (
	"os"

	"github.com/sahilm/jaal/jaal"
	"github.com/sahilm/jaal/web"
)

func main() {
	systemLogger := jaal.NewSystemLogger(os.Stderr, " ")
	eventLogger := jaal.NewEventLogger(os.Stdout, systemLogger, " ")
	webListener := &web.Server{Address: ":9000"}

	jaal.Listen([]jaal.Listener{webListener}, eventLogger, systemLogger)

	select {} //block forever
}
