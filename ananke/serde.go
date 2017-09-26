package ananke

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	multierror "github.com/mspiegel/go-multierror"
)

func (n *Node) MarshalBinary() (data []byte, err error) {
	return n.marshalBinaryLeaf()
}

func (n *Node) UnmarshalBinary(data []byte) error {
	return n.unmarshalBinaryLeaf(data)
}

func (n *Node) marshalBinaryLeaf() (data []byte, err error) {
	size := NodeMinimumSize
	// length of keys and values
	size += unsafe.Sizeof(int64(0))
	// length of each key and value
	size += 2 * uintptr(len(n.Keys)) * unsafe.Sizeof(int64(0))
	for i := range n.Keys {
		size += uintptr(len(n.Keys[i]))
		size += uintptr(len(n.Values[i]))
	}
	data = make([]byte, 0, size)
	buf := bytes.NewBuffer(data)
	// buf.Write() implementation never returns an error
	binary.Write(buf, binary.LittleEndian, n.Internal)
	binary.Write(buf, binary.LittleEndian, n.ID)
	binary.Write(buf, binary.LittleEndian, n.Epoch)
	binary.Write(buf, binary.LittleEndian, int64(len(n.Keys)))
	for i := range n.Keys {
		binary.Write(buf, binary.LittleEndian, int64(len(n.Keys[i])))
	}
	for i := range n.Values {
		binary.Write(buf, binary.LittleEndian, int64(len(n.Values[i])))
	}
	for i := range n.Keys {
		buf.Write(n.Keys[i])
	}
	for i := range n.Values {
		buf.Write(n.Values[i])
	}
	if uintptr(buf.Len()) != size {
		panic(fmt.Sprintf("Expected %d bytes and wrote %d bytes", size, buf.Len()))
	}
	data = data[:size]
	return
}

func (n *Node) unmarshalBinaryLeaf(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	var ignore bool
	var size int64
	err = multierror.Append(err, binary.Read(buf, binary.LittleEndian, &ignore))
	err = multierror.Append(err, binary.Read(buf, binary.LittleEndian, &n.ID))
	err = multierror.Append(err, binary.Read(buf, binary.LittleEndian, &n.Epoch))
	err = multierror.Append(err, binary.Read(buf, binary.LittleEndian, &size))
	n.Keys = make([][]byte, size)
	n.Values = make([][]byte, size)
	for i := range n.Keys {
		err = multierror.Append(err, binary.Read(buf, binary.LittleEndian, &size))
		n.Keys[i] = make([]byte, size)
	}
	for i := range n.Values {
		err = multierror.Append(err, binary.Read(buf, binary.LittleEndian, &size))
		n.Values[i] = make([]byte, size)
	}
	for i := range n.Keys {
		_, e := buf.Read(n.Keys[i])
		err = multierror.Append(err, e)
	}
	for i := range n.Values {
		_, e := buf.Read(n.Values[i])
		err = multierror.Append(err, e)
	}
	return
}
