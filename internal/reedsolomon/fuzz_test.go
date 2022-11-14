package reedsolomon

import "testing"

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

func FuzzAssociativeLawAdd(f *testing.F) {
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

func FuzzAssociativeLawMul(f *testing.F) {
	f.Add(uint8(0), uint8(1), uint8(2))
	f.Add(uint8(1), uint8(2), uint8(3))

	// test a + (b + c) = (a + b) + c
	f.Fuzz(func(t *testing.T, i, j, k uint8) {
		a := element(i)
		b := element(j)
		c := element(k)
		x := mul(a, mul(b, c))
		y := mul(mul(a, b), c)
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
		a := element(i)
		b := element(j)
		c := element(k)
		x := mul(a, add(b, c))
		y := add(mul(a, b), mul(a, c))
		if x != y {
			t.Errorf("a(b + c) != ab + ac: %02x, %02x, %02x", i, j, k)
		}
	})
}

func FuzzSum(f *testing.F) {
	f.Add(10, []byte{
		0b0001_0000, 0b0010_0000, 0b0000_1100, 0b0101_0110,
		0b0110_0001, 0b1000_0000,

		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
	})

	f.Add(5, []byte{
		0b0100_0000, 0b0001_1000, 0b1010_1100, 0b1100_0011,
		0b0000_0000,
	})

	f.Add(17, []byte{
		0b0010_0000, 0b0100_0001, 0b1100_1101, 0b0100_0101,
		0b0010_1001, 0b1101_1100, 0b0010_1110, 0b1000_0000,
		0b1110_1100,
	})

	f.Fuzz(func(t *testing.T, n int, data []byte) {
		if n < 2 || n > 68 {
			return
		}
		w := New(n)
		if _, err := w.Write(data); err != nil {
			t.Fatal(err)
		}
		w.Sum(make([]byte, 0, n))
	})
}
