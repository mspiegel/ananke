package ananke

import (
	"bytes"
	"sort"
)

type Leaf struct {
	Root   bool
	Epoch  int64
	Keys   [][]byte
	Values [][]byte
}

func (l Leaf) Get(key []byte) (val []byte, success bool, err error) {
	idx := sort.Search(len(l.Keys), func(i int) bool {
		return bytes.Compare(key, l.Keys[i]) <= 0
	})
	if idx < len(l.Keys) && bytes.Compare(l.Keys[idx], key) == 0 {
		val = l.Values[idx]
		success = true
	}
	return
}

func (l Leaf) Scan(s ScannerFunc) error {
	for i := range l.Keys {
		err := s(l.Keys[i], l.Values[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (l Leaf) Upsert(key []byte, msg Message) (result Leaf, err error) {
	var val []byte
	idx := sort.Search(len(l.Keys), func(i int) bool {
		return bytes.Compare(key, l.Keys[i]) <= 0
	})
	if idx < len(l.Keys) && bytes.Compare(l.Keys[idx], key) == 0 {
		val, err = msg.Apply(l.Values[idx])
		if err != nil {
			return
		}
		l.Values[idx] = val
	} else {
		val, err = msg.Create()
		if err != nil {
			return
		}
		l.Keys = append(l.Keys, nil)
		l.Values = append(l.Values, nil)
		copy(l.Keys[idx+1:], l.Keys[idx:])
		copy(l.Values[idx+1:], l.Values[idx:])
		l.Keys[idx] = key
		l.Values[idx] = val
	}
	result = l
	return
}
