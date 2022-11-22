package element

import "testing"

func FuzzAddZero(f *testing.F) {
	f.Add(uint8(0))
	f.Add(uint8(1))
	f.Add(uint8(255))

	// test a + 0 = 0 + a = a
	f.Fuzz(func(t *testing.T, i uint8) {
		a := Element(i)
		x := Add(a, Zero)
		y := Add(Zero, a)
		if x != a {
			t.Errorf("%02x: a + 0 != a", i)
		}
		if y != a {
			t.Errorf("%02x: 0 + a != a", i)
		}
	})
}

func FuzzMulOne(f *testing.F) {
	f.Add(uint8(0))
	f.Add(uint8(1))
	f.Add(uint8(255))

	// test a * 1 = 1 * a = a
	f.Fuzz(func(t *testing.T, i uint8) {
		a := Element(i)
		x := Mul(a, One)
		y := Mul(One, a)
		if x != a {
			t.Errorf("%02x: a * 1 != a: %02x", i, x)
		}
		if y != a {
			t.Errorf("%02x: 1 * a != a: %02x", i, y)
		}
	})
}

func FuzzAssociativeLawAdd(f *testing.F) {
	f.Add(uint8(0), uint8(1), uint8(2))
	f.Add(uint8(1), uint8(2), uint8(3))

	// test a + (b + c) = (a + b) + c
	f.Fuzz(func(t *testing.T, i, j, k uint8) {
		a := Element(i)
		b := Element(j)
		c := Element(k)
		x := Add(a, Add(b, c))
		y := Add(Add(a, b), c)
		if x != y {
			t.Errorf("a + (b + c) != (a + b) + c: %02x, %02x, %02x", i, j, k)
		}
	})
}

func FuzzAssociativeLawMul(f *testing.F) {
	f.Add(uint8(0), uint8(1), uint8(2))
	f.Add(uint8(1), uint8(2), uint8(3))

	// test a + (b + c) = (a + b) + c
	f.Fuzz(func(t *testing.T, i, j, k uint8) {
		a := Element(i)
		b := Element(j)
		c := Element(k)
		x := Mul(a, Mul(b, c))
		y := Mul(Mul(a, b), c)
		if x != y {
			t.Errorf("a(bc) != (ab)c: %02x, %02x, %02x", i, j, k)
		}
	})
}

func FuzzDistributiveLaw(f *testing.F) {
	f.Add(uint8(0), uint8(1), uint8(2))
	f.Add(uint8(1), uint8(2), uint8(3))

	// test a + (b + c) = (a + b) + c
	f.Fuzz(func(t *testing.T, i, j, k uint8) {
		a := Element(i)
		b := Element(j)
		c := Element(k)
		x := Mul(a, Add(b, c))
		y := Add(Mul(a, b), Mul(a, c))
		if x != y {
			t.Errorf("a(b + c) != ab + ac: %02x, %02x, %02x", i, j, k)
		}
	})
}
