package main

import (
	"os"

	"github.com/sahilm/jaal/jaal"
	"github.com/sahilm/jaal/web"
)

func main() {
	errLogger := jaal.NewErrLogger(os.Stderr, " ")
	eventLogger := jaal.NewEventLogger(os.Stdout, errLogger, " ")
	webListener := &web.Server{Address: ":9000"}

	jaal.Listen([]jaal.Listener{webListener}, eventLogger, errLogger)

	select {} //block forever
}
