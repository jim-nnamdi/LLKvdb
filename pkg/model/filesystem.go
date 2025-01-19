package model

import (
	"fmt"
	"sync"
)

var (
	ErrWALWrite = "could not write to WAL file"
)

type Filesys struct {
	memtable   *Memtable
	sstables   []*SSTable
	aheadLog   *WAL
	mutex_t    sync.RWMutex
	maxmemsize int
}

func NewFilesys(walLoc string, maxsize int) *Filesys {
	wal, err := NewWAL(walLoc)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%s\n", "Write ahead Log File error")
	}
	return &Filesys{
		memtable:   Newmemtable(),
		sstables:   []*SSTable{},
		aheadLog:   wal,
		maxmemsize: maxsize,
	}
}

func (Fsys *Filesys) Put(key int64, value string) {
	Fsys.mutex_t.Lock()
	defer Fsys.mutex_t.Unlock()
	walerror := Fsys.aheadLog.Write(key, value)
	if walerror != nil {
		fmt.Println(walerror)
		fmt.Printf("%s\n", ErrWALWrite)
	}
	Fsys.memtable.Put(key, value)
	if len(Fsys.memtable.data) > Fsys.maxmemsize {
		memflush := Fsys.memtable.Flush()
		sstable_t := Newsstable(fmt.Sprintf("sstable-%d", len(Fsys.sstables)))
		_ = sstable_t.Write(memflush)
		Fsys.sstables = append(Fsys.sstables, sstable_t)
	}
}

func (Fsys *Filesys) Read(key int64) (string, bool) {
	Fsys.mutex_t.RLock()
	defer Fsys.mutex_t.RUnlock()
	if val, exists := Fsys.memtable.Read(key); exists {
		return val, true
	}
	for i := len(Fsys.sstables) - 1; i >= 0; i-- {
		block, _ := Fsys.sstables[i].Get(key)
		if val, exists := block[key]; exists {
			return val, true
		}
	}
	return emptystring(), false
}
