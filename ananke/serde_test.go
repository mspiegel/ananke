package ananke

import (
	"reflect"
	"testing"
)

func TestRoundTripLeaf(t *testing.T) {
	var output Node
	input := Node{
		ID:     0xf,
		Epoch:  0xff,
		Keys:   [][]byte{[]byte("foo")},
		Values: [][]byte{[]byte("bar")},
	}
	data, err := input.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	err = output.UnmarshalBinary(data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(input, output) {
		t.Errorf("Expected %+v received %+v", input, output)
	}
}
