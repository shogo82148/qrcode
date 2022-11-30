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
