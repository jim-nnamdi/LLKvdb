package tests

import (
	"testing"

	"github.com/jim-nnamdi/Lkvs/pkg/model"
)

func TestWalWrite(t *testing.T) {
	wal, _ := model.NewWAL("testwal.txt")
	memt := model.Newmemtable()
	wal.Write(1, "santa claus")
	wal.Write(2, "elliptic curves")
	val, exist := memt.Dump()
	if exist != nil {
		t.Fatalf("in memory does not exist: '%s'", exist)
	}
	for _, vals := range val {
		if vals.Key <= 0 {
			t.Errorf("key should be a non negative value")
		}
	}
}

func TestWalDelete(t *testing.T) {
	wal, _ := model.NewWAL("testwal.txt")
	wal.Write(5, "five")
	wal.Write(7, "seven")
	exists, err := wal.Delete(5)
	if err != nil || !exists {
		t.Errorf("invalid data selection because data doesnt exist")
	}
	val, _ := wal.WALDump(5)
	if val == "five" {
		t.Errorf("invalid call to non-existent value. value should have been removed already")
	}
}
