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
		bits byte
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
		{
			in: []seq{
				{0b0000, 4},
				{0b1000_1110, 8},
				{0b1111, 4},
			},
			want: []byte{0b0000_1000, 0b1110_1111},
		},
	}

	for _, tt := range tests {
		var buf Buffer
		for _, in := range tt.in {
			if err := buf.WriteBitsLSB(in.bits, in.n); err != nil {
				t.Fatal(err)
			}
		}
		got := buf.Bytes()
		if !bytes.Equal(got, tt.want) {
			t.Errorf("got %x, want %x", got, tt.want)
		}
	}
}
