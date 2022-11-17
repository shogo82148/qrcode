package image

import (
	"image"
	"image/color"
	"math/bits"
)

var _ image.Image = (*Binary)(nil)

type Color bool

const White Color = false
const Black Color = true

func (c Color) RGBA() (r, g, b, a uint32) {
	if c {
		return 0, 0, 0, 0xffff
	}
	return 0xffff, 0xffff, 0xffff, 0xffff
}

// Binary is a binary image.
type Binary struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func New(r image.Rectangle) *Binary {
	stride := (r.Dx() + 7) / 8
	return &Binary{
		Pix:    make([]uint8, r.Dy()*stride),
		Stride: stride,
		Rect:   r,
	}
}

func (img *Binary) ColorModel() color.Model {
	return color.GrayModel
}

func (img *Binary) Bounds() image.Rectangle {
	return img.Rect
}

func (img *Binary) At(x, y int) color.Color {
	return img.BinaryAt(x, y)
}

func (img *Binary) BinaryAt(x, y int) Color {
	if !(image.Point{x, y}).In(img.Rect) {
		return White
	}
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := 7 - (x-img.Rect.Min.X)%8
	return Color((img.Pix[offset]>>shift)&0x01 != 0)
}

func (img *Binary) SetBinary(x, y int, c Color) {
	if !(image.Point{x, y}).In(img.Rect) {
		return
	}
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := (x - img.Rect.Min.X) % 8
	mask := byte(0x80 >> shift)
	if c {
		img.Pix[offset] |= mask
	} else {
		img.Pix[offset] &^= mask
	}
}

func (img *Binary) XorBinary(x, y int, c Color) {
	if !(image.Point{x, y}).In(img.Rect) {
		return
	}
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := (x - img.Rect.Min.X) % 8
	mask := byte(0x80 >> shift)
	if c {
		img.Pix[offset] ^= mask
	}
}

func (img *Binary) OnesCount() int64 {
	var cnt int64
	for _, b := range img.Pix {
		cnt += int64(bits.OnesCount8(b))
	}
	return cnt
}
