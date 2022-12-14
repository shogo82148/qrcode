package rmqr

import (
	"bytes"
	"testing"
)

func TestEncodeToBitmap1(t *testing.T) {
	qr := &QRCode{
		Version: 0,
		Level:   LevelM,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("123456789012"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}
