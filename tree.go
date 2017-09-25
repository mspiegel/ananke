package ananke

type ScannerFunc func(key []byte, val []byte) error
