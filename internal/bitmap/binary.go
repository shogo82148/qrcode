package bitmap

import (
	"fmt"
	"image"
	"io"
	"math/bits"

	"github.com/shogo82148/go-imaging/bitmap"
)

type Color = bitmap.Color

const White Color = false
const Black Color = true

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

func Import(img *bitmap.Image) *Image {
	return &Image{
		Pix:    img.Pix,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

func (img *Image) Export() *bitmap.Image {
	return &bitmap.Image{
		Pix:    img.Pix,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

func (img *Image) BinaryAt(x, y int) Color {
	if !(image.Point{x, y}).In(img.Rect) {
		return White
	}
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := 7 - (x-img.Rect.Min.X)%8
	return Color((img.Pix[offset]>>shift)&0x01 != 0)
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

func (img *Image) XorBinary(x, y int, c Color) {
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

func (img *Image) Mask(in, mask, pattern *Image) *Image {
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
func (img *Image) Clone() *Image {
	pix := append([]byte(nil), img.Pix...)
	return &Image{
		Pix:    pix,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

func (img *Image) Copy(from *Image) *Image {
	img.Pix = append(img.Pix[:0], from.Pix...)
	img.Stride = from.Stride
	img.Rect = from.Rect
	return img
}

// OnesCount returns the number of 1-pixels (black-pixels).
func (img *Image) OnesCount() int {
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

func (img *Image) EncodePBM(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "P1"); err != nil {
		return err
	}
	dx := img.Rect.Dx()
	dy := img.Rect.Dy()

	if _, err := fmt.Fprintf(w, "%d %d\n", dx, dy); err != nil {
		return err
	}
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			if x != 0 {
				if _, err := fmt.Fprint(w, " "); err != nil {
					return err
				}
			}
			v := 0
			if img.BinaryAt(x+img.Rect.Min.X, y+img.Rect.Min.Y) {
				v = 1
			}
			if _, err := fmt.Fprintf(w, "%d", v); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

func (img *Image) Point() int {
	g := newGrid(img)
	return g.finderPattern() + g.longRunLengthCount() + g.blockCount() + img.pointOnesCount()
}

// grid is an unpacked copy of Image for fast random access.
// It is used for calculating the penalty score.
type grid struct {
	pix  []uint8 // 1 is black, 0 is white
	rect image.Rectangle
	w, h int
}

func newGrid(img *Image) *grid {
	w, h := img.Rect.Dx(), img.Rect.Dy()
	pix := make([]uint8, w*h)
	for y := 0; y < h; y++ {
		row := img.Pix[y*img.Stride:]
		line := pix[y*w : (y+1)*w]
		for x := range line {
			line[x] = (row[x/8] >> (7 - x%8)) & 0x01
		}
	}
	return &grid{
		pix:  pix,
		rect: img.Rect,
		w:    w,
		h:    h,
	}
}

// at is same as BinaryAt, but faster.
func (g *grid) at(x, y int) Color {
	x -= g.rect.Min.X
	y -= g.rect.Min.Y
	if uint(x) >= uint(g.w) || uint(y) >= uint(g.h) {
		return White
	}
	return g.pix[y*g.w+x] != 0
}

func (img *Image) PointMicro() int {
	var sum1, sum2 int
	for x := img.Rect.Min.X + 1; x < img.Rect.Max.X; x++ {
		if img.BinaryAt(x, img.Rect.Max.Y-1) {
			sum1++
		}
	}
	for y := img.Rect.Min.Y + 1; y < img.Rect.Max.Y; y++ {
		if img.BinaryAt(img.Rect.Max.X-1, y) {
			sum2++
		}
	}
	if sum1 > sum2 {
		sum1, sum2 = sum2, sum1
	}
	return sum1*16 + sum2
}

// isSquare reports whether the transposed coordinates used by the penalty rules
// point to the same pixels as the normal coordinates.
func (g *grid) isSquare() bool {
	return g.w == g.h && g.rect.Min.X == g.rect.Min.Y
}

func (g *grid) longRunLengthCount() int {
	if g.isSquare() {
		return g.longRunLengthCountSquare()
	}

	var cnt int
	for y := g.rect.Min.Y; y < g.rect.Max.Y; y++ {
		var length int
		c0 := g.at(g.rect.Min.X, y)
		for x := g.rect.Min.X; x < g.rect.Max.X; x++ {
			c := g.at(x, y)
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

	for x := g.rect.Min.Y; x < g.rect.Max.Y; x++ {
		var length int
		c0 := g.at(x, g.rect.Min.X)
		for y := g.rect.Min.X; y < g.rect.Max.X; y++ {
			c := g.at(x, y)
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

func (g *grid) longRunLengthCountSquare() int {
	var cnt int
	w := g.w
	if w == 0 {
		return 0
	}

	// horizontal
	for y := 0; y < w; y++ {
		line := g.pix[y*w : (y+1)*w]
		var length int
		c0 := line[0]
		for _, c := range line {
			length, cnt = runLength(length, cnt, c^c0)
			c0 = c
		}
	}

	// vertical
	for x := 0; x < w; x++ {
		var length int
		c0 := g.pix[x]
		for i := x; i < len(g.pix); i += w {
			c := g.pix[i]
			length, cnt = runLength(length, cnt, c^c0)
			c0 = c
		}
	}

	return cnt
}

// runLength is a branch-less version of the following code:
//
//	if diff == 0 {
//		length++
//	} else {
//		if length >= 5 {
//			cnt += length - 5 + 3
//		}
//		length = 0
//	}
func runLength(length, cnt int, diff uint8) (int, int) {
	d := int(diff)                                    // 1 if the color changes, otherwise 0
	ge5 := int(uint(4-length) >> (bits.UintSize - 1)) // 1 if length >= 5, otherwise 0
	cnt += (length - 2) & -(d & ge5)
	length = (length + 1) & (d - 1)
	return length, cnt
}

func (g *grid) blockCount() int {
	if g.isSquare() {
		return g.blockCountSquare()
	}

	var cnt int
	for y := g.rect.Min.Y; y < g.rect.Max.Y-1; y++ {
		for x := g.rect.Min.X; x < g.rect.Max.X-1; x++ {
			c1 := g.at(y, x)
			c2 := g.at(y, x+1)
			c3 := g.at(y+1, x)
			c4 := g.at(y+1, x+1)
			if c1 == c2 && c1 == c3 && c1 == c4 {
				cnt++
			}
		}
	}
	return cnt * 3
}

func (g *grid) blockCountSquare() int {
	var cnt int
	w := g.w
	for y := 0; y < w-1; y++ {
		line0 := g.pix[y*w : (y+1)*w]
		line1 := g.pix[(y+1)*w : (y+2)*w]
		for x := 0; x < w-1; x++ {
			// all pixels have the same color if the sum is 0 or 4.
			s := line0[x] + line0[x+1] + line1[x] + line1[x+1]
			cnt += int(((s&3)+3)>>2) ^ 1
		}
	}
	return cnt * 3
}

func (g *grid) finderPattern() int {
	var cnt int

	// vertical: 1:1:3:1:1 pattern
	for x := g.rect.Min.X; x < g.rect.Max.X; x++ {
		// window holds the pixels from (x, y-3) to (x, y+3).
		// the pixel (x, y-3) is the most significant bit.
		var window uint
		for yy := g.rect.Min.Y - 3; yy < g.rect.Max.Y+3; yy++ {
			window = (window<<1 | b2u(g.at(x, yy))) & 0x7f
			y := yy - 3
			if y < g.rect.Min.Y || window != 0b1011101 {
				continue
			}
			c := !g.at(x, y-4) && !g.at(x, y-5) && !g.at(x, y-6) && !g.at(x, y-7)
			c = c || !g.at(x, y+4) && !g.at(x, y+5) && !g.at(x, y+6) && !g.at(x, y+7)
			if c {
				cnt++
			}
		}
	}

	// horizontal
	// NOTE: it checks only the pixels from (x-3, y) to (x, y),
	// because the pixels (x-1, y), (x-2, y) and (x-3, y) are checked twice.
	for y := g.rect.Min.Y; y < g.rect.Max.Y; y++ {
		// window holds the pixels from (x-3, y) to (x, y).
		// the pixel (x-3, y) is the most significant bit.
		var window uint
		for x := g.rect.Min.X - 3; x < g.rect.Max.X; x++ {
			window = (window<<1 | b2u(g.at(x, y))) & 0xf
			if x < g.rect.Min.X || window != 0b1011 {
				continue
			}
			c := !g.at(x-4, y) && !g.at(x-5, y) && !g.at(x-6, y) && !g.at(x-7, y-7)
			c = c || !g.at(x+4, y) && !g.at(x-5, y) && !g.at(x+6, y) && !g.at(x, y+7)
			if c {
				cnt++
			}
		}
	}
	return cnt * 40
}

func b2u(c Color) uint {
	if c {
		return 1
	}
	return 0
}

func (img *Image) pointOnesCount() int {
	total := img.Rect.Dx() * img.Rect.Dy()
	cnt := img.OnesCount()
	p := float64(cnt)/float64(total) - 0.5
	if p < 0 {
		p = -p
	}
	return int(p*20) * 10
}
