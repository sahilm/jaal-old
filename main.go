package main

import (
	"os"

	"github.com/sahilm/jaal/jaal"
	"github.com/sahilm/jaal/web"
)

func main() {
	eventLogger := jaal.NewEventLogger(os.Stdout, " ")
	errLogger := jaal.NewErrLogger(os.Stderr, " ")
	webListener := &web.Server{Address: ":8080"}

	jaal.Listen([]jaal.Listener{webListener}, eventLogger, errLogger)

	select {} //block forever
}
