package ananke

import (
	"bytes"
	"testing"
)

func TestAssign(t *testing.T) {
	msg := Message{
		Op:   ASSIGN,
		Data: []byte{'f', 'o', 'o'},
	}
	val, err := msg.Create()
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(val, []byte{'f', 'o', 'o'}) != 0 {
		t.Errorf("ASSIGN creation failed: %v", val)
	}
	val, err = msg.Apply([]byte{'b', 'a', 'r'})
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(val, []byte{'f', 'o', 'o'}) != 0 {
		t.Errorf("ASSIGN creation failed: %v", val)
	}
}

func TestInvalidOperation(t *testing.T) {
	msg := Message{
		Op:   -1,
		Data: nil,
	}
	_, err := msg.Create()
	if err == nil {
		t.Error("Expected error on invalid operation")
	}
	_, err = msg.Apply([]byte{'f', 'o', 'o'})
	if err == nil {
		t.Error("Expected error on invalid operation")
	}
}
