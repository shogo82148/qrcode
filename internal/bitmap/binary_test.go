package bitmap

import (
	"bytes"
	"image"
	"math"
	"testing"

	"github.com/shogo82148/go-imaging/bitmap"
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

func TestImport(t *testing.T) {
	src := bitmap.New(image.Rect(0, 0, 11, 11))
	src.SetBinary(10, 10, Black)
	img, err := Import(src)
	if err != nil {
		t.Fatal(err)
	}
	if img.BinaryAt(10, 10) != Black {
		t.Errorf("got %v, want %v", img.BinaryAt(10, 10), Black)
	}
	img.SetBinary(0, 0, Black)
	if src.BinaryAt(0, 0) != White {
		t.Error("Import must copy the pixel buffer")
	}
}

func TestImportInvalid(t *testing.T) {
	tests := []struct {
		name string
		img  *bitmap.Image
	}{
		{
			name: "empty",
			img:  &bitmap.Image{Rect: image.Rect(0, 0, 0, 0)},
		},
		{
			name: "overflowing bounds",
			img:  &bitmap.Image{Rect: image.Rect(math.MinInt, 0, math.MaxInt, 1), Stride: 1, Pix: []byte{0}},
		},
		{
			// it must be rejected without allocating 512 MiB.
			name: "large bounds with empty buffer",
			img:  &bitmap.Image{Rect: image.Rect(0, 0, 65536, 65536), Stride: 8192},
		},
		{
			name: "small stride",
			img:  &bitmap.Image{Rect: image.Rect(0, 0, 11, 11), Stride: 1, Pix: make([]byte, 22)},
		},
		{
			name: "short buffer",
			img:  &bitmap.Image{Rect: image.Rect(0, 0, 11, 11), Stride: 2, Pix: make([]byte, 21)},
		},
		{
			// (Dy-1)*Stride overflows int.
			name: "huge stride",
			img:  &bitmap.Image{Rect: image.Rect(0, 0, 11, 11), Stride: 1 << 62, Pix: make([]byte, 22)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Import(tt.img); err == nil {
				t.Error("want error, got nil")
			}
		})
	}

	allocs := testing.AllocsPerRun(1, func() {
		Import(&bitmap.Image{Rect: image.Rect(0, 0, 65536, 65536), Stride: 8192})
	})
	if allocs > 1 {
		t.Errorf("Import allocated %v times for an invalid image", allocs)
	}
}
