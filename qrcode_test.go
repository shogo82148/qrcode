package qrcode

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func TestQRCode(t *testing.T) {
	var buf bytes.Buffer
	img := Generate()
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("qrcode.png", buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}
