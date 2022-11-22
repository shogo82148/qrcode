package qrcode

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/shogo82148/qrcode/internal/binimage"
)

func TestDecode(t *testing.T) {
	// from https://commons.wikimedia.org/wiki/File:Qr-1.png
	r, err := os.Open("testdata/version1.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := binimage.New(image.Rect(0, 0, 21, 21))
	for y := 0; y <= 21; y++ {
		for x := 0; x <= 21; x++ {
			X := float64(x)*(143.0/21.0) + 40
			Y := float64(y)*(143.0/21.0) + 40
			binimg.SetBinary(x, y, imageAt(img, X, Y))
		}
	}

	if _, err := Decode(binimg); err != nil {
		t.Fatal(err)
	}
}
