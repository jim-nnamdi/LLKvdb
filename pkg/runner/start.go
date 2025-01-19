package runner

import (
	"github.com/jim-nnamdi/Lkvs/pkg/handlers"
	"github.com/jim-nnamdi/Lkvs/pkg/model"
	"github.com/jim-nnamdi/Lkvs/pkg/server"
	"github.com/urfave/cli/v2"
)

type StartRunner struct {
	ListenAddr string
}

func (runner *StartRunner) Run(c *cli.Context) error {
	var (
		Fsys = model.NewFilesys("wal.txt", 1024)
	)
	server := &server.GracefulShutdownServer{
		HTTPListenAddr:      runner.ListenAddr,
		PutHandler:          handlers.NewPutHandler(Fsys),
		ReadHandler:         handlers.NewReadHandler(Fsys),
		ReadKeyRangeHandler: handlers.NewReadKeyRangeHandler(Fsys),
		BatchPutHandler:     handlers.NewBatchPutHandler(Fsys),
		DeleteHandler:       handlers.NewDeleteHandler(Fsys),
	}
	server.Start()
	return nil
}
