package microqr

import (
	"bytes"
	"testing"
)

func TestEncode1(t *testing.T) {
	qr, err := New(LevelL, []byte("MICROQR"))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 3 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 3)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeAlphanumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeAlphanumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("MICROQR")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "MICROQR")
	}
}

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

func TestEncodeToBitmap2(t *testing.T) {
	qr := &QRCode{
		Version: 3,
		Level:   LevelM,
		Mask:    Mask3,
		Segments: []Segment{
			{
				Mode: ModeAlphanumeric,
				Data: []byte("MICROQR"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010, 0b10000010,
		0b11011110, 0b10111010, 0b10111100,
		0b10111010, 0b10011010, 0b10111010,
		0b00000110, 0b10000010, 0b00010100,
		0b11111110, 0b01001010, 0b00000000,
		0b00101000, 0b10001001, 0b10101010,
		0b01000001, 0b10011010, 0b11111111,
		0b10001000, 0b00010101, 0b11000010,
		0b10100110, 0b11001100, 0b01111100,
		0b11010100, 0b11000100, 0b11011100,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}
