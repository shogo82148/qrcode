package element

// Element is GF(2^8) performed modulo x^8 + x^4 + x^3 + x^2 + 1.
type Element uint8

const Zero = Element(0)
const One = Element(1)

// Add returns x + y.
func Add(x, y Element) Element {
	return x ^ y
}

// Mul returns x * y.
func Mul(x, y Element) Element {
	if x == Zero || y == Zero {
		return Zero
	}
	xx := logTable[x]
	yy := logTable[y]
	zz := (xx + yy) % 255
	return expTable[zz]
}

func Log(x Element) int {
	if x == Zero {
		panic("element: log of zero")
	}
	return logTable[x]
}

func Exp(n int) Element {
	return expTable[n%255]
}

func Inv(x Element) Element {
	if x == Zero {
		panic("element: divided by zero")
	}
	return expTable[255-logTable[x]]
}

// AddMulExp returns x + a^y * a^z
func AddMulExp(x Element, y, z int) Element {
	return Add(x, expTable[(y+z)%255])
}
