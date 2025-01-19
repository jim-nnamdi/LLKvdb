package model

import (
	"sort"
	"sync"
)

type Memtable struct {
	data  map[int64]string
	mutex sync.RWMutex
}

func Newmemtable() *Memtable {
	return &Memtable{data: make(map[int64]string)}
}

func (mem *Memtable) Put(key int64, value string) {
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
		sortedData = append(sortedData, KeyValue{Key: key, Value: mem.data[key]})
	}
	mem.data = make(map[int64]string)
	return sortedData
}

func (mem *Memtable) FlushTableBenchMarkTest() []KeyValue {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	keys := make([]int64, 0, len(mem.data))
	for key := range mem.data {
		keys = append(keys, key)
	}
	QuicksortAlgorithm(keys)
	// sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	sorted := make([]KeyValue, 0, len(mem.data))
	for _, key := range keys {
		sorted = append(sorted, KeyValue{Key: key, Value: mem.data[key]})
	}
	mem.data = make(map[int64]string)
	return sorted
}

func (mem *Memtable) Dump() ([]KeyValue, error) {
	mem.mutex.RLock()
	defer mem.mutex.RUnlock()
	keys := make([]int64, len(mem.data))
	for key := range mem.data {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	sorted := make([]KeyValue, 0, len(mem.data))
	for _, sval := range keys {
		sorted = append(sorted, KeyValue{Key: sval, Value: mem.data[sval]})
	}
	return sorted, nil
}

func QuicksortAlgorithm(arr []int64) {
	if len(arr) < 2 {
		return
	}
	pivot := arr[0]
	Left := []int64{}
	Right := []int64{}
	for _, v := range arr {
		if v < pivot {
			Left = append(Left, v)
		}
		if v > pivot {
			Right = append(Right, v)
		}
	}
	QuicksortAlgorithm(Left)
	QuicksortAlgorithm(Right)
	copy(arr, append(append(Left, pivot), Right...))
}
