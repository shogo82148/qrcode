package bitmap

import (
	"image"
	"image/color"
	"image/draw"
)

// Color is binary color.
type Color bool

var _ color.Color = White

// White is white color.
const White Color = false

// Black is black color.
const Black Color = true

// RGBA implements [image/color.Color].
func (c Color) RGBA() (r, g, b, a uint32) {
	if c {
		return 0, 0, 0, 0xffff
	}
	return 0xffff, 0xffff, 0xffff, 0xffff
}

var ColorModel color.Model = color.ModelFunc(binaryModel)

func binaryModel(c color.Color) color.Color {
	bin, ok := c.(Color)
	if ok {
		return bin
	}

	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	return Color(y < 0x8000)
}

var _ image.Image = (*Image)(nil)
var _ draw.Image = (*Image)(nil)

// Image is a binary image.
type Image struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func New(r image.Rectangle) *Image {
	stride := (r.Dx() + 7) / 8
	return &Image{
		Pix:    make([]uint8, r.Dy()*stride),
		Stride: stride,
		Rect:   r,
	}
}

func (img *Image) ColorModel() color.Model {
	return ColorModel
}

func (img *Image) Bounds() image.Rectangle {
	return img.Rect
}

func (img *Image) At(x, y int) color.Color {
	return img.BinaryAt(x, y)
}

func (img *Image) BinaryAt(x, y int) Color {
	if !(image.Point{x, y}).In(img.Rect) {
		return White
	}
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := 7 - (x-img.Rect.Min.X)%8
	return Color((img.Pix[offset]>>shift)&0x01 != 0)
}

func (img *Image) Set(x, y int, c color.Color) {
	img.SetBinary(x, y, binaryModel(c).(Color))
}

func (img *Image) SetBinary(x, y int, c Color) {
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
