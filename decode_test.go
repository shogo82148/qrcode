package qrcode

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/shogo82148/qrcode/internal/binimage"
)

func TestDecode1(t *testing.T) {
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
			X := float64(x)*(143/21.0) + 40
			Y := float64(y)*(143/21.0) + 40
			binimg.SetBinary(x, y, imageAt(img, X, Y))
		}
	}

	if _, err := Decode(binimg); err != nil {
		t.Fatal(err)
	}
}

func TestDecode2(t *testing.T) {
	// from https://commons.wikimedia.org/wiki/File:Qr-2.png
	r, err := os.Open("testdata/version2.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := binimage.New(image.Rect(0, 0, 25, 25))
	for y := 0; y <= 25; y++ {
		for x := 0; x <= 25; x++ {
			X := float64(x)*(150/25) + 38
			Y := float64(y)*(150/25) + 38
			binimg.SetBinary(x, y, imageAt(img, X, Y))
		}
	}

	if _, err := Decode(binimg); err != nil {
		t.Fatal(err)
	}
}

func TestDecode3(t *testing.T) {
	// from https://commons.wikimedia.org/wiki/File:Qr-3.png
	r, err := os.Open("testdata/version3.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := binimage.New(image.Rect(0, 0, 29, 29))
	for y := 0; y <= 29; y++ {
		for x := 0; x <= 29; x++ {
			X := float64(x)*(150/29) + 39
			Y := float64(y)*(150/29) + 39
			if imageAt(img, X, Y) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
			binimg.SetBinary(x, y, imageAt(img, X, Y))
		}
		fmt.Println()
	}

	if _, err := Decode(binimg); err != nil {
		t.Fatal(err)
	}
}

func TestDecode4(t *testing.T) {
	// from https://commons.wikimedia.org/wiki/File:Qr-4.png
	r, err := os.Open("testdata/version4.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := binimage.New(image.Rect(0, 0, 33, 33))
	for y := 0; y < 33; y++ {
		for x := 0; x < 33; x++ {
			X := float64(x)*(170/33) + 29
			Y := float64(y)*(170/33) + 29
			binimg.SetBinary(x, y, imageAt(img, X, Y))
		}
	}

	if _, err := Decode(binimg); err != nil {
		t.Fatal(err)
	}
}

func TestDecode10(t *testing.T) {
	// from https://commons.wikimedia.org/wiki/File:Qr-code-ver-10.png
	r, err := os.Open("testdata/version10.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := binimage.New(image.Rect(0, 0, 57, 57))
	for y := 0; y < 57; y++ {
		for x := 0; x < 57; x++ {
			X := float64(x)*(180/57) + 25
			Y := float64(y)*(180/57) + 25
			binimg.SetBinary(x, y, imageAt(img, X, Y))
		}
	}

	if _, err := Decode(binimg); err != nil {
		t.Fatal(err)
	}
}
