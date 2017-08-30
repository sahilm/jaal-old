package main

import (
	"github.com/sahilm/jaal/jaal"
	"github.com/sahilm/jaal/web"
	"github.com/sirupsen/logrus"
)

func main() {
	webListener := &web.Server{Address: ":8080"}

	jaal.Listen([]jaal.Listener{webListener}, logrus.New())

	select {} //block forever
}
