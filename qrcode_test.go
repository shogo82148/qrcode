package qrcode

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/shogo82148/qrcode/internal/bitstream"
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

func TestSegment_encodeNumber(t *testing.T) {
	s := &Segment{
		Mode: ModeNumber,
		Data: []byte("01234567"),
	}
	var buf bitstream.Buffer

	if err := s.encode(1, &buf); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0b0001_0000, 0b001000_00, 0b00001100, 0b01010110, 0b01_100001, 0b1_0000000,
	}
	got := buf.Bytes()

	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}
