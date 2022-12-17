package rmqr

import (
	"bytes"
	"testing"
)

func TestEncodeToBitmap1(t *testing.T) {
	qr := &QRCode{
		Version: 0,
		Level:   LevelM,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("123456789012"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010, 0b10101110, 0b10101010, 0b10101010, 0b11100000,
		0b10000010, 0b01010111, 0b11101010, 0b00001000, 0b11011000, 0b10100000,
		0b10111010, 0b10111000, 0b10011110, 0b11010110, 0b10111111, 0b11100000,
		0b10111010, 0b01100110, 0b11100000, 0b00011111, 0b11000010, 0b00100000,
		0b10111010, 0b00101100, 0b01111110, 0b01101111, 0b10110010, 0b10100000,
		0b10000010, 0b11110111, 0b11111010, 0b11010010, 0b10111010, 0b00100000,
		0b11111110, 0b10101010, 0b10101110, 0b10101010, 0b10101011, 0b11100000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}
