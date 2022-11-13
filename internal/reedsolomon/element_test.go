package reedsolomon

import (
	"testing"
)

func FuzzAddZero(f *testing.F) {
	f.Add(uint8(0))
	f.Add(uint8(1))
	f.Add(uint8(255))

	// test a + 0 = 0 + a = a
	f.Fuzz(func(t *testing.T, i uint8) {
		a := element(i)
		x := add(a, zero)
		y := add(zero, a)
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
		a := element(i)
		x := mul(a, one)
		y := mul(one, a)
		if x != a {
			t.Errorf("%02x: a * 1 != a: %02x", i, x)
		}
		if y != a {
			t.Errorf("%02x: 1 * a != a: %02x", i, y)
		}
	})
}

func FuzzAssociativeLaw(f *testing.F) {
	f.Add(uint8(0), uint8(1), uint8(2))
	f.Add(uint8(1), uint8(2), uint8(3))

	// test a + (b + c) = (a + b) + c
	f.Fuzz(func(t *testing.T, i, j, k uint8) {
		a := element(i)
		b := element(j)
		c := element(k)
		x := add(a, add(b, c))
		y := add(add(a, b), c)
		if x != y {
			t.Errorf("a + (b + c) != (a + b) + c: %02x, %02x, %02x", i, j, k)
		}
	})
}
