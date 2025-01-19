package model

import (
	"fmt"
	"os"
)

var (
	ErrFileProblems = "could not create or append data to files"
	ErrGenericWAL   = "write ahead log encountered a read or write error"
)

type WAL struct {
	file *os.File
}

func NewWAL(walLoc string) *WAL {
	file, err := os.OpenFile(walLoc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	WriteAheadLogError(err)
	defer file.Close()
	return &WAL{file: file}
}

func (wal *WAL) Write(key int64, value string) error {
	memsize, err := wal.file.WriteString(fmt.Sprintf("%d:%s\n", key, value))
	GenericWriteAheadLogError(err)
	fmt.Printf("bytes written to WALog: '%d'\n", memsize)
	return err
}

func WriteAheadLogError(err error) {
	if err != nil {
		fmt.Printf("WAL:'%s'\n", err)
		fmt.Println(ErrFileProblems)
	}
}

func GenericWriteAheadLogError(err error) {
	if err != nil {
		fmt.Printf("WAL:'%s'\n", err)
		fmt.Println(ErrGenericWAL)
	}
}
