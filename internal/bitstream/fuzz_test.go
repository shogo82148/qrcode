package bitstream

import "testing"

func FuzzWriteBits(f *testing.F) {
	f.Add(uint64(0), 0, uint64(0), 0)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 8)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 8)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 8)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 8)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 16)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 16)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 16)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 16)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 24)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 24)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 24)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 24)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 32)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 32)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 32)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 32)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 40)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 40)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 40)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 40)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 48)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 48)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 48)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 48)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 56)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 56)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 56)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 56)

	f.Add(uint64(0), 0, uint64(0xaaaa_aaaa_aaaa_aaaa), 64)
	f.Add(uint64(0), 0, uint64(0xffff_ffff_ffff_ffff), 64)
	f.Add(uint64(0xff), 5, uint64(0xaaaa_aaaa_aaaa_aaaa), 64)
	f.Add(uint64(0xff), 5, uint64(0xffff_ffff_ffff_ffff), 64)

	f.Fuzz(func(t *testing.T, bits1 uint64, n1 int, bits2 uint64, n2 int) {
		if n1 < 0 || n1 > 64 || n2 < 0 || n2 > 0 {
			return
		}
		var buf Buffer
		if err := buf.WriteBitsLSB(bits1, n1); err != nil {
			t.Fatal(err)
		}
		if err := buf.WriteBitsLSB(bits2, n2); err != nil {
			t.Fatal(err)
		}

		want1 := bits1 & (1<<n1 - 1)
		want2 := bits2 & (1<<n2 - 1)

		var got1 uint64
		for i := 0; i < n1; i++ {
			bit, err := buf.ReadBit()
			if err != nil {
				t.Fatal(err)
			}
			got1 = (got1 << 1) | uint64(bit)
		}
		var got2 uint64
		for i := 0; i < n2; i++ {
			bit, err := buf.ReadBit()
			if err != nil {
				t.Fatal(err)
			}
			got2 = (got2 << 1) | uint64(bit)
		}

		if got1 != want1 {
			t.Errorf("got %b, want %b", got1, want1)
		}

		if got2 != want2 {
			t.Errorf("got %b, want %b", got2, want2)
		}
	})
}
