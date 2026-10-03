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

// Point returns the penalty score of the QR code.
// The lower score is better.
// See ISO/IEC 18004:2015 7.8.3 Evaluation of data masking results.
func (img *Image) Point() int {
	g := newGrid(img)
	return g.longRunLengthCount() + g.blockCount() + g.finderPattern() + img.pointOnesCount()
}

// grid is an unpacked copy of Image for fast random access.
// It is used for calculating the penalty score.
type grid struct {
	pix  []uint8 // 1 is black, 0 is white
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
		pix: pix,
		w:   w,
		h:   h,
	}
}

// at returns the color of the pixel at (x, y).
// (x, y) is the relative position from the top-left corner.
// The pixels outside of the image are white (light modules of the quiet zone).
func (g *grid) at(x, y int) uint8 {
	if uint(x) >= uint(g.w) || uint(y) >= uint(g.h) {
		return 0
	}
	return g.pix[y*g.w+x]
}

// PointMicro returns the evaluation score of the Micro QR code.
// The higher score is better.
// See ISO/IEC 18004:2015 7.8.3.2 Evaluation of Micro QR Code symbols.
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

// longRunLengthCount calculates the penalty of
// adjacent modules in row/column in same color.
// The penalty is 3 + i points for each run of 5 + i modules.
func (g *grid) longRunLengthCount() int {
	var cnt int
	w, h := g.w, g.h
	if w == 0 || h == 0 {
		return 0
	}

	// horizontal
	for y := 0; y < h; y++ {
		line := g.pix[y*w : (y+1)*w]
		length := 0
		c0 := line[0]
		for _, c := range line {
			length, cnt = runLength(length, cnt, c^c0)
			c0 = c
		}
		cnt += runLengthPoint(length)
	}

	// vertical
	for x := 0; x < w; x++ {
		length := 0
		c0 := g.pix[x]
		for i := x; i < len(g.pix); i += w {
			c := g.pix[i]
			length, cnt = runLength(length, cnt, c^c0)
			c0 = c
		}
		cnt += runLengthPoint(length)
	}

	return cnt
}

// runLength is a branch-less version of the following code:
//
//	if diff != 0 {
//		cnt += runLengthPoint(length)
//		length = 0
//	}
//	length++
func runLength(length, cnt int, diff uint8) (int, int) {
	d := int(diff)                                    // 1 if the color changes, otherwise 0
	ge5 := int(uint(4-length) >> (bits.UintSize - 1)) // 1 if length >= 5, otherwise 0
	cnt += (length - 2) & -(d & ge5)
	length = (length & (d - 1)) + 1
	return length, cnt
}

func runLengthPoint(length int) int {
	if length >= 5 {
		return length - 5 + 3
	}
	return 0
}

// blockCount calculates the penalty of blocks of modules in same color.
// The penalty is 3 points for each 2x2 block.
func (g *grid) blockCount() int {
	var cnt int
	w, h := g.w, g.h
	for y := 0; y < h-1; y++ {
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

// finderPattern calculates the penalty of finder-like patterns.
// The penalty is 40 points for each 1:1:3:1:1 (dark:light:dark:light:dark) pattern
// in row/column, preceded or followed by light area 4 modules wide.
func (g *grid) finderPattern() int {
	const pattern = 0b1011101
	var cnt int

	// vertical
	for x := 0; x < g.w; x++ {
		// window holds the pixels from (x, y-3) to (x, y+3).
		// the pixel (x, y-3) is the most significant bit.
		var window uint
		for yy := -3; yy < g.h+3; yy++ {
			window = (window<<1 | uint(g.at(x, yy))) & 0x7f
			y := yy - 3
			if window != pattern {
				continue
			}
			if g.at(x, y-4)|g.at(x, y-5)|g.at(x, y-6)|g.at(x, y-7) == 0 ||
				g.at(x, y+4)|g.at(x, y+5)|g.at(x, y+6)|g.at(x, y+7) == 0 {
				cnt++
			}
		}
	}

	// horizontal
	for y := 0; y < g.h; y++ {
		// window holds the pixels from (x-3, y) to (x+3, y).
		// the pixel (x-3, y) is the most significant bit.
		var window uint
		for xx := -3; xx < g.w+3; xx++ {
			window = (window<<1 | uint(g.at(xx, y))) & 0x7f
			x := xx - 3
			if window != pattern {
				continue
			}
			if g.at(x-4, y)|g.at(x-5, y)|g.at(x-6, y)|g.at(x-7, y) == 0 ||
				g.at(x+4, y)|g.at(x+5, y)|g.at(x+6, y)|g.at(x+7, y) == 0 {
				cnt++
			}
		}
	}
	return cnt * 40
}

// pointOnesCount calculates the penalty of the proportion of dark modules.
// The penalty is 10 * k points, where k is the rating of
// the deviation of the proportion of dark modules from 50% in steps of 5%.
func (img *Image) pointOnesCount() int {
	total := img.Rect.Dx() * img.Rect.Dy()
	if total == 0 {
		return 0
	}
	cnt := img.OnesCount()

	// k = floor(|cnt / total * 100 - 50| / 5) = floor(|20 * cnt - 10 * total| / total)
	d := 20*cnt - 10*total
	if d < 0 {
		d = -d
	}
	return d / total * 10
}
