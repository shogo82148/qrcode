package reedsolomon

import "testing"

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
