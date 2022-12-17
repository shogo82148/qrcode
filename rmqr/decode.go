// Package rmqr handles rMRQ Codes.
package rmqr

import (
	"errors"
	"io"
	"math/bits"

	"github.com/shogo82148/qrcode/bitmap"
	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func DecodeBitmap(img *bitmap.Image) (*QRCode, error) {
	binimg := internalbitmap.Import(img)
	bounds := img.Bounds()
	w := bounds.Dx() - 1
	h := bounds.Dy() - 1

	version, level, err := decodeFormat(binimg)
	_ = level
	if err != nil {
		return nil, err
	}
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
	bitLength := capacityTable[version][level].BitLength
LOOP:
	for {
		mode, err := stream.ReadBits(3)
		if err == io.EOF {
			break
		}
		switch Mode(mode) {
		case ModeNumeric:
			seg, err := decodeNumber(bitLength, stream)
			if err != nil {
				return nil, err
			}
			segments = append(segments, seg)
		case ModeAlphanumeric:
			seg, err := decodeAlphanumeric(bitLength, stream)
			if err != nil {
				return nil, err
			}
			segments = append(segments, seg)
		case ModeBytes:
			seg, err := decodeBytes(bitLength, stream)
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
		Level:    level,
		Segments: segments,
	}, nil
}

func decodeFormat(img *internalbitmap.Image) (Version, Level, error) {
	bounds := img.Rect
	w := bounds.Dx() - 1
	h := bounds.Dy() - 1

	// search version info around finder pattern
	var rawVersion uint
	for i := 0; i < 18; i++ {
		if img.BinaryAt(8+i/5, 1+i%5) {
			rawVersion |= 1 << i
		}
	}
	version, level, ok := decodeFormat0(rawVersion ^ 0b011111101010110010)
	if !ok {
		return version, level, nil
	}

	// search version info around sub-finder pattern
	rawVersion = 0
	for i := 0; i < 15; i++ {
		if img.BinaryAt(w-7+i/5, h-5+i%5) {
			rawVersion |= 1 << i
		}
	}
	if img.BinaryAt(w-4, h-5) {
		rawVersion |= 1 << 15
	}
	if img.BinaryAt(w-3, h-5) {
		rawVersion |= 1 << 16
	}
	if img.BinaryAt(w-2, h-5) {
		rawVersion |= 1 << 17
	}
	version, level, ok = decodeFormat0(rawVersion ^ 0b100000101001111011)
	if ok {
		return version, level, nil
	}

	return 0, 0, errors.New("rmqr: rMRQ not found")
}

func decodeFormat0(data uint) (Version, Level, bool) {
	var idx, min int
	min = bits.OnesCount(encodedVersion[0] ^ data)
	for i, v := range encodedVersion {
		diff := bits.OnesCount(v ^ data)
		if diff < min {
			idx = i
			min = diff
		}
	}
	if min >= 3 {
		return 0, 0, false
	}
	return Version(idx & 0x1f), Level((idx >> 6) & 1), true
}

func decodeFromBits(version Version, level Level, buf []byte) []block {
	capacity := capacityTable[version][level]
	blocks := []block{}
	for _, blockCapacity := range capacity.Blocks {
		for i := 0; i < blockCapacity.Num; i++ {
			blocks = append(blocks, block{
				data:       make([]byte, blockCapacity.Data),
				correction: make([]byte, blockCapacity.Total-blockCapacity.Data),
				maxError:   blockCapacity.MaxError,
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

func decodeNumber(bitLength [5]int, buf *bitstream.Buffer) (Segment, error) {
	n := bitLength[ModeNumeric]
	length, err := buf.ReadBits(n)
	if err != nil {
		return Segment{}, err
	}
	data := make([]byte, length)
	if err := bitstream.DecodeNumeric(buf, data); err != nil {
		return Segment{}, err
	}

	return Segment{
		Mode: ModeNumeric,
		Data: data,
	}, nil
}

func decodeAlphanumeric(bitLength [5]int, buf *bitstream.Buffer) (Segment, error) {
	n := bitLength[ModeAlphanumeric]
	length, err := buf.ReadBits(n)
	if err != nil {
		return Segment{}, err
	}
	data := make([]byte, length)
	if err := bitstream.DecodeAlphanumeric(buf, data); err != nil {
		return Segment{}, err
	}

	return Segment{
		Mode: ModeAlphanumeric,
		Data: data,
	}, nil
}

func decodeBytes(bitLength [5]int, buf *bitstream.Buffer) (Segment, error) {
	n := bitLength[ModeBytes]
	length, err := buf.ReadBits(n)
	if err != nil {
		return Segment{}, err
	}
	data := make([]byte, length)
	if err := bitstream.DecodeBytes(buf, data); err != nil {
		return Segment{}, err
	}

	return Segment{
		Mode: ModeBytes,
		Data: data,
	}, nil
}
