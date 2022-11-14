package reedsolomon

// element is GF(2^8) performed modulo x^8 + x^4 + x^3 + x^2 + 1.
type element uint8

const zero = element(0)
const one = element(1)

// add returns x + y.
func add(x, y element) element {
	return x ^ y
}

// mul returns x * y.
func mul(x, y element) element {
	if x == zero || y == zero {
		return zero
	}
	xx := logTable[x]
	yy := logTable[y]
	zz := (xx + yy) % 255
	return expTable[zz]
}

// AddMul sets v = x + a^y * a^z
func (v *element) AddMulExp(x element, y, z int) {
	*v = add(x, expTable[(y+z)%255])
}
