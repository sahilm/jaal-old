package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/sahilm/jaal/jaal"
	"github.com/sahilm/jaal/secureshell"
	"github.com/sahilm/jaal/web"
)

var version = "latest"

func main() {
	var opts struct {
		SSHHostKeyFile string `long:"ssh-host-key-file" description:"path to the ssh host key file"`
		SSHPort        uint   `long:"ssh-port" description:"port to listen on for ssh traffic" default:"22"`
		HTTPPort       uint   `long:"http-port" description:"port to listen on for http traffic" default:"80"`
		Version        func() `long:"version" description:"print version and exit"`
	}

	opts.Version = func() {
		fmt.Fprintf(os.Stderr, "%v\n", version)
		os.Exit(0)
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	systemLogger := jaal.NewSystemLogger(os.Stderr)
	eventLogger := jaal.NewEventLogger(os.Stdout, systemLogger, " ")
	webListener := web.NewServer(fmt.Sprintf(":%v", opts.HTTPPort))
	sshListener := secureshell.NewServer(fmt.Sprintf(":%v", opts.SSHPort), opts.SSHHostKeyFile)

	jaal.Listen([]jaal.Listener{webListener, sshListener}, eventLogger, systemLogger)
}
