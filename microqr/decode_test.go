package microqr

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/shogo82148/go-imaging/bitmap"
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
	// from https://www.qrcode.com/img/rmqr/gra2.jpg
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
	// from https://www.qrcode.com/img/rmqr/gra2.jpg
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

	qr, err := DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("1haicso")) {
		t.Errorf("want %q, got %q", []byte("1haicso"), qr.Segments[0].Data)
	}
}

func TestDecodeBitmap5(t *testing.T) {
	r, err := os.Open("testdata/04.png")
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
			X := float64(x)*(89.0/14.0) + 15
			Y := float64(y)*(81.0/14.0) + 18
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}

	qr, err := DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
	if qr.Version != 3 {
		t.Errorf("unexpected version: got %d, want %d", qr.Version, 3)
	}
	if qr.Level != LevelL {
		t.Errorf("unexpected level: got %d, want %d", qr.Level, LevelL)
	}
	if qr.Mask != Mask1 {
		t.Errorf("unexpected mask: got %d, want %d", qr.Mask, Mask1)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("AINIX")) {
		t.Errorf("want %q, got %q", []byte("AINIX"), qr.Segments[0].Data)
	}
	if !bytes.Equal(qr.Segments[1].Data, []byte("12345")) {
		t.Errorf("want %q, got %q", []byte("12345"), qr.Segments[1].Data)
	}
}

func TestDecodeBitmap6(t *testing.T) {
	// from https://docs.esko.com/docs/ja-jp/dynamicbarcodes-for-ai/12/userguide/assets/bar/ss_bar_Micro_QR.png
	// https://docs.esko.com/docs/ja-jp/dynamicbarcodes-for-ai/12/userguide/ja-jp/common/bar/reference/re_bar_MicroQR.html
	r, err := os.Open("testdata/05.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 11, 11))
	for y := 0; y < 11; y++ {
		for x := 0; x < 11; x++ {
			X := float64(x)*(28.0/10.0) + 5
			Y := float64(y)*(28.0/10.0) + 6
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}

	qr, err := DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
	if qr.Version != 1 {
		t.Errorf("unexpected version: got %d, want %d", qr.Version, 3)
	}
	if qr.Level != LevelCheck {
		t.Errorf("unexpected level: got %d, want %d", qr.Level, LevelCheck)
	}
	if qr.Mask != Mask2 {
		t.Errorf("unexpected mask: got %d, want %d", qr.Mask, Mask2)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("0000")) {
		t.Errorf("want %q, got %q", []byte("0000"), qr.Segments[0].Data)
	}
}

func round(x float64) int {
	return int(math.Round(x))
}

// copyBitmap returns a bitmap with the bounds rect that has the same pixels as src.
func copyBitmap(src *bitmap.Image, rect image.Rectangle) *bitmap.Image {
	img := bitmap.New(rect)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.SetBinary(x, y, src.BinaryAt(x, y))
		}
	}
	return img
}

// padStride returns a copy of src whose Stride is larger than necessary.
func padStride(src *bitmap.Image) *bitmap.Image {
	stride := src.Stride + 3
	pix := make([]byte, stride*src.Rect.Dy())
	for y := 0; y < src.Rect.Dy(); y++ {
		copy(pix[y*stride:], src.Pix[y*src.Stride:(y+1)*src.Stride])
	}
	return &bitmap.Image{Pix: pix, Stride: stride, Rect: src.Rect}
}

func TestDecodeInvalidSize(t *testing.T) {
	// a valid M4 symbol, used to put correct format information into malformed bitmaps.
	qr, err := New([]byte("12345"), WithLevel(LevelL), WithKanji(false))
	if err != nil {
		t.Fatal(err)
	}
	qr.Version = 4
	src, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}

	tests := []image.Rectangle{
		image.Rect(0, 0, 0, 0),
		image.Rect(0, 0, 11, 11), // M1 size
		image.Rect(0, 0, 13, 13), // M2 size
		image.Rect(0, 0, 30, 30),
		image.Rect(0, 0, 17, 18), // not square
		image.Rect(1, 1, 18, 18), // non-zero origin
	}
	for _, rect := range tests {
		img := copyBitmap(src, rect)
		if _, err := DecodeBitmap(img); err == nil {
			t.Errorf("%v: want error, got nil", rect)
		}
	}
}

func TestDecodeBitmapStride(t *testing.T) {
	qr, err := New([]byte("12345"), WithLevel(LevelL), WithKanji(false))
	if err != nil {
		t.Fatal(err)
	}
	src, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	img := padStride(src)
	orig := bytes.Clone(img.Pix)

	got, err := DecodeBitmap(img)
	if err != nil {
		t.Fatal(err)
	}
	if string(got.Segments[0].Data) != "12345" {
		t.Errorf("unexpected data: got %q, want %q", got.Segments[0].Data, "12345")
	}
	if !bytes.Equal(img.Pix, orig) {
		t.Error("DecodeBitmap modified the input image")
	}

	// the pixel buffer is too small for the bounds.
	img.Pix = img.Pix[:len(img.Pix)-img.Stride]
	if _, err := DecodeBitmap(img); err == nil {
		t.Error("want error, got nil")
	}
}
