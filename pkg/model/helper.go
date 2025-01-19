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

func BinarySearch(kv []KeyValue, key int64) int64 {
	low := 0
	high := len(kv) - 1
	if low == high {
		return key
	}

	for low < high {
		mid := (low + high) / 2
		if kv[mid].key == key {
			return kv[mid].key
		}
		if kv[mid].key > key {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}
