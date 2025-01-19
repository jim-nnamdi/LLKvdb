package model

import "fmt"

func emptystring() string {
	return ""
}

type KeyValue struct {
	Key   int64
	Value string
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
		if kv[mid].Key == key {
			return kv[mid].Key
		}
		if kv[mid].Key > key {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}
