package binimage

import (
	"image"
	"testing"
)

func TestSetBinary(t *testing.T) {
	img := New(image.Rect(0, 0, 8, 8))
	img.SetBinary(0, 0, Black)

	if img.Pix[0] != 0x80 {
		t.Errorf("got %02x, want %02x", img.Pix[0], 0x80)
	}
	if img.BinaryAt(0, 0) != Black {
		t.Errorf("got %v, want %v", img.BinaryAt(0, 0), Black)
	}
}
