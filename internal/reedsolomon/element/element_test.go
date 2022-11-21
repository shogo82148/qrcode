package element

import "testing"

func TestInverse(t *testing.T) {
	for i := 1; i < 256; i++ {
		x := Element(i)
		y := Inv(x)
		z := Mul(x, y)
		if z != One {
			t.Errorf("%08b * %08b is not one", x, y)
		}
	}
}
