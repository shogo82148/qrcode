package main

import (
	"bytes"
	"fmt"
	"go/format"
	"image"
	"log"
	"os"

	"github.com/shogo82148/qrcode/internal/bitmap"
)

var size = [32]struct{ h, w int }{
	{7, 43},
	{7, 59},
	{7, 77},
	{7, 99},
	{7, 139},

	{9, 43},
	{9, 59},
	{9, 77},
	{9, 99},
	{9, 139},

	{11, 27},
	{11, 43},
	{11, 59},
	{11, 77},
	{11, 99},
	{11, 139},

	{13, 27},
	{13, 43},
	{13, 59},
	{13, 77},
	{13, 99},
	{13, 139},

	{15, 43},
	{15, 59},
	{15, 77},
	{15, 99},
	{15, 139},

	{17, 43},
	{17, 59},
	{17, 77},
	{17, 99},
	{17, 139},
}

var alignmentPatternPositions = map[int][]int{
	27:  {},
	43:  {21},
	59:  {19, 39},
	77:  {25, 51},
	99:  {23, 49, 75},
	139: {27, 55, 83, 111},
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by genbase/main.go; DO NOT EDIT.\n\n")
	fmt.Fprintf(&buf, "package rmqr\n\n")
	fmt.Fprintf(&buf, "import (\n")
	fmt.Fprintf(&buf, "\"image\"\n")
	fmt.Fprintf(&buf, "\"github.com/shogo82148/qrcode/internal/bitmap\"\n")
	fmt.Fprintf(&buf, ")\n\n")

	genMaskList(&buf)
	genBaseList(&buf)

	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("base_gen.go", out, 0o644); err != nil {
		log.Fatal(err)
	}
}

func genMaskList(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "var precomputedMask = &bitmap.Image{\n")
	genMask(buf, func(i, j int) int { return (i/2 + j/3) % 2 })
	fmt.Fprintf(buf, "}\n\n")
}

func genMask(buf *bytes.Buffer, f func(i, j int) int) {
	w, h := 144, 17
	img := bitmap.New(image.Rect(0, 0, w, h))
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			img.SetBinary(j, i, f(i, j) == 0)
		}
	}
	writeImageContents(buf, img)
}

func genBaseList(buf *bytes.Buffer) {
	imgList := []*bitmap.Image{}
	usedList := []*bitmap.Image{}

	for version := 0; version < 32; version++ {
		img, used := newBase(version)
		imgList = append(imgList, img)
		usedList = append(usedList, used)
	}

	fmt.Fprintf(buf, "var baseList = []*bitmap.Image{\n")
	for version := 0; version < 32; version++ {
		fmt.Fprintf(buf, "\n// version %d\n", version)
		writeImage(buf, imgList[version])
	}
	fmt.Fprintf(buf, "}\n\n")

	fmt.Fprintf(buf, "var usedList = []*bitmap.Image{\n")
	for version := 0; version < 32; version++ {
		fmt.Fprintf(buf, "\n// version %d\n", version)
		writeImage(buf, usedList[version])
	}
	fmt.Fprintf(buf, "}\n")
}

func newBase(version int) (*bitmap.Image, *bitmap.Image) {
	w := size[version].w - 1
	h := size[version].h - 1
	img := bitmap.New(image.Rect(0, 0, w+1, h+1))
	used := bitmap.New(image.Rect(0, 0, w+1, h+1))

	// timing pattern
	for i := 0; i <= w; i++ {
		img.SetBinary(i, 0, i%2 == 0)
		img.SetBinary(i, h, i%2 == 0)
		used.SetBinary(i, 0, bitmap.Black)
		used.SetBinary(i, h, bitmap.Black)
	}
	for _, pos := range append([]int{0, w}, alignmentPatternPositions[w+1]...) {
		for i := 0; i <= h; i++ {
			img.SetBinary(pos, i, i%2 == 0)
			used.SetBinary(pos, i, bitmap.Black)
		}
	}

	// corner finder pattern
	// bottom left
	img.SetBinary(1, h-0, bitmap.Black)
	img.SetBinary(0, h-0, bitmap.Black)
	img.SetBinary(0, h-1, bitmap.Black)
	img.SetBinary(1, h-1, bitmap.White)
	used.SetBinary(1, h-0, bitmap.Black)
	used.SetBinary(0, h-0, bitmap.Black)
	used.SetBinary(0, h-1, bitmap.Black)
	used.SetBinary(1, h-1, bitmap.Black)

	// Top right
	img.SetBinary(w-1, 0, bitmap.Black)
	img.SetBinary(w-0, 0, bitmap.Black)
	img.SetBinary(w-0, 1, bitmap.Black)
	img.SetBinary(w-1, 1, bitmap.White)
	used.SetBinary(w-1, 0, bitmap.Black)
	used.SetBinary(w-0, 0, bitmap.Black)
	used.SetBinary(w-0, 1, bitmap.Black)
	used.SetBinary(w-1, 1, bitmap.Black)

	// finder pattern
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			d := max(abs(x-3), abs(y-3))
			c := bitmap.Color(d != 2 && d != 4)
			img.SetBinary(x, y, c)
			used.SetBinary(x, y, true)
		}
	}

	// sub-finder pattern
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			d := max(abs(x-2), abs(y-2))
			c := bitmap.Color(d != 1)
			img.SetBinary(w-x, h-y, c)
			used.SetBinary(w-x, h-y, true)
		}
	}

	// alignment patterns
	for _, pos := range alignmentPatternPositions[w+1] {
		for y := 0; y <= 2; y++ {
			for x := -1; x <= 1; x++ {
				img.SetBinary(pos+x, y, x != 0 || y != 1)
				img.SetBinary(pos+x, h-y, x != 0 || y != 1)
				used.SetBinary(pos+x, y, bitmap.Black)
				used.SetBinary(pos+x, h-y, bitmap.Black)
			}
		}
	}

	// format information
	for i := 0; i < 18; i++ {
		// around finder pattern
		used.SetBinary(8+i/5, 1+i%5, bitmap.Black)

		// around sub-finder pattern
		used.SetBinary(w-7+i/5, h-5+i%5, bitmap.Black)
	}
	used.SetBinary(w-4, h-5, bitmap.Black)
	used.SetBinary(w-3, h-5, bitmap.Black)
	used.SetBinary(w-2, h-5, bitmap.Black)

	return img, used
}

func writeImage(buf *bytes.Buffer, img *bitmap.Image) {
	fmt.Fprintf(buf, "{\n")
	writeImageContents(buf, img)
	fmt.Fprintf(buf, "},\n")
}

func writeImageContents(buf *bytes.Buffer, img *bitmap.Image) {
	fmt.Fprintf(buf, "Stride: %d,\n", img.Stride)
	fmt.Fprintf(buf, "Rect: image.Rect(%d, %d, %d, %d),\n",
		img.Rect.Min.X, img.Rect.Min.Y, img.Rect.Max.X, img.Rect.Max.Y)
	fmt.Fprintf(buf, "Pix: []byte{\n")
	for i := 0; i < len(img.Pix); i += img.Stride {
		for _, b := range img.Pix[i : i+img.Stride] {
			fmt.Fprintf(buf, "0b%08b, ", b)
		}
		fmt.Fprintln(buf)
	}
	fmt.Fprintf(buf, "},\n")
}
