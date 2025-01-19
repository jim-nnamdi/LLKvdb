package runner

import (
	"github.com/jim-nnamdi/Lkvs/pkg/server"
	"github.com/urfave/cli/v2"
)

type StartRunner struct {
	ListenAddr string
}

func (runner *StartRunner) Run(c *cli.Context) error {

	server := &server.GracefulShutdownServer{
		HTTPListenAddr: runner.ListenAddr,
		// HomeHandler:    handlers.NewHomeHandler(),
	}
	server.Start()
	return nil
}
