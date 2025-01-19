package tests

import (
	"testing"

	"github.com/jim-nnamdi/Lkvs/pkg/model"
)

func TestMemTablePutAndRead(t *testing.T) {
	memtable := model.Newmemtable()
	memtable.Put(1, "jim")
	memtable.Put(2, "moniepoint")

	val, exists := memtable.Read(1)
	if !exists || val != "jim" {
		t.Errorf("Expected key 1 to have value 'jim' got '%s'", val)
	}

	val, exists = memtable.Read(2)
	if !exists || val != "moniepoint" {
		t.Errorf("Expected key 1 to have value 'moniepoint' got '%s'", val)
	}

	_, exists = memtable.Read(3)
	if exists {
		t.Errorf("Expected key 3 not to exist, but it was found")
	}
}

func TestMemTableDelete(t *testing.T) {
	memtable := model.Newmemtable()
	memtable.Put(1, "samuel")
	memtable.Delete(1)
	_, exists := memtable.Read(1)
	if exists {
		t.Errorf("Expected key 1 not to exist again but it was found")
	}
}

func TestMemTableFlushSSTable(t *testing.T) {
	memtable := model.Newmemtable()
	memtable.Put(34, "a")
	memtable.Put(37, "b")
	flushed := memtable.Flush()
	expected := []model.KeyValue{
		{Key: 34, Value: "a"},
		{Key: 37, Value: "b"},
	}
	if len(flushed) != len(expected) {
		t.Fatalf("expected flushed len of '%d' got '%d'", len(expected), len(flushed))
	}

	for i, kv := range flushed {
		if kv != expected[i] {
			t.Errorf("Expected flushed['%d'] = %+v got %+v", i, expected[i], kv)
		}
	}
}
