package model

import (
	"sort"
	"sync"
)

type Memtable struct {
	data  map[int64]string
	mutex sync.RWMutex
}

type KeyValue struct {
	key   int64
	value string
}

func Newmemtable() *Memtable {
	return &Memtable{data: make(map[int64]string)}
}

func (mem *Memtable) Write(key int64, value string) {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	mem.data[key] = value
}

func (mem *Memtable) Read(key int64) (string, bool) {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	value, exists := mem.data[key]
	return value, exists
}

func (mem *Memtable) Delete(key int64) {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	delete(mem.data, key)
}

func (mem *Memtable) Flush() []KeyValue {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	keys := make([]int64, 0, len(mem.data))
	for key := range mem.data {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })
	sortedData := make([]KeyValue, 0, len(mem.data))
	for _, key := range keys {
		sortedData = append(sortedData, KeyValue{key: key, value: mem.data[key]})
	}
	mem.data = make(map[int64]string)
	return sortedData
}
