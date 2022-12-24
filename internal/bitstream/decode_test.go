package bitstream

import (
	"bytes"
	"testing"
)

func TestDecodeNumeric(t *testing.T) {
	tests := []struct {
		in   []byte
		want []byte
	}{
		{
			in:   []byte{0b0000_0000, 0b0000_0000},
			want: []byte{'0', '0', '0'},
		},
		{
			in:   []byte{0b0000_0000, 0b0100_0000},
			want: []byte{'0', '0', '1'},
		},
		{
			in:   []byte{0b1000_0000, 0b0000_0000},
			want: []byte{'5', '1', '2'},
		},
		{
			in:   []byte{999 >> 2, (999 & 0b11) << 6},
			want: []byte{'9', '9', '9'},
		},
		{
			in:   []byte{0b0000_0000},
			want: []byte{'0', '0'},
		},
		{
			in:   []byte{0b0000_0010},
			want: []byte{'0', '1'},
		},
		{
			in:   []byte{0b1000_0000},
			want: []byte{'6', '4'},
		},
		{
			in:   []byte{99 << 1},
			want: []byte{'9', '9'},
		},
		{
			in:   []byte{0b0000_0000},
			want: []byte{'0'},
		},
		{
			in:   []byte{0b0001_0000},
			want: []byte{'1'},
		},
		{
			in:   []byte{0b1000_0000},
			want: []byte{'8'},
		},
		{
			in:   []byte{0b1001_0000},
			want: []byte{'9'},
		},
	}

	for i, tt := range tests {
		buf := NewBuffer(tt.in)
		got := make([]byte, len(tt.want))
		if err := DecodeNumeric(buf, got); err != nil {
			t.Errorf("%d: error %v", i, err)
			continue
		}
		if !bytes.Equal(tt.want, got) {
			t.Errorf("%d: got %q, want %q", i, got, tt.want)
		}
	}
}

func TestDecodeAlphanumeric(t *testing.T) {
	tests := []struct {
		in   []byte
		want []byte
	}{
		{
			in:   []byte{0b0000_0000, 0b0000_0000},
			want: []byte{'0', '0'},
		},
		{
			in:   []byte{0b0000_0000, 0b0010_0000},
			want: []byte{'0', '1'},
		},
		{
			in:   []byte{0b1111_1101, 0b0000_0000},
			want: []byte{':', ':'},
		},
		{
			in:   []byte{0 << 2},
			want: []byte{'0'},
		},
		{
			in:   []byte{44 << 2},
			want: []byte{':'},
		},
	}

	for i, tt := range tests {
		buf := NewBuffer(tt.in)
		got := make([]byte, len(tt.want))
		if err := DecodeAlphanumeric(buf, got); err != nil {
			t.Errorf("%d: error %v", i, err)
			continue
		}
		if !bytes.Equal(tt.want, got) {
			t.Errorf("%d: got %q, want %q", i, got, tt.want)
		}
	}
}

func TestDecodeKanji(t *testing.T) {
	tests := []struct {
		in    []byte
		count int
		want  []byte
	}{
		{
			in:    []byte{0b01101100, 0b11111000},
			count: 1,
			want:  []byte("点"),
		},
		{
			in:    []byte{0b11010101, 0b01010000},
			count: 1,
			want:  []byte("茗"),
		},
	}

	for i, tt := range tests {
		buf := NewBuffer(tt.in)
		got, err := DecodeKanji(buf, tt.count)
		if err != nil {
			t.Errorf("%d: error %v", i, err)
			continue
		}
		if !bytes.Equal(tt.want, got) {
			t.Errorf("%d: got %q, want %q", i, got, tt.want)
		}
	}
}
