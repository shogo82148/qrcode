package bitmap

import (
	"image"
	"math/rand"
	"testing"
)

// The reference implementation of the penalty score, kept to verify that
// the optimized implementation returns exactly the same results.

func (img *Image) refPoint() int {
	return img.refFinderPattern() + img.refLongRunLengthCount() + img.refBlockCount() + img.pointOnesCount()
}

func (img *Image) refLongRunLengthCount() int {
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

func (img *Image) refBlockCount() int {
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

func (img *Image) refFinderPattern() int {
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

func TestPoint(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	rects := []image.Rectangle{
		image.Rect(0, 0, 11, 11),
		image.Rect(0, 0, 13, 13),
		image.Rect(0, 0, 17, 17),
		image.Rect(0, 0, 21, 21),
		image.Rect(0, 0, 25, 25),
		image.Rect(0, 0, 57, 57),
		image.Rect(0, 0, 177, 177),

		// non-square images and images with non-zero origin
		image.Rect(0, 0, 27, 11),
		image.Rect(0, 0, 11, 27),
		image.Rect(2, 3, 23, 24),
		image.Rect(3, 3, 24, 24),
	}
	for _, rect := range rects {
		for _, density := range []int{2, 4, 8, 16} {
			for i := 0; i < 20; i++ {
				img := New(rect)
				for y := rect.Min.Y; y < rect.Max.Y; y++ {
					for x := rect.Min.X; x < rect.Max.X; x++ {
						// sparse and dense images exercise the finder pattern and the long run length rules.
						img.SetBinary(x, y, r.Intn(density) == 0 != (density > 4 && i%2 == 0))
					}
				}
				// embed finder-like patterns
				for j := 0; j < 3; j++ {
					x0 := rect.Min.X + r.Intn(rect.Dx())
					y0 := rect.Min.Y + r.Intn(rect.Dy())
					for k, c := range []Color{Black, White, Black, Black, Black, White, Black} {
						img.SetBinary(x0+k, y0, c)
						img.SetBinary(x0, y0+k, c)
					}
				}
				if got, want := img.Point(), img.refPoint(); got != want {
					t.Errorf("rect %v, density %d, #%d: got %d, want %d", rect, density, i, got, want)
				}
			}
		}
	}
}

func BenchmarkPoint(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	img := New(image.Rect(0, 0, 177, 177))
	for y := 0; y < 177; y++ {
		for x := 0; x < 177; x++ {
			img.SetBinary(x, y, r.Intn(2) == 0)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.Point()
	}
}
