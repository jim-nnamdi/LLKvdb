package model

import "sync"

type Memtable struct {
	data  map[int64]string
	mutex sync.RWMutex
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

func (mem *Memtable) Flush() map[int64]string {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	flushed := mem.data
	mem.data = make(map[int64]string)
	return flushed
}
