package ananke

import (
	"bytes"
	"testing"
)

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
