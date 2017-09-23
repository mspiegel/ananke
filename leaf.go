package ananke

import (
	"bytes"
	"sort"
)

type Leaf struct {
	Keys   [][]byte
	Values [][]byte
}

func (l Leaf) Upsert(key []byte, msg Message) (Leaf, error) {
	idx := sort.Search(len(l.Keys), func(i int) bool {
		return bytes.Compare(key, l.Keys[i]) <= 0
	})
	if idx < len(l.Keys) && bytes.Compare(l.Keys[idx], key) == 0 {
		val, err := msg.Apply(l.Values[idx])
		if err != nil {
			return l, err
		}
		l.Values[idx] = val
	} else {
		val, err := msg.Create()
		if err != nil {
			return l, err
		}
		l.Keys = append(l.Keys, nil)
		l.Values = append(l.Values, nil)
		copy(l.Keys[idx+1:], l.Keys[idx:])
		copy(l.Values[idx+1:], l.Values[idx:])
		l.Keys[idx] = key
		l.Values[idx] = val
	}
	return l, nil
}
