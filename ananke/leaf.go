package ananke

import (
	"bytes"
	"sort"
)

func (l *Node) GetLeaf(key []byte) (val []byte, success bool, err error) {
	idx := sort.Search(len(l.Keys), func(i int) bool {
		return bytes.Compare(key, l.Keys[i]) <= 0
	})
	if idx < len(l.Keys) && bytes.Compare(l.Keys[idx], key) == 0 {
		val = l.Values[idx]
		success = true
	}
	return
}

func (l *Node) ScanLeaf(scan ScannerFunc) (err error) {
	for i := range l.Keys {
		err = scan(l.Keys[i], l.Values[i])
		if err != nil {
			return
		}
	}
	return
}

func (l *Node) UpsertLeaf(tree *Tree, msg Message) (sib *Node, err error) {
	var val []byte
	idx := sort.Search(len(l.Keys), func(i int) bool {
		return bytes.Compare(msg.Key, l.Keys[i]) <= 0
	})
	if idx < len(l.Keys) && bytes.Compare(l.Keys[idx], msg.Key) == 0 {
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
		if len(l.Keys) == tree.MaxItems {
			sib = l.splitLeaf(tree, msg.Key, val, idx)
		} else {
			l.putLeaf(msg.Key, val, idx)
		}
	}
	return
}

func (left *Node) splitLeaf(tree *Tree, key []byte, val []byte, idx int) (right *Node) {
	right = new(Node)
	split := len(left.Keys) / 2
	right.Keys = left.Keys[split:]
	right.Values = left.Values[split:]
	left.Keys = left.Keys[:split]
	left.Values = left.Values[:split]
	if idx > split {
		right.putLeaf(key, val, idx-split)
	} else {
		left.putLeaf(key, val, idx)
	}
	return
}

func (l *Node) putLeaf(key []byte, val []byte, idx int) {
	l.Keys = append(l.Keys, nil)
	l.Values = append(l.Values, nil)
	copy(l.Keys[idx+1:], l.Keys[idx:])
	copy(l.Values[idx+1:], l.Values[idx:])
	l.Keys[idx] = key
	l.Values[idx] = val
}
