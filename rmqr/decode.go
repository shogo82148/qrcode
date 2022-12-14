// Package rmqr handles rMRQ Codes.
package rmqr

import (
	"errors"
	"log"
	"math/bits"

	"github.com/shogo82148/qrcode/bitmap"
	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
)

func DecodeBitmap(img *bitmap.Image) (*QRCode, error) {
	binimg := internalbitmap.Import(img)

	var rawVersion uint
	for i := 0; i < 18; i++ {
		if binimg.BinaryAt(8+i/5, 1+i%5) {
			rawVersion |= 1 << i
		}
	}
	version, _, ok := decodeFormat(rawVersion ^ 0b011111101010110010)
	if !ok {
		return nil, errors.New("rmqr: rMQR not found")
	}
	// TODO: check around sub-finder pattern

	bounds := img.Bounds()
	w := bounds.Dx() - 1
	h := bounds.Dy() - 1

	used := usedList[version]
	binimg.Mask(binimg, used, precomputedMask)

	var buf bitstream.Buffer
	dy := -1
	x, y := w-1, h-5
	for {
		if !used.BinaryAt(x, y) {
			if binimg.BinaryAt(x, y) {
				buf.WriteBit(1)
			} else {
				buf.WriteBit(0)
			}
		}
		x--
		if x < 1 { // +1 is for avoiding time pattern
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
		if y < +1 || y > h-1 { // +1 and -1 are for avoiding time pattern
			dy *= -1
			x, y = x-2, y+dy
		}
		if x < 1 { // +1 is for avoiding time pattern
			break
		}
	}

	// data := buf.Bytes()
	// if err := reedsolomon.Decode(data, 2); err != nil {
	// 	return nil, err
	// }

	mode, err := buf.ReadBits(3)
	if err != nil {
		return nil, err
	}
	_ = mode
	length, err := buf.ReadBits(4)
	if err != nil {
		return nil, err
	}
	data := make([]byte, length)
	if err := bitstream.DecodeNumeric(&buf, data); err != nil {
		return nil, err
	}
	log.Println(string(data))

	return &QRCode{}, nil
}

func decodeFormat(data uint) (Version, Level, bool) {
	var version, min int
	min = bits.OnesCount(encodedVersion[0] ^ data)
	for i, v := range encodedVersion[1:] {
		diff := bits.OnesCount(v ^ data)
		if diff < min {
			version = i
			min = diff
		}
	}
	if min >= 3 {
		return 0, 0, false
	}
	return Version(version & 0x1f), Level((version >> 5) & 1), true
}
