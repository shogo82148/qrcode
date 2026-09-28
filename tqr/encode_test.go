package tqr

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    []byte("123456789012"),
			wantErr: false,
		},
		{
			name:    "invalid length",
			data:    []byte("123"),
			wantErr: true,
		},
		{
			name:    "non-numeric data",
			data:    []byte("12345678901a"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncodeToBitmap(t *testing.T) {
	data := []byte("000000000000")
	qr, err := New(data)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Errorf("EncodeToBitmap() error = %v", err)
	}

	want := []byte{
		0b11111111, 0b11111111, 0b11111110,
		0b10000000, 0b00000000, 0b00000010,
		0b10111111, 0b10110011, 0b11111010,
		0b10100000, 0b10100010, 0b00001010,
		0b10101110, 0b10011010, 0b11101010,
		0b10101110, 0b10010010, 0b11101010,
		0b10101110, 0b10100010, 0b11101010,
		0b10100000, 0b10100010, 0b00001010,
		0b10111111, 0b10101011, 0b11111010,
		0b10000000, 0b00011000, 0b00000000,
		0b10101010, 0b10100001, 0b10011010,
		0b10111000, 0b01100011, 0b10001000,
		0b10000111, 0b10011110, 0b01110010,
		0b10000000, 0b00111100, 0b01111000,
		0b10111111, 0b10100001, 0b11001010,
		0b10100000, 0b10100011, 0b10001000,
		0b10101110, 0b10001100, 0b01110010,
		0b10101110, 0b10011100, 0b01010010,
		0b10101110, 0b10101010, 0b10001010,
		0b10100000, 0b10100011, 0b10001010,
		0b10111111, 0b10011100, 0b01111010,
		0b10000000, 0b00000000, 0b00000010,
		0b11111111, 0b10101010, 0b11111110,
	}
	if !bytes.Equal(img.Pix, want) {
		t.Errorf("EncodeToBitmap() = %v, want %v", img.Pix, want)
	}
}
