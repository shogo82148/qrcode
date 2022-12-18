package microqr

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"testing"

	bitmap "github.com/shogo82148/qrcode/bitmap"
)

func TestDecodeBitmap1(t *testing.T) {
	binimg := &bitmap.Image{
		Pix: []byte{
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
		},
		Stride: 2,
		Rect:   image.Rect(0, 0, 13, 13),
	}

	qr, err := DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("01234567")) {
		t.Errorf("want %q, got %q", []byte("01234567"), qr.Segments[0].Data)
	}
}

func TestDecodeBitmap2(t *testing.T) {
	// from https://www.qrcode.com/codes/microqr.html
	r, err := os.Open("testdata/01.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 15, 15))
	for y := 0; y <= 15; y++ {
		for x := 0; x <= 15; x++ {
			X := float64(x)*(55.0/15.0) + 2
			Y := float64(y)*(55.0/15.0) + 2
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}

	qr, err := DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("MICROQR")) {
		t.Errorf("want %q, got %q", []byte("MICROQR"), qr.Segments[0].Data)
	}
}

func TestDecodeBitmap3(t *testing.T) {
	// from https://www.qrcode.com/codes/microqr.html
	r, err := os.Open("testdata/02.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 13, 13))
	for y := 0; y < 13; y++ {
		for x := 0; x < 13; x++ {
			X := float64(x)*(21.0/13.0) + 5.5
			Y := float64(y)*(21.0/13.0) + 5.0
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}

	qr, err := DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("12345")) {
		t.Errorf("want %q, got %q", []byte("12345"), qr.Segments[0].Data)
	}
}

func TestDecodeBitmap4(t *testing.T) {
	// from https://www.qrcode.com/codes/microqr.html
	r, err := os.Open("testdata/03.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 15, 15))
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			X := float64(x)*(26.0/15.0) + 5.0
			Y := float64(y)*(27.0/15.0) + 5.0
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	log.Printf("%08b", binimg.Pix)

	qr, err := DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("1haicso")) {
		t.Errorf("want %q, got %q", []byte("1haicso"), qr.Segments[0].Data)
	}
}

func round(x float64) int {
	return int(math.Round(x))
}
