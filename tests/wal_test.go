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
