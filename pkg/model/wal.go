package model

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

var (
	ErrOpeningFile  = "could not open write ahead log file"
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

func (wal *WAL) Replay() ([]KeyValue, error) {
	fsk, err := os.Open(wal.file.Name())
	if err != nil {
		fmt.Println(ErrOpeningFile)
		return nil, errors.New(ErrOpeningFile)
	}
	defer fsk.Close()
	var data []KeyValue
	scanner := bufio.NewScanner(fsk)
	for scanner.Scan() {
		lines := scanner.Text()
		var key int64
		var value string
		fmt.Sscanf(lines, "%d:%s", key, value)
		data = append(data, KeyValue{key: key, value: value})
	}
	if err = scanner.Err(); err != nil {
		return nil, errors.New(err.Error())
	}
	return data, nil
}

func (wal *WAL) Walclose() error {
	return wal.file.Close()
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
