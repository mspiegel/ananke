package ananke

import "fmt"

type MemStorage struct {
	Store map[int64]Node
}

func (m *MemStorage) Write(node Node) error {
	m.Store[node.ID] = node
	return nil
}

func (m *MemStorage) Read(id int64) (Node, error) {
	node, ok := m.Store[id]
	if !ok {
		return Node{}, fmt.Errorf(NoSuchNodeError, id)
	}
	return node, nil
}
