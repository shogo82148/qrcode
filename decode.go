package qrcode

import (
	"errors"
	"fmt"
	"image"
	"io"
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
	version := Version((bounds.Dx() - 17) / 4)
	w := 16 + 4*int(version)
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
	level, mask, ok := decodeFormat(rawFormat1)
	if !ok {
		return nil, errors.New("qr code not found")
	}

	used := usedList[version]

	// mask
	binimg.Mask(binimg, used, maskList[mask])

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

	// un-interleave
	blocks := decodeFromBits(version, level, buf.Bytes())

	// error correction
	var result []byte
	for _, blk := range blocks {
		data := append(blk.data, blk.correction...)
		if err := reedsolomon.Decode(data, 2); err != nil {
			return nil, err
		}
		result = append(result, data[:len(blk.data)]...)
	}

	// decode segments
	stream := bitstream.NewBuffer(result)
	segments := make([]Segment, 0)
LOOP:
	for {
		mode, err := stream.ReadBits(4)
		if err == io.EOF {
			break
		}
		switch Mode(mode) {
		case ModeNumber:
			seg, err := decodeNumber(version, stream)
			if err != nil {
				return nil, err
			}
			segments = append(segments, seg)
		case ModeAlphanumeric:
			seg, err := decodeAlphanumeric(version, stream)
			if err != nil {
				return nil, err
			}
			segments = append(segments, seg)
		case ModeBytes:
			seg, err := decodeBytes(version, stream)
			if err != nil {
				return nil, err
			}
			segments = append(segments, seg)
		case ModeTerminated:
			break LOOP
		}
	}

	return &QRCode{
		Version:  version,
		Mask:     mask,
		Level:    level,
		Segments: segments,
	}, nil
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

func decodeFromBits(version Version, level Level, buf []byte) []block {
	capacity := capacityTable[version][level]
	blocks := []block{}
	for _, blockCapacity := range capacity.Blocks {
		for i := 0; i < blockCapacity.Num; i++ {
			blocks = append(blocks, block{
				data:       make([]byte, blockCapacity.Data),
				correction: make([]byte, blockCapacity.Total-blockCapacity.Data),
			})
		}
	}

	i := 0
	for _, b := range buf[:capacity.Data] {
		for {
			if i/len(blocks) < len(blocks[i%len(blocks)].data) {
				blocks[i%len(blocks)].data[i/len(blocks)] = b
				i++
				break
			}
			i++
		}
	}

	i = 0
	for _, b := range buf[capacity.Data:capacity.Total] {
		for {
			if i/len(blocks) < len(blocks[i%len(blocks)].correction) {
				blocks[i%len(blocks)].correction[i/len(blocks)] = b
				i++
				break
			}
			i++
		}
	}
	return blocks
}

func decodeNumber(version Version, buf *bitstream.Buffer) (Segment, error) {
	var n int
	switch {
	case version <= 0 || version > 40:
		return Segment{}, fmt.Errorf("qrcode: invalid version: %d", version)
	case version < 10:
		n = 10
	case version < 27:
		n = 12
	default:
		n = 14
	}
	length, err := buf.ReadBits(n)
	if err != nil {
		return Segment{}, err
	}
	data := make([]byte, 0, length)

	for i := uint64(0); i+2 < length; i += 3 {
		bits, err := buf.ReadBits(10)
		if err != nil {
			return Segment{}, err
		}
		if bits >= 1000 {
			return Segment{}, errors.New("invalid digit")
		}
		n1 := bits / 1000
		n2 := bits / 100 % 10
		n3 := bits % 10
		data = append(data, byte(n1+'0'), byte(n2+'0'), byte(n3+'0'))
	}

	switch len(data) % 3 {
	case 1:
		bits, err := buf.ReadBits(4)
		if err != nil {
			return Segment{}, err
		}
		if bits >= 10 {
			return Segment{}, errors.New("invalid digit")
		}
		data = append(data, byte(bits+'0'))
	case 2:
		bits, err := buf.ReadBits(7)
		if err != nil {
			return Segment{}, err
		}
		if bits >= 100 {
			return Segment{}, errors.New("invalid digit")
		}
		n1 := bits / 10
		n2 := bits % 10
		data = append(data, byte(n1+'0'), byte(n2+'0'))
	}

	return Segment{
		Mode: ModeNumber,
		Data: data,
	}, nil
}

func decodeAlphanumeric(version Version, buf *bitstream.Buffer) (Segment, error) {
	var n int
	switch {
	case version <= 0 || version > 40:
		return Segment{}, fmt.Errorf("qrcode: invalid version: %d", version)
	case version < 10:
		n = 9
	case version < 27:
		n = 11
	default:
		n = 13
	}

	length, err := buf.ReadBits(n)
	if err != nil {
		return Segment{}, err
	}
	data := make([]byte, 0)
	for i := uint64(0); i+1 < length; i += 2 {
		bits, err := buf.ReadBits(11)
		if err != nil {
			return Segment{}, err
		}
		n1 := int(bits) / 45
		n2 := int(bits) % 45
		if n1 >= 45 {
			return Segment{}, errors.New("invalid digit")
		}
		data = append(data, bitToAlphanumeric[n1], bitToAlphanumeric[n2])
	}

	if length%2 != 0 {
		bits, err := buf.ReadBits(6)
		if err != nil {
			return Segment{}, err
		}
		if bits >= 45 {
			return Segment{}, errors.New("invalid digit")
		}
		data = append(data, bitToAlphanumeric[bits])
	}

	return Segment{
		Mode: ModeAlphanumeric,
		Data: data,
	}, nil
}

func decodeBytes(version Version, buf *bitstream.Buffer) (Segment, error) {
	var n int
	switch {
	case version <= 0 || version > 40:
		return Segment{}, fmt.Errorf("qrcode: invalid version: %d", version)
	case version < 10:
		n = 8
	default:
		n = 16
	}
	length, err := buf.ReadBits(n)
	if err != nil {
		return Segment{}, err
	}
	data := make([]byte, length)
	for i := uint64(0); i < length; i++ {
		bits, err := buf.ReadBits(8)
		if err != nil {
			return Segment{}, err
		}
		data[i] = byte(bits)
	}

	return Segment{
		Mode: ModeBytes,
		Data: data,
	}, nil
}
