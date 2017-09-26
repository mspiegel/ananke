package ananke

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	multierror "github.com/mspiegel/go-multierror"
)

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
	writeInt64(buf, n.ID)
	writeInt64(buf, n.Epoch)
	writeInt64(buf, int64(len(n.Keys)))
	writeLengths(buf, n.Keys)
	writeLengths(buf, n.Values)
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
	reader := bytes.NewReader(data)
	var ignore bool
	var size int64
	err = multierror.Append(err, binary.Read(reader, binary.LittleEndian, &ignore))
	err = multierror.Append(err, binary.Read(reader, binary.LittleEndian, &n.ID))
	err = multierror.Append(err, binary.Read(reader, binary.LittleEndian, &n.Epoch))
	err = multierror.Append(err, binary.Read(reader, binary.LittleEndian, &size))
	n.Keys = make([][]byte, size)
	n.Values = make([][]byte, size)
	lengths := make([]uint64, size*2)
	keyTotal, valTotal, e := readLengths(reader, lengths)
	err = multierror.Append(err, e)
	readByteArray(data[len(data)-int(valTotal)-int(keyTotal):], n.Keys, lengths)
	readByteArray(data[len(data)-int(valTotal):], n.Values, lengths[size:])
	return
}

func writeInt64(buf *bytes.Buffer, v int64) {
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

func writeLengths(buf *bytes.Buffer, array [][]byte) {
	b := make([]byte, len(array)*int(unsafe.Sizeof(uint64(0))))
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

func readLengths(reader *bytes.Reader, lengths []uint64) (uint64, uint64, error) {
	b := make([]byte, len(lengths)*int(unsafe.Sizeof(uint64(0))))
	mid := len(lengths) / 2
	tot1 := uint64(0)
	tot2 := uint64(0)
	idx := 0
	_, err := reader.Read(b)
	if err != nil {
		return 0, 0, err
	}
	for i := range lengths {
		lengths[i] = uint64(b[idx]) | uint64(b[idx+1])<<8 |
			uint64(b[idx+2])<<16 | uint64(b[idx+3])<<24 |
			uint64(b[idx+4])<<32 | uint64(b[idx+5])<<40 |
			uint64(b[idx+6])<<48 | uint64(b[idx+7])<<56
		idx += 8
		if i < mid {
			tot1 += lengths[i]
		} else {
			tot2 += lengths[i]
		}
	}
	return tot1, tot2, nil
}

func readByteArray(src []byte, dest [][]byte, lengths []uint64) {
	idx := 0
	for i := range dest {
		dest[i] = src[idx : idx+int(lengths[i])]
		idx += int(lengths[i])
	}
}
