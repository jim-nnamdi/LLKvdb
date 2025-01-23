package model

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrWALWrite          = "could not write to WAL file"
	ErrWALBatchWrite     = "could not write Batch Data to WAL file"
	ErrMTBatchWrite      = "could not write Batch Data to memtable file"
	ErrSSTBatchWrite     = "could not write Batch Data to sstable file"
	ErrCompactionSSTable = "Error during compaction process"
)

type Filesys struct {
	memtable   *AVLTree
	sstables   []*SSTable
	aheadLog   *WAL
	compacts   *MinHeap
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
		memtable:   NewAVLTree(),
		sstables:   []*SSTable{},
		aheadLog:   wal,
		maxmemsize: maxsize,
	}
}

func (Fsys *Filesys) sstableDirCount(dirname string) int {
	files, err := os.ReadDir(dirname)
	if err != nil {
		fmt.Println(err)
	}
	count := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "sstable-") {
			count++
		}
	}
	if count >= 5 {
		names := make([]string, len(files))
		for i, name := range files {
			names[i] = name.Name()
		}
		// TODO: hash the value passed as compact index
		Fsys.compacts.MergeSSTables(names, fmt.Sprintf("sstable-compact-%d", count-1))
	}

	for _, zname := range files {
		build_files := dirname + "/" + zname.Name()
		if strings.HasPrefix(build_files, "sstable-") {
			if err := os.Remove(build_files); err != nil {
				fmt.Println("error removing sstable", err)
			}
		}
	}
	return count
}

func (Fsys *Filesys) Put(key int64, value string) {
	Fsys.mutex_t.Lock()
	defer Fsys.mutex_t.Unlock()
	walerror := Fsys.aheadLog.Write(key, value)
	if walerror != nil {
		fmt.Println(walerror)
		fmt.Printf("%s\n", ErrWALWrite)
	}
	// perform compaction first before writing to AVLTree
	sstable_c := Fsys.sstableDirCount(".")
	fmt.Printf("sstables are %d in number\n", sstable_c+1)
	Fsys.memtable.Insert(key, value)
	Fsys.memtable.InOrderTraversal()
	if Fsys.memtable.Size > Fsys.maxmemsize {
		memflush := Fsys.memtable.Flush()
		sstable_t := Newsstable(fmt.Sprintf("sstable-%d", len(Fsys.sstables)))
		_ = sstable_t.Write(memflush)
		Fsys.sstables = append(Fsys.sstables, sstable_t)
	}
}

func (Fsys *Filesys) Read(key int64) (string, bool) {
	Fsys.mutex_t.RLock()
	defer Fsys.mutex_t.RUnlock()
	if val, exists := Fsys.memtable.Search(key); exists {
		fmt.Println("val", val)
		return val, true
	}
	for i := len(Fsys.sstables) - 1; i >= 0; i-- {
		block, _ := Fsys.sstables[i].Get(key)
		if val, exists := block[key]; exists {
			fmt.Println("val", val)
			return val, true
		}
	}
	extra := Fsys.memtable.InOrderTraversal()
	for _, kv := range extra {
		if key == kv.Key {
			fmt.Println("val", kv.Value)
			return kv.Value, true
		}
	}

	return emptystring(), false
}

func (Fsys *Filesys) ReadKeyRange(startkey int64, endkey int64) ([]KeyValue, error) {
	Fsys.mutex_t.Lock()
	defer Fsys.mutex_t.Unlock()
	var results []KeyValue
	for key, val := range Fsys.memtable.ReadKeyRange(startkey, endkey) {
		if int64(key) >= startkey && int64(key) <= endkey {
			results = append(results, val)
		}
	}

	for _, sstable := range Fsys.sstables {
		data, err := sstable.Load()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _, kv := range data {
			if startkey >= kv.Key && endkey <= kv.Key {
				results = append(results, kv)
			}
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results, nil
}

func (Fsys *Filesys) BatchPut(keys []int64, vals []string) {
	Fsys.mutex_t.Lock()
	defer Fsys.mutex_t.Unlock()
	if len(keys) != len(vals) {
		fmt.Println("Length of keys and values are different. they should be the same")
		return
	}
	for i := 0; i < len(keys); i++ {
		keys := keys[i]
		vals := vals[i]
		if wal := Fsys.aheadLog.Write(keys, vals); wal != nil {
			fmt.Println("error writing batch data to disk")
			return
		}
		Fsys.memtable.Insert(keys, vals)
		if Fsys.memtable.Size > Fsys.maxmemsize {
			flushed := Fsys.memtable.Flush()
			sstable := Newsstable(fmt.Sprintf("sstable-%d", keys))
			_ = sstable.Write(flushed)
			Fsys.sstables = append(Fsys.sstables, sstable)
		}
	}
}

func (Fsys *Filesys) Recover() {
	wal, err := Fsys.aheadLog.Replay()
	if err != nil {
		fmt.Println(err)
	}
	for _, vals := range wal {
		Fsys.memtable.Insert(vals.Key, vals.Value)
	}
	if Fsys.memtable.Size > Fsys.maxmemsize {
		Fsys.memtable.Flush()
	}
}

func (Fsys *Filesys) Delete(key int64) {
	Fsys.aheadLog.Delete(key)
	Fsys.memtable.Delete(key)
}

func (Fsys *Filesys) Compact(sstables []string, ofile string) {
	if err := Fsys.compacts.MergeSSTables(sstables, ofile); err != nil {
		fmt.Println("compaction error", err)
	}
}

func (Fsys *Filesys) Compaction(sst *SSTable) error {
	/*
		when multiple tables are created for the sstables
		which happens when the memtable flushes info out
		because it has reached a certain memory threshold
		which would of course be declared by the system ..
		[Dynamic allocation] [in this case am statically allocating memory]
		we would need to merge multiple sstables as one to minimize amount
		of sstable-(N) files created

		--- current state of my program [85 % Efficiency]
		WAL -> memtable [ 4 / 1024 * 1024] -> Flush -> sstable-0
		WAL -> memtable [ 4 / 1024 * 1024] -> Flush -> sstable-1
		WAL -> memtable [ N / 1024 * 1024] -> Flush-N -> sstable-N
	*/
	var all []KeyValue
	for _, sstable := range Fsys.sstables {
		data, err := sstable.Load()
		if err != nil {
			fmt.Println(err)
			return err
		}
		all = append(all, data...)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Key < all[j].Key })
	merge := make([]KeyValue, 0)
	seenk := make(map[int64]bool)
	for i := len(all) - 1; i > 0; i-- {
		if !seenk[all[i].Key] {
			merge = append(merge, all[i])
			seenk[all[i].Key] = true
		}
	}
	sort.Slice(merge, func(i, j int) bool { return merge[i].Key < merge[j].Key })
	newsspath := Newsstable(fmt.Sprintf("sstable-compact-%d", time.Now().UnixNano()))
	if err := newsspath.Write(merge); err != nil {
		fmt.Println(ErrCompactionSSTable)
		return err
	}
	Fsys.sstables = []*SSTable{newsspath}
	return nil
}
