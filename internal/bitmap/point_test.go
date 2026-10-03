package bitmap

import (
	"image"
	"math/rand"
	"strings"
	"testing"
)

// refPoint is a straightforward implementation of the penalty score.
// It is used for verifying the optimized implementation.
func (img *Image) refPoint() int {
	r := img.Rect
	var n1, n2, n3 int

	// N1: adjacent modules in row/column in same color
	runs := func(line []Color) int {
		var p int
		length := 1
		for i := 1; i <= len(line); i++ {
			if i < len(line) && line[i] == line[i-1] {
				length++
				continue
			}
			if length >= 5 {
				p += 3 + length - 5
			}
			length = 1
		}
		return p
	}

	// N3: 1:1:3:1:1 pattern preceded or followed by light area 4 modules wide
	finders := func(line []Color) int {
		at := func(i int) Color {
			if i < 0 || i >= len(line) {
				return White
			}
			return line[i]
		}
		pattern := []Color{Black, White, Black, Black, Black, White, Black}
		var p int
	LOOP:
		for i := 0; i+len(pattern) <= len(line); i++ {
			for j, c := range pattern {
				if at(i+j) != c {
					continue LOOP
				}
			}
			before := !at(i-1) && !at(i-2) && !at(i-3) && !at(i-4)
			after := !at(i+7) && !at(i+8) && !at(i+9) && !at(i+10)
			if before || after {
				p += 40
			}
		}
		return p
	}

	for y := r.Min.Y; y < r.Max.Y; y++ {
		var line []Color
		for x := r.Min.X; x < r.Max.X; x++ {
			line = append(line, img.BinaryAt(x, y))
		}
		n1 += runs(line)
		n3 += finders(line)
	}
	for x := r.Min.X; x < r.Max.X; x++ {
		var line []Color
		for y := r.Min.Y; y < r.Max.Y; y++ {
			line = append(line, img.BinaryAt(x, y))
		}
		n1 += runs(line)
		n3 += finders(line)
	}

	// N2: block of modules in same color
	for y := r.Min.Y; y < r.Max.Y-1; y++ {
		for x := r.Min.X; x < r.Max.X-1; x++ {
			c := img.BinaryAt(x, y)
			if c == img.BinaryAt(x+1, y) && c == img.BinaryAt(x, y+1) && c == img.BinaryAt(x+1, y+1) {
				n2 += 3
			}
		}
	}

	return n1 + n2 + n3 + img.pointOnesCount()
}

func parseImage(s string) *Image {
	lines := strings.Fields(s)
	img := New(image.Rect(0, 0, len(lines[0]), len(lines)))
	for y, line := range lines {
		for x, c := range line {
			img.SetBinary(x, y, c == '#')
		}
	}
	return img
}

func TestPoint_Known(t *testing.T) {
	tests := []struct {
		name string
		img  string
		want int
	}{
		{
			// N1: 5 rows + 5 columns of length 5 = 10 * 3
			// N2: 4 * 4 blocks = 16 * 3
			// N4: 0% dark = 10 * 10
			name: "all white",
			img: `
.....
.....
.....
.....
.....`,
			want: 30 + 48 + 100,
		},
		{
			// N1: 6 rows + 6 columns of length 6 = 12 * 4
			// N2: 5 * 5 blocks = 25 * 3
			// N4: 100% dark = 10 * 10
			name: "all black",
			img: `
######
######
######
######
######
######`,
			want: 48 + 75 + 100,
		},
		{
			// N3: one finder-like pattern
			// N4: 5/11 dark = 45.45% -> k = 0
			name: "horizontal finder pattern, light area before",
			img:  `....#.###.#`,
			want: 40,
		},
		{
			name: "horizontal finder pattern, light area after",
			img:  `#.###.#....`,
			want: 40,
		},
		{
			// light area on both sides is counted once.
			// N4: 5/15 = 33.3% -> k = 3
			name: "horizontal finder pattern, light area on both sides",
			img:  `....#.###.#....`,
			want: 40 + 30,
		},
		{
			name: "horizontal finder pattern without light area",
			img:  `#...#.###.#...#`,
			want: 0 + 0,
		},
		{
			name: "vertical finder pattern",
			img: `
.
.
.
.
#
.
#
#
#
.
#`,
			want: 40,
		},
		{
			// N1: a run of length 5 at the end of the row.
			// N4: 3/11 = 27.3% -> k = 4
			name: "long run at the end of the row",
			img:  `.#.#.#.....`,
			want: 3 + 40,
		},
		{
			// N4: 9/20 = exactly 45% -> k = 1
			name: "45% dark",
			img:  `#.#.#.#.#.#.#.#.#...`,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := parseImage(tt.img)
			if got := img.Point(); got != tt.want {
				t.Errorf("Point() = %d, want %d", got, tt.want)
			}
			if got := img.refPoint(); got != tt.want {
				t.Errorf("refPoint() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPoint(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	rects := []image.Rectangle{
		image.Rect(0, 0, 1, 1),
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
