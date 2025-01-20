package command

import (
	"github.com/jim-nnamdi/Lkvs/pkg/runner"
	"github.com/urfave/cli/v2"
)

func StartCommand() *cli.Command {
	var (
		startRunner = &runner.StartRunner{}
	)

	cmd := &cli.Command{
		Name:  "start",
		Usage: "starts the LLKVdb server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "listen-addr",
				EnvVars:     []string{"LISTEN_ADDR"},
				Usage:       "the address that the server will listen for request on",
				Destination: &startRunner.ListenAddr,
				Value:       ":7009",
			},
		},
		Action: startRunner.Run,
	}
	return cmd
}
