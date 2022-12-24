package bitstream

import (
	"bytes"
	"testing"
	"unicode/utf8"
)

func TestEncodeNumeric(t *testing.T) {
	tests := []struct {
		in   []byte
		want []byte
	}{
		{
			in:   []byte{'0', '0', '0'},
			want: []byte{0b0000_0000, 0b0000_0000},
		},
		{
			in:   []byte{'0', '0', '1'},
			want: []byte{0b0000_0000, 0b0100_0000},
		},
		{
			in:   []byte{'5', '1', '2'},
			want: []byte{0b1000_0000, 0b0000_0000},
		},
		{
			in:   []byte{'9', '9', '9'},
			want: []byte{999 >> 2, (999 & 0b11) << 6},
		},
		{
			in:   []byte{'0', '0'},
			want: []byte{0b0000_0000},
		},
		{
			in:   []byte{'0', '1'},
			want: []byte{0b0000_0010},
		},
		{
			in:   []byte{'6', '4'},
			want: []byte{0b1000_0000},
		},
		{
			in:   []byte{'9', '9'},
			want: []byte{99 << 1},
		},
		{
			in:   []byte{'0'},
			want: []byte{0b0000_0000},
		},
		{
			in:   []byte{'1'},
			want: []byte{0b0001_0000},
		},
		{
			in:   []byte{'8'},
			want: []byte{0b1000_0000},
		},
		{
			in:   []byte{'9'},
			want: []byte{0b1001_0000},
		},
	}

	for i, tt := range tests {
		var buf Buffer
		if err := EncodeNumeric(&buf, tt.in); err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		got := buf.Bytes()
		if !bytes.Equal(tt.want, got) {
			t.Errorf("%d: got %08b, want %08b", i, got, tt.want)
		}
	}
}

func TestEncodeAlphanumeric(t *testing.T) {
	tests := []struct {
		in   []byte
		want []byte
	}{
		{
			in:   []byte{'0', '0'},
			want: []byte{0b0000_0000, 0b0000_0000},
		},
		{
			in:   []byte{'0', '1'},
			want: []byte{0b0000_0000, 0b0010_0000},
		},
		{
			in:   []byte{':', ':'},
			want: []byte{0b1111_1101, 0b0000_0000},
		},
		{
			in:   []byte{'0'},
			want: []byte{0 << 2},
		},
		{
			in:   []byte{':'},
			want: []byte{44 << 2},
		},
	}

	for i, tt := range tests {
		var buf Buffer
		if err := EncodeAlphanumeric(&buf, tt.in); err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		got := buf.Bytes()
		if !bytes.Equal(tt.want, got) {
			t.Errorf("%d: got %08b, want %08b", i, got, tt.want)
		}
	}
}

func TestEncodeKanji(t *testing.T) {
	tests := []struct {
		in   []byte
		want []byte
	}{
		{
			in:   []byte("点"),
			want: []byte{0b01101100, 0b11111000},
		},
		{
			in:   []byte("茗"),
			want: []byte{0b11010101, 0b01010000},
		},
	}

	for i, tt := range tests {
		var buf Buffer
		if err := EncodeKanji(&buf, tt.in); err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		got := buf.Bytes()
		if !bytes.Equal(tt.want, got) {
			t.Errorf("%d: got %08b, want %08b", i, got, tt.want)
		}
	}
}

func TestEncodeKanji_Error(t *testing.T) {
	tests := []struct {
		in   []byte
		want []byte
	}{
		{
			in: []byte("😎"),
		},
		{
			in: []byte(string(utf8.RuneError)),
		},
	}

	for i, tt := range tests {
		var buf Buffer
		if err := EncodeKanji(&buf, tt.in); err == nil {
			t.Errorf("%d: want error, but not", i)
			continue
		}
	}
}
