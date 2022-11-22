package poly

import (
	"errors"

	"github.com/shogo82148/qrcode/internal/reedsolomon/element"
)

// Poly is a polynomial that is p[0]*x^(len(p)-1) + p[1]*x^(len(p)-2) + ... + p[len(p)-1].
type Poly []element.Element

func Zero() Poly {
	return Poly{}
}

func One() Poly {
	return Poly{element.One}
}

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

func (p Poly) Add(q Poly) Poly {
	r := p
	if len(p) < len(q) {
		r = make(Poly, len(q))
		copy(r[len(q)-len(p):], p)
	}
	for i := range q {
		r[len(r)-i-1] = element.Add(r[len(r)-i-1], q[len(q)-i-1])
	}
	return r
}

func (p Poly) Mul(q Poly) Poly {
	ret := make(Poly, len(p)+len(q)-1)
	for i, pe := range p {
		for j, qe := range q {
			ret[i+j] = element.Add(ret[i+j], element.Mul(pe, qe))
		}
	}
	return ret
}

func (p Poly) MulMonomial(degree int, coefficient element.Element) Poly {
	if coefficient == element.Zero {
		return Zero()
	}
	ret := make(Poly, degree+len(p))
	for i, e := range p {
		ret[i] = element.Mul(e, coefficient)
	}
	return ret
}

func (p Poly) MulElement(v element.Element) Poly {
	for i := range p {
		p[i] = element.Mul(p[i], v)
	}
	return p
}

func EuclideanAlgorithm(a, b Poly, R int) (sigma Poly, omega Poly, err error) {
	// from https://github.com/zxing/zxing/blob/99e9b34f5afc21fdaeead283d5ed0bc1314cbec1/core/src/main/java/com/google/zxing/common/reedsolomon/ReedSolomonDecoder.java#L88-L141
	// Assume a's degree is >= b's
	if a.Degree() < b.Degree() {
		a, b = b, a
	}

	rLast := a
	r := b
	tLast := Zero()
	t := One()
	for 2*r.Degree() >= R {
		r, rLast = rLast, r
		t, tLast = tLast, t

		// Divide rLastLast by rLast, with quotient in q and remainder in r
		q := Poly{}
		denominatorLeadingTerm := rLast.Coefficient(rLast.Degree())
		dltInverse := element.Inv(denominatorLeadingTerm)
		for r.Degree() >= rLast.Degree() {
			degreeDiff := r.Degree() - rLast.Degree()
			scale := element.Mul(r.Coefficient(r.Degree()), dltInverse)
			q = q.Add(NewMonomial(degreeDiff, scale))
			r = r.Add(rLast.MulMonomial(degreeDiff, scale))
		}
		t = q.Mul(tLast).Add(t)
	}

	sigmaTildeAtZero := t.Coefficient(0)
	if sigmaTildeAtZero == element.Zero {
		return nil, nil, errors.New("sigmaTilde(0) was zero")
	}
	inv := element.Inv(sigmaTildeAtZero)
	sigma = t.MulElement(inv)
	omega = r.MulElement(inv)
	return
}
