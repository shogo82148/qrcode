package image

import (
	"image"
	"image/color"
)

var _ image.Image = (*Binary)(nil)

type Color bool

const White Color = true
const Black Color = false

func (c Color) RGBA() (r, g, b, a uint32) {
	if c {
		return 0xffff, 0xffff, 0xffff, 0xffff
	}
	return 0, 0, 0, 0xffff
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
	return color.RGBAModel
}

func (img *Binary) Bounds() image.Rectangle {
	return img.Rect
}

func (img *Binary) At(x, y int) color.Color {
	return img.BinaryAt(x, y)
}

func (img *Binary) BinaryAt(x, y int) Color {
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := (x - img.Rect.Min.X) % 8
	return Color(img.Pix[offset]>>shift&0x01 != 0)
}

func (img *Binary) SetBinary(x, y int, c Color) {
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := (x - img.Rect.Min.X) % 8
	mask := byte(0x80 >> shift)
	if c {
		img.Pix[offset] |= mask
	} else {
		img.Pix[offset] &^= mask
	}
}
