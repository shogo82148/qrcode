package microqr

import (
	"bytes"
	"testing"
)

func TestEncodeToBitmap1(t *testing.T) {
	qr := &QRCode{
		Version: 2,
		Level:   LevelL,
		Mask:    Mask1,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("01234567"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	// X 0510 : 2018
	// 附属書I (参考)
	// シンボルの符号化例
	want := []byte{
		0b11111110, 0b10101000,
		0b10000010, 0b11101000,
		0b10111010, 0b01101000,
		0b10111010, 0b01111000,
		0b10111010, 0b11100000,
		0b10000010, 0b10001000,
		0b11111110, 0b01111000,
		0b00000000, 0b01100000,
		0b11010000, 0b10001000,
		0b01101010, 0b10101000,
		0b11100111, 0b11110000,
		0b00010100, 0b00110000,
		0b11101001, 0b10111000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}
