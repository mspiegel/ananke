package ananke

import "unsafe"

type Node struct {
	ID    int64
	Epoch int64
	Keys  [][]byte
	// leaf nodes have values
	Values [][]byte
	// internal nodes have message buffer and children
	Buffer   []Message
	Children []int64
	Internal bool
}

// NodeMinimumSize is the minimum amount of bytes
// to store a node on disk. ID, Epoch, Internal
const NodeMinimumSize = 2*unsafe.Sizeof(int64(0)) + unsafe.Sizeof(false)
