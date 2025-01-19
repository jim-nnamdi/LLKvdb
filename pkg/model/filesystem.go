package model

import (
	"fmt"
	"sort"
	"sync"
)

var (
	ErrWALWrite      = "could not write to WAL file"
	ErrWALBatchWrite = "could not write Batch Data to WAL file"
	ErrMTBatchWrite  = "could not write Batch Data to memtable file"
	ErrSSTBatchWrite = "could not write Batch Data to sstable file"
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

func (Fsys *Filesys) ReadKeyRange(startkey int64, endkey int64) ([]KeyValue, error) {
	Fsys.mutex_t.Lock()
	defer Fsys.mutex_t.Unlock()
	var results []KeyValue
	for key, val := range Fsys.memtable.data {
		if key >= startkey && key <= endkey {
			results = append(results, KeyValue{key: key, value: val})
		}
	}

	for _, sstable := range Fsys.sstables {
		data, err := sstable.Load()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _, kv := range data {
			if startkey >= kv.key && endkey <= kv.key {
				results = append(results, kv)
			}
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].key < results[j].key
	})

	return results, nil
}

func (Fsys *Filesys) BatchPut(batch []map[int64]string) {
	Fsys.mutex_t.Lock()
	defer Fsys.mutex_t.Unlock()
	for i := 0; i < len(batch); i++ {
		for keys, vals := range batch {
			if wal := Fsys.aheadLog.Write(int64(keys), vals[int64(keys)]); wal != nil {
				fmt.Println()
			}
			Fsys.memtable.Put(int64(keys), vals[int64(keys)])
			if len(Fsys.memtable.data) > Fsys.maxmemsize {
				memflush := Fsys.memtable.Flush()
				sstable := Newsstable(fmt.Sprintf("sstable-%d", keys))
				_ = sstable.Write(memflush)
				Fsys.sstables = append(Fsys.sstables, sstable)
			}
		}
	}
}

func (Fsys *Filesys) Delete(key int64) {
	Fsys.memtable.Delete(key)
}
