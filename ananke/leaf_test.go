package ananke

import (
	"bytes"
	"testing"
)

func TestGetLeaf(t *testing.T) {
	leaf := Node{
		Keys:   [][]byte{[]byte("foo")},
		Values: [][]byte{[]byte("bar")},
	}
	val, success, err := leaf.GetLeaf([]byte("foo"))
	if bytes.Compare(val, []byte("bar")) != 0 {
		t.Errorf("get on success incorrect value %v", val)
	}
	if !success {
		t.Error("get on success incorrect success")
	}
	if err != nil {
		t.Error("get on success incorrect error")
	}
	val, success, err = leaf.GetLeaf([]byte("fo"))
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
func TestScanLeaf(t *testing.T) {
	leaf := Node{
		Keys:   [][]byte{[]byte("foo")},
		Values: [][]byte{[]byte("bar")},
	}
	scanner := func(key []byte, val []byte) error {
		if bytes.Compare(key, []byte("foo")) != 0 {
			t.Errorf("scanner incorrect key %v", key)
		}
		if bytes.Compare(val, []byte("bar")) != 0 {
			t.Errorf("scanner incorrect value %v", val)
		}
		return nil
	}
	err := leaf.ScanLeaf(scanner)
	if err != nil {
		t.Error("scanner incorrect error")
	}
}

func TestUpsertLeaf(t *testing.T) {
	tree, err := Builder(100).Build()
	if err != nil {
		t.Error(err)
	}
	leaf := Node{}
	msg := Message{
		Key:  []byte("world"),
		Op:   ASSIGN,
		Data: []byte("foo"),
	}
	sibling, err := leaf.UpsertLeaf(tree, msg)
	if err != nil {
		t.Error(err)
	}
	if len(leaf.Keys) != 1 {
		t.Error("Upsert failed: first key not inserted")
	}
	if len(leaf.Values) != 1 {
		t.Error("Upsert failed: first value not inserted")
	}
	if sibling != nil {
		t.Error("Upsert failed. Sibling created")
	}
	msg = Message{
		Key:  []byte("hello"),
		Op:   ASSIGN,
		Data: []byte("bar"),
	}
	sibling, err = leaf.UpsertLeaf(tree, msg)
	if err != nil {
		t.Error(err)
	}
	if len(leaf.Keys) != 2 {
		t.Error("Upsert failed: second key not inserted")
	}
	if len(leaf.Values) != 2 {
		t.Error("Upsert failed: second value not inserted")
	}
	if sibling != nil {
		t.Error("Upsert failed. Sibling created")
	}
	msg = Message{
		Key:  []byte("world"),
		Op:   ASSIGN,
		Data: []byte("baz"),
	}
	sibling, err = leaf.UpsertLeaf(tree, msg)
	if err != nil {
		t.Error(err)
	}
	if len(leaf.Keys) != 2 {
		t.Error("Upsert failed: third key inserted")
	}
	if len(leaf.Values) != 2 {
		t.Error("Upsert failed: third value inserted")
	}
	if bytes.Compare(leaf.Keys[0], []byte("hello")) != 0 {
		t.Error("First key is not 'hello'")
	}
	if bytes.Compare(leaf.Keys[1], []byte("world")) != 0 {
		t.Error("Second key is not 'world'")
	}
	if bytes.Compare(leaf.Values[0], []byte("bar")) != 0 {
		t.Error("First value is not 'bar'")
	}
	if bytes.Compare(leaf.Values[1], []byte("baz")) != 0 {
		t.Error("Second value is not 'baz'")
	}
	if sibling != nil {
		t.Error("Upsert failed. Sibling created")
	}
}

func TestSplitLeafInsertLeft(t *testing.T) {
	left := Node{
		Keys:   [][]byte{[]byte("hello"), []byte("world")},
		Values: [][]byte{[]byte("foo"), []byte("bar")},
	}
	msg := Message{
		Key:  []byte("middle"),
		Op:   ASSIGN,
		Data: []byte("baz"),
	}
	tree := &Tree{
		MaxItems: 2,
	}
	right, err := left.UpsertLeaf(tree, msg)
	if err != nil {
		t.Error(err)
	}
	if len(left.Keys) != 2 {
		t.Errorf("left node has %d elements", len(left.Keys))
	}
	if len(right.Keys) != 1 {
		t.Errorf("left node has %d elements", len(right.Keys))
	}
}

func TestSplitLeafInsertRight(t *testing.T) {
	left := Node{
		Keys:   [][]byte{[]byte("hello"), []byte("world")},
		Values: [][]byte{[]byte("foo"), []byte("bar")},
	}
	msg := Message{
		Key:  []byte("yesss"),
		Op:   ASSIGN,
		Data: []byte("baz"),
	}
	tree := &Tree{
		MaxItems: 2,
	}
	right, err := left.UpsertLeaf(tree, msg)
	if err != nil {
		t.Error(err)
	}
	if len(left.Keys) != 1 {
		t.Errorf("left node has %d elements", len(left.Keys))
	}
	if len(right.Keys) != 2 {
		t.Errorf("left node has %d elements", len(right.Keys))
	}
}
