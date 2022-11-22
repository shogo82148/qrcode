package reedsolomon

import (
	"fmt"
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
	// from https://github.com/zxing/zxing/blob/99e9b34f5afc21fdaeead283d5ed0bc1314cbec1/core/src/main/java/com/google/zxing/common/reedsolomon/ReedSolomonDecoder.java#L49-L86

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
	sigma, omega, err := poly.EuclideanAlgorithm(poly.NewMonomial(twoS, element.One), syndrome, twoS)
	if err != nil {
		return fmt.Errorf("reedsolomon: failed to decode: %w", err)
	}
	errorLocations := findErrorLocations(sigma)
	errorMagnitudes := findErrorMagnitudes(omega, errorLocations)

	for i := range errorLocations {
		pos := len(data) - 1 - element.Log(errorLocations[i])
		data[pos] = byte(element.Add(element.Element(data[pos]), errorMagnitudes[i]))
	}
	return nil
}

func findErrorLocations(sigma poly.Poly) []element.Element {
	ret := []element.Element{}
	for i := 1; i < 256; i++ {
		e := element.Element(i)
		if sigma.Eval(e) == element.Zero {
			ret = append(ret, element.Inv(e))
		}
	}
	return ret
}

func findErrorMagnitudes(omega poly.Poly, errorLocations []element.Element) []element.Element {
	ret := make([]element.Element, len(errorLocations))
	for i, loc1 := range errorLocations {
		xiInverse := element.Inv(loc1)
		denominator := element.One
		for j, loc2 := range errorLocations {
			if i != j {
				term := element.Mul(loc2, xiInverse)
				term = element.Add(term, element.One)
				denominator = element.Mul(denominator, term)
			}
		}
		ret[i] = element.Mul(omega.Eval(xiInverse), element.Inv(denominator))
	}
	return ret
}
