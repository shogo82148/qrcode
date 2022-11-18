package qrcode

import (
	"image/png"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	r, err := os.Open("testdata/wikipedia-version1.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	qr, err := Decode(img)
	if err != nil {
		t.Fatal(err)
	}
	_ = qr
}
