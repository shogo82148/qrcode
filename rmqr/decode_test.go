package rmqr

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/shogo82148/qrcode/bitmap"
)

func TestDecode1(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/rmqr2.png
	r, err := os.Open("testdata/rmqr2.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 59, 15))
	for y := 0; y <= 15; y++ {
		for x := 0; x <= 59; x++ {
			X := float64(x)*(164/59.0) + 9
			Y := float64(y)*(46/15.0) + 9
			binimg.Set(x, y, img.At(round(X), round(Y)))
			if binimg.BinaryAt(x, y) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func TestDecode2(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r7x43.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 43, 7))
	for y := 0; y < 7; y++ {
		for x := 0; x < 43; x++ {
			X := float64(x)*(335/43.0) + 29
			Y := float64(y)*(50/7.0) + 30
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}
