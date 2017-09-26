package ananke

import "fmt"

type Operation int

const UnknownOperation = "Unknown message type %d"

const (
	ASSIGN Operation = iota
)

type Message struct {
	Key  []byte
	Op   Operation
	Data []byte
}

func (msg Message) Apply(value []byte) ([]byte, error) {
	switch msg.Op {
	case ASSIGN:
		value = msg.Data
	default:
		return nil, fmt.Errorf(UnknownOperation, msg.Op)
	}
	return value, nil
}

func (msg Message) Create() ([]byte, error) {
	switch msg.Op {
	case ASSIGN:
		return msg.Data, nil
	default:
		return nil, fmt.Errorf(UnknownOperation, msg.Op)
	}
}
