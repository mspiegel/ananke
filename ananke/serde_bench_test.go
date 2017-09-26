package ananke

import (
	"math/rand"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rng.Intn(len(letterBytes))]
	}
	return b
}

var rng *rand.Rand
var leaf Node
var leafBytes []byte

func setup(items int) {
	rng = rand.New(rand.NewSource(7))
	leaf.ID = rng.Int63()
	leaf.Epoch = rng.Int63()
	leaf.Keys = make([][]byte, items)
	leaf.Values = make([][]byte, items)
	for i := range leaf.Keys {
		leaf.Keys[i] = randStringBytes(rng.Int() % 100)
	}
	for i := range leaf.Values {
		leaf.Values[i] = randStringBytes(rng.Int() % 100)
	}
	leafBytes, _ = leaf.MarshalBinary()
}

func BenchmarkMarshalBinaryLeaf10(b *testing.B)  { benchmarkMarshalBinaryLeaf(b, 10) }
func BenchmarkMarshalBinaryLeaf100(b *testing.B) { benchmarkMarshalBinaryLeaf(b, 100) }

func benchmarkMarshalBinaryLeaf(b *testing.B, items int) {
	setup(items)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		leaf.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinaryLeaf10(b *testing.B)  { benchmarkUnmarshalBinaryLeaf(b, 10) }
func BenchmarkUnmarshalBinaryLeaf100(b *testing.B) { benchmarkUnmarshalBinaryLeaf(b, 100) }

func benchmarkUnmarshalBinaryLeaf(b *testing.B, items int) {
	var node Node
	setup(items)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		node.UnmarshalBinary(leafBytes)
	}
}
