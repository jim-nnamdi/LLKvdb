package model

import "fmt"

func emptystring() string {
	return ""
}

type KeyValue struct {
	key   int64
	value string
}

func TableDiskError(err error) {
	if err != nil {
		fmt.Printf("err:'%s'\n", err)
	}
}
