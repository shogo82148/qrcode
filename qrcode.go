package qrcode

//go:generate go run genbch/main.go

import (
	"image"

	binimage "github.com/shogo82148/qrcode/internal/image"
)

// 0x00000001
// 0x01111101
// 0x01000101
// 0x01000101
// 0x01000101
// 0x01111101
// 0x00000001
// 0x11111111

func Generate() image.Image {
	w := 21
	img := binimage.New(image.Rect(0, 0, w, w))

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

	for i := 6; i < w-6; i++ {
		img.SetBinary(i, 6, i%2 == 0)
		img.SetBinary(6, i, i%2 == 0)
	}
	return img
}
