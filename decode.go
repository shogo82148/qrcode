package qrcode

import (
	"errors"
	"image"
	"math"
	"math/bits"

	"github.com/shogo82148/qrcode/internal/binimage"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func Decode(img image.Image) (*QRCode, error) {
	// TODO: find pattern

	bounds := img.Bounds()
	binimg := binimage.New(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	version := (bounds.Dx() - 17) / 4
	w := 16 + 4*version
	for y := 0; y <= w; y++ {
		for x := 0; x <= w; x++ {
			c := imageAt(img, float64(x), float64(y))
			binimg.SetBinary(x, y, c)
		}
	}

	// decode format
	var rawFormat1, rawFormat2 uint
	for i := 0; i < 8; i++ {
		if binimg.BinaryAt(8, skipTimingPattern(i)) {
			rawFormat1 |= 1 << i
		}
		if binimg.BinaryAt(skipTimingPattern(i), 8) {
			rawFormat1 |= 1 << (14 - i)
		}

		if binimg.BinaryAt(w-i, 8) {
			rawFormat2 |= 1 << i
		}
		if binimg.BinaryAt(8, w-i) {
			rawFormat2 |= 1 << (14 - i)
		}
	}
	_, mask, ok := decodeFormat(rawFormat1)
	if !ok {
		return nil, errors.New("qr code not found")
	}

	used := usedList[version]

	// mask
	var f func(i, j int) int
	switch mask {
	case 0b000:
		f = func(i, j int) int { return (i + j) % 2 }
	case 0b001:
		f = func(i, j int) int { return i % 2 }
	case 0b010:
		f = func(i, j int) int { return j % 3 }
	case 0b011:
		f = func(i, j int) int { return (i + j) % 3 }
	case 0b100:
		f = func(i, j int) int { return (i/2 + j/3) % 2 }
	case 0b101:
		f = func(i, j int) int { return i*j%2 + i*j%3 }
	case 0b110:
		f = func(i, j int) int { return (i*j%2 + i*j%3) % 2 }
	case 0b111:
		f = func(i, j int) int { return ((i+j)%2 + i*j%3) % 2 }
	}
	for i := 0; i <= w; i++ {
		for j := 0; j <= w; j++ {
			binimg.XorBinary(j, i, !used.BinaryAt(j, i) && f(i, j) == 0)
		}
	}

	var buf bitstream.Buffer
	dy := -1
	x, y := w, w
	for {
		if x == timingPatternOffset {
			// skip timing pattern
			x--
			continue
		}
		if !used.BinaryAt(x, y) {
			if binimg.BinaryAt(x, y) {
				buf.WriteBit(1)
			} else {
				buf.WriteBit(0)
			}
		}
		x--
		if x < 0 {
			break
		}

		if !used.BinaryAt(x, y) {
			if binimg.BinaryAt(x, y) {
				buf.WriteBit(1)
			} else {
				buf.WriteBit(0)
			}
		}
		x, y = x+1, y+dy
		if y < 0 || y > w {
			dy *= -1
			x, y = x-2, y+dy
		}
		if x < 0 {
			break
		}
	}

	// TODO: un-interleave

	data := buf.Bytes()
	if err := reedsolomon.Decode(data, 2); err != nil {
		return nil, err
	}

	// TODO: decode segment

	return &QRCode{}, nil
}

func imageAt(img image.Image, x, y float64) binimage.Color {
	x = math.Round(x)
	y = math.Round(y)
	c := img.At(int(x), int(y))
	r, g, b, _ := c.RGBA()
	return (r + g + b) < 128*3
}

func decodeFormat(raw uint) (Level, Mask, bool) {
	idx := 0
	min := bits.OnesCount(encodedFormat[0] ^ raw)
	for i, pattern := range encodedFormat {
		count := bits.OnesCount(pattern ^ raw)
		if count < min {
			idx = i
			min = count
		}
	}
	if min >= 3 {
		return 0, 0, false
	}
	return Level(idx >> 3), Mask(idx & 0b111), true
}
