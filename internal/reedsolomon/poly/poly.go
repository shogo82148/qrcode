package poly

import "github.com/shogo82148/qrcode/internal/reedsolomon/element"

// Poly is a polynomial that is p[0]*x^(len(p)-1) + p[1]*x^(len(p)-2) + ... + p[len(p)-1].
type Poly []element.Element

func NewPoly(data []byte) Poly {
	p := make([]element.Element, len(data))
	for i, b := range data {
		p[i] = element.Element(b)
	}
	return p
}

func NewMonomial(degree int, coefficient element.Element) Poly {
	p := make(Poly, degree+1)
	p[0] = coefficient
	return p
}

func (p Poly) Eval(x element.Element) element.Element {
	ret := element.Zero
	for _, b := range p {
		ret = element.Mul(ret, x)
		ret = element.Add(ret, b)
	}
	return ret
}

func (p Poly) Degree() int {
	for i, e := range p {
		if e != element.Zero {
			return len(p) - i - 1
		}
	}
	return 0
}

func (p Poly) Coefficient(degree int) element.Element {
	if degree >= len(p) {
		return element.Zero
	}
	return p[len(p)-degree-1]
}

func EuclideanAlgorithm(a, b Poly, R int) (Poly, Poly, error) {
	// from https://github.com/zxing/zxing/blob/bc88dd15c879183cb59cc8482772bc2d185a2ad6/core/src/main/java/com/google/zxing/common/reedsolomon/ReedSolomonDecoder.java
	// Assume a's degree is >= b's
	if a.Degree() < b.Degree() {
		a, b = b, a
	}

	rLast := a
	r := b
	// tLast := poly{zero}
	// t := poly{one}
	for 2*r.Degree() >= R {
		r, rLast = rLast, r
		// t, tLast = tLast, t

		// Divide rLastLast by rLast, with quotient in q and remainder in r
		// denominatorLeadingTerm := rLast.coefficient(rLast.degree())
		// dltInverse := inverse(denominatorLeadingTerm)
		// for r.degree() >= rLast.degree() {
		// 	degreeDiff := r.degree() - rLast.degree()
		// 	scale := mul(r.coefficient(r.degree()), dltInverse)
		// 	_ = degreeDiff
		// 	_ = scale
		// }
	}
	return nil, nil, nil
}
