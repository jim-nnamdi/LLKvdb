package model

import "os"

type WAL struct {
	file *os.File
}
