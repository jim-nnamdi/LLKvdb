package model

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrOpeningFile  = "could not open write ahead log file"
	ErrFileProblems = "could not create or append data to files"
	ErrGenericWAL   = "write ahead log encountered a read or write error"
)

type WAL struct {
	file *os.File
}

func NewWAL(walLoc string) (*WAL, error) {
	file, err := os.OpenFile(walLoc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	WriteAheadLogError(err)
	return &WAL{file: file}, nil
}

func (wal *WAL) Write(key int64, value string) error {
	_, err := wal.file.WriteString(fmt.Sprintf("%d:%s\n", key, value))
	GenericWriteAheadLogError(err)
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
		data = append(data, KeyValue{Key: key, Value: value})
	}
	if err = scanner.Err(); err != nil {
		return nil, errors.New(err.Error())
	}
	return data, nil
}

func (wal *WAL) Delete(key int64) (bool, error) {
	fsk, err := os.Open(wal.file.Name())
	if err != nil {
		return false, err
	}
	defer fsk.Close()
	var lines []string
	scanner := bufio.NewScanner(fsk)
	for scanner.Scan() {
		line := scanner.Text()
		skey := strconv.FormatInt(key, 10)
		if strings.HasPrefix(line, skey+":") {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	ofile, err := os.OpenFile(wal.file.Name(), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false, err
	}
	defer ofile.Close()
	writer := bufio.NewWriter(ofile)
	for _, line := range lines {
		fmt.Println(line)
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return false, err
		}
	}
	fmt.Println(lines)
	writer.Flush()
	fmt.Printf("key '%d' removed successfully\n", key)
	return true, nil
}

func (wal *WAL) WALDump(key int64) (string, error) {
	walf, err := os.Open(wal.file.Name())
	if err != nil {
		fmt.Println(err)
		return emptystring(), err
	}
	defer walf.Close()
	scanner := bufio.NewScanner(walf)
	for scanner.Scan() {
		line := scanner.Text()
		skey := strconv.FormatInt(key, 10)
		if strings.HasPrefix(line, skey+":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return parts[1], nil
			}
		}
	}
	return emptystring(), errors.New("cannot find value")
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
