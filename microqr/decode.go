package microqr

import (
	"errors"
	"log"
	"math/bits"

	"github.com/shogo82148/qrcode/bitmap"
	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func DecodeBitmap(img *bitmap.Image) (*QRCode, error) {
	bounds := img.Bounds()
	version := Version((bounds.Dx() - 9) / 2)
	binimg := internalbitmap.Import(img)
	w := 8 + 2*int(version)

	// decode format
	var rawFormat uint
	for i := 0; i < 8; i++ {
		if binimg.BinaryAt(8, i+1) {
			rawFormat |= 1 << i
		}
		if binimg.BinaryAt(i+1, 8) {
			rawFormat |= 1 << (14 - i)
		}
	}
	_, level, mask, ok := decodeFormat(rawFormat)
	if !ok {
		return nil, errors.New("qr code not found")
	}
	_ = level

	used := usedList[version]

	// mask
	binimg.Mask(binimg, used, maskList[mask])

	var buf bitstream.Buffer
	dy := -1
	x, y := w, w
	for {
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

	data := buf.Bytes()
	log.Printf("%08b", data)
	if err := reedsolomon.Decode(data, 2); err != nil {
		return nil, err
	}

	return &QRCode{}, nil
}

func decodeFormat(raw uint) (Version, Level, Mask, bool) {
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
		return 0, 0, 0, false
	}
	format := rawFormatTable[idx>>2]
	return format.version, format.level, Mask(idx & 0b11), true
}
