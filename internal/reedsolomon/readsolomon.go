package reedsolomon

import (
	"hash"

	"github.com/shogo82148/qrcode/internal/reedsolomon/element"
	"github.com/shogo82148/qrcode/internal/reedsolomon/poly"
)

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

func Decode(data []byte, twoS int) error {
	syndrome := make(poly.Poly, twoS)
	noError := true
	p := poly.NewPoly(data)
	for i := 0; i < twoS; i++ {
		ret := p.Eval(element.Exp(i))
		noError = noError && ret == 0
		syndrome[len(syndrome)-1-i] = ret
	}
	if noError {
		return nil
	}
	poly.EuclideanAlgorithm(poly.NewMonomial(twoS, element.One), syndrome, twoS)
	return nil
}
