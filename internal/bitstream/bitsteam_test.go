package bitstream

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestReadBit(t *testing.T) {
	buf := NewBuffer([]byte{
		0b1000_0000, 0b1000_0001,
	})

	got := make([]byte, 0, 16)
	for {
		bit, err := buf.ReadBit()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, bit)
	}

	want := []byte{
		1, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 1,
	}

	if !bytes.Equal(got, want) {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestWriteBit(t *testing.T) {
	in := []byte{
		1, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 1,
	}
	var buf Buffer
	for _, b := range in {
		if err := buf.WriteBit(b); err != nil {
			t.Fatal(err)
		}
	}
	got := buf.Bytes()
	want := []byte{
		0b1000_0000, 0b1000_0001,
	}

	if !bytes.Equal(got, want) {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestWriteBits(t *testing.T) {
	type seq struct {
		bits uint64
		n    int
	}
	tests := []struct {
		in   []seq
		want []byte
	}{
		{
			in: []seq{
				{0b1, 1},
				{0b000, 3},
				{0b1111, 4},
			},
			want: []byte{0b1000_1111},
		},
		{
			in: []seq{
				{0b111, 3},
			},
			want: []byte{0b1110_0000},
		},

		// Cases that cross byte boundaries
		// 1 byte
		{
			in: []seq{
				{0b0000, 4},
				{0b1000_1110, 8},
				{0b1111, 4},
			},
			want: []byte{0b0000_1000, 0b1110_1111},
		},
		{
			in: []seq{
				{0b000, 3},
				{0b1111_1111, 8},
				{0b00000, 5},
			},
			want: []byte{0b0001_1111, 0b1110_0000},
		},
		// 2 bytes
		{
			in: []seq{
				{0b0000, 4},
				{0b1000_0000_0000_1110, 16},
				{0b1111, 4},
			},
			want: []byte{0b0000_1000, 0b00000000, 0b1110_1111},
		},
		{
			in: []seq{
				{0b000, 3},
				{0b1111_1111_1111_1111, 16},
				{0b00000, 5},
			},
			want: []byte{0b000_11111, 0b11111111, 0b111_00000},
		},
		// 3 bytes
		{
			in: []seq{
				{0b0000, 4},
				{0b1000_0000_1111_1111_0000_1110, 24},
				{0b1111, 4},
			},
			want: []byte{0b0000_1000, 0b00001111, 0b11110000, 0b1110_1111},
		},
		{
			in: []seq{
				{0b000, 3},
				{0b1111_1111_1111_1111_1111_1111, 24},
				{0b00000, 5},
			},
			want: []byte{0b000_11111, 0b11111111, 0b11111111, 0b111_00000},
		},
		// 4 bytes
		{
			in: []seq{
				{0b0000, 4},
				{0b1000_0000_1111_0000_0000_1111_0000_1110, 32},
				{0b1111, 4},
			},
			want: []byte{0b0000_1000, 0b00001111, 0b00000000, 0b11110000, 0b1110_1111},
		},
		{
			in: []seq{
				{0b000, 3},
				{0b1111_1111_1111_0000_0000_1111_1111_1111, 32},
				{0b00000, 5},
			},
			want: []byte{0b000_11111, 0b11111110, 0b00000001, 0b11111111, 0b111_00000},
		},
	}

	for i, tt := range tests {
		var buf Buffer
		for _, in := range tt.in {
			if err := buf.WriteBitsLSB(in.bits, in.n); err != nil {
				t.Fatal(err)
			}
		}
		got := buf.Bytes()
		if !bytes.Equal(got, tt.want) {
			t.Errorf("%d: got %b, want %b", i, got, tt.want)
		}
	}
}

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
