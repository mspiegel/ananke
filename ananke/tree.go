package ananke

import (
	"errors"
	"math"
	"sync/atomic"

	multierror "github.com/mspiegel/go-multierror"
)

type ScannerFunc func(key []byte, val []byte) error

var (
	ErrEpsilonOutOfBounds = errors.New("epsilon must be between 0.0 and 1.0")
	ErrNegativeCapacity   = errors.New("capacity must be positive integer")
)

type Tree struct {
	ReadRoot  atomic.Value
	WriteRoot *Node
	MaxKeys   int
	MaxItems  int
	Epoch     int64
	NextID    int64
	Store     Storage
}

type TreeBuilder struct {
	Cap   int
	Eps   float64
	Store Storage
	Err   error
}

func Builder(capacity int) *TreeBuilder {
	b := new(TreeBuilder)
	b.Cap = capacity
	b.Eps = 0.5
	b.Store = new(MemStorage)
	if capacity < 0 {
		b.Err = ErrNegativeCapacity
	}
	return b
}

func (b *TreeBuilder) Epsilon(eps float64) *TreeBuilder {
	if eps < 0.0 || eps > 1.0 {
		b.Err = multierror.Append(b.Err, ErrEpsilonOutOfBounds)
	} else {
		b.Eps = eps
	}
	return b
}

func (b *TreeBuilder) Storage(store Storage) *TreeBuilder {
	b.Store = store
	return b
}

func (b *TreeBuilder) Build() (*Tree, error) {
	if b.Err != nil {
		return nil, b.Err
	}
	tree := new(Tree)
	tree.MaxKeys = int(math.Pow(float64(b.Cap), b.Eps))
	tree.MaxItems = b.Cap - tree.MaxKeys
	tree.NextID = 1
	tree.Store = b.Store
	return tree, nil
}
