package qrcode

//go:generate go run genbch/main.go

import (
	"image"

	"github.com/shogo82148/qrcode/internal/bitstream"
	binimage "github.com/shogo82148/qrcode/internal/image"
)

func Generate() image.Image {
	w := 21
	img := binimage.New(image.Rect(0, 0, w, w))
	used := binimage.New(image.Rect(0, 0, w, w))

	for i := 0; i < 7; i++ {
		img.SetBinary(i, 0, binimage.Black)
		img.SetBinary(0, i, binimage.Black)
		img.SetBinary(i, 6, binimage.Black)
		img.SetBinary(6, i, binimage.Black)

		img.SetBinary(w-i-1, 0, binimage.Black)
		img.SetBinary(w-0-1, i, binimage.Black)
		img.SetBinary(w-i-1, 6, binimage.Black)
		img.SetBinary(w-6-1, i, binimage.Black)

		img.SetBinary(0, w-i-1, binimage.Black)
		img.SetBinary(i, w-0-1, binimage.Black)
		img.SetBinary(6, w-i-1, binimage.Black)
		img.SetBinary(i, w-6-1, binimage.Black)
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			img.SetBinary(i+2, j+2, binimage.Black)
			img.SetBinary(w-i-3, j+2, binimage.Black)
			img.SetBinary(i+2, w-j-3, binimage.Black)
		}
	}

	// finder pattern
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			used.SetBinary(i, j, true)
			used.SetBinary(w-i-1, j, true)
			used.SetBinary(i, w-j-1, true)
		}
	}

	// model number
	for i := 0; i < 8; i++ {
		used.SetBinary(i, 8, true)
		used.SetBinary(8, i, true)
		used.SetBinary(8, w-i-1, true)
		used.SetBinary(w-i-1, 8, true)
	}
	used.SetBinary(8, 8, true)

	// timing pattern
	for i := 6; i < w-6; i++ {
		img.SetBinary(i, 6, i%2 == 0)
		img.SetBinary(6, i, i%2 == 0)
		used.SetBinary(i, 6, true)
		used.SetBinary(6, i, true)
	}

	buffer := bitstream.NewBuffer([]byte{
		0b0001_0000, 0b0010_0000, 0b0000_1100, 0b0101_0110,
		0b0110_0001, 0b1000_0000,

		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,

		0b1010_0101, 0b0010_0100, 0b1101_0100, 0b1100_0001,
		0b1110_1101, 0b0011_0110, 0b1100_0111, 0b1000_0111,
		0b0010_1100, 0b0101_0101,
	})

	dy := -1
	x, y := 20, 20
	for {
		if x == 6 {
			// skip timing pattern
			x--
			continue
		}
		if !used.BinaryAt(x, y) {
			bit, err := buffer.ReadBit()
			if err != nil {
				break
			}
			img.SetBinary(x, y, bit != 0)
		}
		x--
		if x < 0 {
			break
		}

		if !used.BinaryAt(x, y) {
			bit, err := buffer.ReadBit()
			if err != nil {
				break
			}
			img.SetBinary(x, y, bit != 0)
		}
		x, y = x+1, y+dy
		if y < 0 || y >= w {
			dy *= -1
			x, y = x-2, y+dy
		}
		if x < 0 {
			break
		}
	}

	model := encodedModel[0b00_010]
	img.SetBinary(0, 8, (model>>14)&1 != 0)
	img.SetBinary(1, 8, (model>>13)&1 != 0)
	img.SetBinary(2, 8, (model>>12)&1 != 0)
	img.SetBinary(3, 8, (model>>11)&1 != 0)
	img.SetBinary(4, 8, (model>>10)&1 != 0)
	img.SetBinary(5, 8, (model>>9)&1 != 0)
	img.SetBinary(7, 8, (model>>8)&1 != 0)
	img.SetBinary(8, 8, (model>>7)&1 != 0)
	img.SetBinary(8, 7, (model>>6)&1 != 0)
	img.SetBinary(8, 5, (model>>5)&1 != 0)
	img.SetBinary(8, 4, (model>>4)&1 != 0)
	img.SetBinary(8, 3, (model>>3)&1 != 0)
	img.SetBinary(8, 2, (model>>2)&1 != 0)
	img.SetBinary(8, 1, (model>>1)&1 != 0)
	img.SetBinary(8, 0, (model>>0)&1 != 0)

	img.SetBinary(8, w-1, (model>>14)&1 != 0)
	img.SetBinary(8, w-2, (model>>13)&1 != 0)
	img.SetBinary(8, w-3, (model>>12)&1 != 0)
	img.SetBinary(8, w-4, (model>>11)&1 != 0)
	img.SetBinary(8, w-5, (model>>10)&1 != 0)
	img.SetBinary(8, w-6, (model>>9)&1 != 0)
	img.SetBinary(8, w-7, (model>>8)&1 != 0)
	img.SetBinary(8, w-8, binimage.Black)
	img.SetBinary(w-8, 8, (model>>7)&1 != 0)
	img.SetBinary(w-7, 8, (model>>6)&1 != 0)
	img.SetBinary(w-6, 8, (model>>5)&1 != 0)
	img.SetBinary(w-5, 8, (model>>4)&1 != 0)
	img.SetBinary(w-4, 8, (model>>3)&1 != 0)
	img.SetBinary(w-3, 8, (model>>2)&1 != 0)
	img.SetBinary(w-2, 8, (model>>1)&1 != 0)
	img.SetBinary(w-1, 8, (model>>0)&1 != 0)

	for i := 0; i < w; i++ {
		for j := 0; j < w; j++ {
			if used.BinaryAt(i, j) {
				continue
			}
			if i%3 == 0 {
				c := img.BinaryAt(i, j)
				img.SetBinary(i, j, !c)
			}
		}
	}
	return img
}
