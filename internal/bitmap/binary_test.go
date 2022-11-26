package bitmap

import (
	"bytes"
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

func TestMask(t *testing.T) {
	tests := []struct {
		w, h    int
		in      []byte
		mask    []byte
		pattern []byte
		want    []byte
	}{
		{
			w: 8,
			h: 4,
			in: []byte{
				0b0000_0000,
				0b0000_0000,
				0b0000_0000,
				0b0000_0000,
			},
			mask: []byte{
				0b1111_0000,
				0b1111_0000,
				0b1111_0000,
				0b1111_0000,
			},
			pattern: []byte{
				0b1010_1010,
				0b0101_0101,
				0b1010_1010,
				0b0101_0101,
			},
			want: []byte{
				0b0000_1010,
				0b0000_0101,
				0b0000_1010,
				0b0000_0101,
			},
		},
	}

	for i, tt := range tests {
		in := &Image{
			Pix:    tt.in,
			Stride: (tt.w + 7) / 8,
			Rect:   image.Rect(0, 0, tt.w, tt.h),
		}
		mask := &Image{
			Pix:    tt.mask,
			Stride: (tt.w + 7) / 8,
			Rect:   image.Rect(0, 0, tt.w, tt.h),
		}
		pattern := &Image{
			Pix:    tt.pattern,
			Stride: (tt.w + 7) / 8,
			Rect:   image.Rect(0, 0, tt.w, tt.h),
		}
		var img Image
		img.Mask(in, mask, pattern)
		if !bytes.Equal(img.Pix, tt.want) {
			t.Errorf("%d: got %08b, want %08b", i, img.Pix, tt.want)
		}
	}
}
