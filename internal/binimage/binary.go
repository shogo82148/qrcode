package binimage

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

func (img *Binary) Mask(in, mask, pattern *Binary) *Binary {
	if !in.Rect.Eq(mask.Rect) {
		panic("binimage: in and mask must have same bounds")
	}

	img.Copy(in)
	dx, dy := img.Rect.Dx(), img.Rect.Dy()
	edge := byte(int(0xFF00 >> (dx % 8)))
	dx -= dx % 8
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x += 8 {
			offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
			offsetPattern := (y-pattern.Rect.Min.Y)*pattern.Stride + (x-pattern.Rect.Min.X)/8
			img.Pix[offset] ^= ^mask.Pix[offset] & pattern.Pix[offsetPattern]
		}
		if edge != 0 {
			x := dx
			offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
			offsetPattern := (y-pattern.Rect.Min.Y)*pattern.Stride + (x-pattern.Rect.Min.X)/8
			img.Pix[offset] ^= ^mask.Pix[offset] & pattern.Pix[offsetPattern] & edge
		}
	}
	return img
}

// Clone returns a clone of img.
func (img *Binary) Clone() *Binary {
	pix := append([]byte(nil), img.Pix...)
	return &Binary{
		Pix:    pix,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

func (img *Binary) Copy(from *Binary) *Binary {
	img.Pix = append(img.Pix[:0], from.Pix...)
	img.Stride = from.Stride
	img.Rect = from.Rect
	return img
}

// OnesCount returns the number of 1-pixels (black-pixels).
func (img *Binary) OnesCount() int {
	var cnt int
	dx := img.Rect.Dx()
	length := dx / 8
	dy := img.Rect.Dy()
	mask := byte(int(0xff00) >> (dx % 8))
	for y := 0; y < dy; y++ {
		for x := 0; x < length; x++ {
			offset := y*img.Stride + x
			cnt += bits.OnesCount8(img.Pix[offset])
		}
		if mask != 0 {
			cnt += bits.OnesCount8(img.Pix[y*img.Stride+length] & mask)
		}
	}
	return cnt
}

func (img *Binary) Point() int {
	return img.finderPattern() + img.longRunLengthCount() + img.blockCount() + img.pointOnesCount()
}

func (img *Binary) longRunLengthCount() int {
	var cnt int
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		var length int
		c0 := img.BinaryAt(img.Rect.Min.X, y)
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.BinaryAt(x, y)
			if c == c0 {
				length++
			} else {
				if length >= 5 {
					cnt += length - 5 + 3
				}
				c0 = c
				length = 0
			}
		}
	}

	for x := img.Rect.Min.Y; x < img.Rect.Max.Y; x++ {
		var length int
		c0 := img.BinaryAt(x, img.Rect.Min.X)
		for y := img.Rect.Min.X; y < img.Rect.Max.X; y++ {
			c := img.BinaryAt(x, y)
			if c == c0 {
				length++
			} else {
				if length >= 5 {
					cnt += length - 5 + 3
				}
				c0 = c
				length = 0
			}
		}
	}

	return cnt
}

func (img *Binary) blockCount() int {
	var cnt int
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y-1; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X-1; x++ {
			c1 := img.BinaryAt(y, x)
			c2 := img.BinaryAt(y, x+1)
			c3 := img.BinaryAt(y+1, x)
			c4 := img.BinaryAt(y+1, x+1)
			if c1 == c2 && c1 == c3 && c1 == c4 {
				cnt++
			}
		}
	}
	return cnt * 3
}

func (img *Binary) finderPattern() int {
	var cnt int
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			var c1, c2, c3, c4, c5, c6, c7 Color
			c1 = img.BinaryAt(x, y-3)
			c2 = img.BinaryAt(x, y-2)
			c3 = img.BinaryAt(x, y-1)
			c4 = img.BinaryAt(x, y)
			c5 = img.BinaryAt(x, y+1)
			c6 = img.BinaryAt(x, y+2)
			c7 = img.BinaryAt(x, y+3)
			if c1 && !c2 && c3 && c4 && c5 && !c6 && c7 {
				c := !img.BinaryAt(x, y-4) && !img.BinaryAt(x, y-5) && !img.BinaryAt(x, y-6) && !img.BinaryAt(x, y-7)
				c = c || !img.BinaryAt(x, y+4) && !img.BinaryAt(x, y+5) && !img.BinaryAt(x, y+6) && !img.BinaryAt(x, y+7)
				if c {
					cnt++
				}
			}

			c1 = img.BinaryAt(x-3, y)
			c2 = img.BinaryAt(x-2, y)
			c3 = img.BinaryAt(x-1, y)
			c4 = img.BinaryAt(x, y)
			c5 = img.BinaryAt(x-1, y)
			c6 = img.BinaryAt(x-2, y)
			c7 = img.BinaryAt(x-3, y)
			if c1 && !c2 && c3 && c4 && c5 && !c6 && c7 {
				c := !img.BinaryAt(x-4, y) && !img.BinaryAt(x-5, y) && !img.BinaryAt(x-6, y) && !img.BinaryAt(x-7, y-7)
				c = c || !img.BinaryAt(x+4, y) && !img.BinaryAt(x-5, y) && !img.BinaryAt(x+6, y) && !img.BinaryAt(x, y+7)
				if c {
					cnt++
				}
			}
		}
	}
	return cnt * 40
}

func (img *Binary) pointOnesCount() int {
	total := img.Rect.Dx() * img.Rect.Dy()
	cnt := img.OnesCount()
	p := float64(cnt)/float64(total) - 0.5
	if p < 0 {
		p = -p
	}
	return int(p*20) * 10
}
