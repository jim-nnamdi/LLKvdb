package model

import (
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

func (sstable *SSTable) Write(data map[int64]string) error {
	fsk, err := os.Create(sstable.diskfile)
	TableDiskError(err)
	defer fsk.Close()
	store := json.NewEncoder(fsk)
	return store.Encode(data)
}

func (sstable *SSTable) Read() map[int64]string {
	fsk, err := os.Open(sstable.diskfile)
	TableDiskError(err)
	var data map[int64]string
	decode := json.NewDecoder(fsk)
	decode.Decode(&data)
	return data
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

func TableDiskError(err error) {
	if err != nil {
		fmt.Printf("err:'%s'\n", err)
	}
}
