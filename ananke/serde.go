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

func putInt64(buf *bytes.Buffer, v int64) {
	var b [8]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	buf.Write(b[:])
}

func putLengths(buf *bytes.Buffer, array [][]byte) {
	b := make([]byte, len(array)*8)
	idx := 0
	for i := range array {
		v := int64(len(array[i]))
		b[idx] = byte(v)
		b[idx+1] = byte(v >> 8)
		b[idx+2] = byte(v >> 16)
		b[idx+3] = byte(v >> 24)
		b[idx+4] = byte(v >> 32)
		b[idx+5] = byte(v >> 40)
		b[idx+6] = byte(v >> 48)
		b[idx+7] = byte(v >> 56)
		idx += 8
	}
	buf.Write(b)
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
	if n.Internal {
		buf.WriteByte(byte(1))
	} else {
		buf.WriteByte(byte(0))
	}
	putInt64(buf, n.ID)
	putInt64(buf, n.Epoch)
	putInt64(buf, int64(len(n.Keys)))
	putLengths(buf, n.Keys)
	putLengths(buf, n.Values)
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
