package model

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	ErrInexistentKey = "key does not exists"
)

type SSTable struct {
	diskfile string
}

func Newsstable(diskfile string) *SSTable {
	return &SSTable{diskfile: diskfile}
}

func (sstable *SSTable) Write(data []KeyValue) error {
	fsk, err := os.Create(sstable.diskfile)
	TableDiskError(err)
	defer fsk.Close()
	for _, kv := range data {
		_, err := fsk.WriteString(fmt.Sprintf("%d:%s\n", kv.key, kv.value))
		if err != nil {
			return err
		}
	}
	return nil
}

func (sstable *SSTable) Load() ([]KeyValue, error) {
	fsk, err := os.Open(sstable.diskfile)
	TableDiskError(err)
	defer fsk.Close()
	var data []KeyValue
	scanner := bufio.NewScanner(fsk)
	for scanner.Scan() {
		line := scanner.Text()
		var key int64
		var value string
		fmt.Sscanf(line, "%d:%s", key, value)
		data = append(data, KeyValue{key: key, value: value})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func (sstable *SSTable) Get(key int64) (string, error) {
	data, err := sstable.Load()
	if err != nil {
		fmt.Println(ErrInexistentKey)
		return emptystring(), err
	}
	Left := 0
	Right := len(data) - 1
	if Left == Right {
		result := data[key]
		return result.value, nil
	}
	for Left < Right {
		sMid := (Left + Right) / 2
		if data[sMid].key == key {
			result := data[data[sMid].key]
			return result.value, nil
		}
		if data[sMid].key < key {
			Left = sMid + 1
		}
		if data[sMid].key > key {
			Right = sMid - 1
		}
	}
	return emptystring(), nil
}

func (sstable *SSTable) ReadOne(key int64) (string, error) {
	fsk, err := os.Open(sstable.diskfile)
	TableDiskError(err)
	defer fsk.Close()

	var data map[int64]string
	decode := json.NewDecoder(fsk)
	results := decode.Decode(&data)
	for datakey, value := range data {
		if datakey == key {
			return value, results
		}
	}
	fmt.Println(ErrInexistentKey)
	return ErrInexistentKey, errors.New(ErrInexistentKey)
}
