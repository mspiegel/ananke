package ananke

import (
	"bytes"
	"testing"
)

func TestGet(t *testing.T) {
	leaf := Leaf{
		Root:   true,
		Epoch:  0,
		Keys:   [][]byte{[]byte{'f', 'o', 'o'}},
		Values: [][]byte{[]byte{'b', 'a', 'r'}},
	}
	val, success, err := leaf.Get([]byte{'f', 'o', 'o'})
	if bytes.Compare(val, []byte{'b', 'a', 'r'}) != 0 {
		t.Errorf("get on success incorrect value %v", val)
	}
	if !success {
		t.Error("get on success incorrect success")
	}
	if err != nil {
		t.Error("get on success incorrect error")
	}
	val, success, err = leaf.Get([]byte{'f', 'o'})
	if bytes.Compare(val, []byte{}) != 0 {
		t.Errorf("get on failure incorrect value %v", val)
	}
	if success {
		t.Error("get on failure incorrect success")
	}
	if err != nil {
		t.Error("get on failure incorrect error")
	}
}
func TestScan(t *testing.T) {
	leaf := Leaf{
		Root:   true,
		Epoch:  0,
		Keys:   [][]byte{[]byte{'f', 'o', 'o'}},
		Values: [][]byte{[]byte{'b', 'a', 'r'}},
	}
	scanner := func(key []byte, val []byte) error {
		if bytes.Compare(key, []byte{'f', 'o', 'o'}) != 0 {
			t.Errorf("scanner incorrect key %v", key)
		}
		if bytes.Compare(val, []byte{'b', 'a', 'r'}) != 0 {
			t.Errorf("scanner incorrect value %v", val)
		}
		return nil
	}
	err := leaf.Scan(scanner)
	if err != nil {
		t.Error("scanner incorrect error")
	}
}

func TestUpsert(t *testing.T) {
	var leaf Leaf
	var err error
	msg := Message{
		Op:   ASSIGN,
		Data: []byte{'f', 'o', 'o'},
	}
	leaf, err = leaf.Upsert([]byte{'w', 'o', 'r', 'l', 'd'}, msg)
	if err != nil {
		t.Error(err)
	}
	if len(leaf.Keys) != 1 {
		t.Error("Upsert failed: first key not inserted")
	}
	if len(leaf.Values) != 1 {
		t.Error("Upsert failed: first value not inserted")
	}
	msg = Message{
		Op:   ASSIGN,
		Data: []byte{'b', 'a', 'r'},
	}
	leaf, err = leaf.Upsert([]byte{'h', 'e', 'l', 'l', 'o'}, msg)
	if err != nil {
		t.Error(err)
	}
	if len(leaf.Keys) != 2 {
		t.Error("Upsert failed: second key not inserted")
	}
	if len(leaf.Values) != 2 {
		t.Error("Upsert failed: second value not inserted")
	}
	msg = Message{
		Op:   ASSIGN,
		Data: []byte{'b', 'a', 'z'},
	}
	leaf, err = leaf.Upsert([]byte{'w', 'o', 'r', 'l', 'd'}, msg)
	if err != nil {
		t.Error(err)
	}
	if len(leaf.Keys) != 2 {
		t.Error("Upsert failed: third key inserted")
	}
	if len(leaf.Values) != 2 {
		t.Error("Upsert failed: third value inserted")
	}
	if bytes.Compare(leaf.Keys[0], []byte{'h', 'e', 'l', 'l', 'o'}) != 0 {
		t.Error("First key is not 'hello'")
	}
	if bytes.Compare(leaf.Keys[1], []byte{'w', 'o', 'r', 'l', 'd'}) != 0 {
		t.Error("Second key is not 'world'")
	}
	if bytes.Compare(leaf.Values[0], []byte{'b', 'a', 'r'}) != 0 {
		t.Error("First value is not 'bar'")
	}
	if bytes.Compare(leaf.Values[1], []byte{'b', 'a', 'z'}) != 0 {
		t.Error("Second value is not 'baz'")
	}
}
