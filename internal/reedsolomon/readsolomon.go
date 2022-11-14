package reedsolomon

import "hash"

//go:generate go run gen/main.go

func New(n int) hash.Hash {
	if n < 2 {
		panic("negative")
	}
	if n >= len(coders) {
		panic("too large")
	}
	return coders[n]()
}
