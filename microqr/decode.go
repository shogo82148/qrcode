package microqr

import (
	"errors"
	"io"
	"math/bits"
	"strconv"

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
	version, level, mask, ok := decodeFormat(rawFormat)
	if !ok {
		return nil, errors.New("qr code not found")
	}

	w = 8 + 2*int(version)
	used := usedList[version]

	// mask
	binimg.Mask(binimg, used, maskList[mask])

	qrCapacity := capacityTable[version][level]
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
		if buf.Len() == qrCapacity.DataBits {
			for buf.Len()%8 != 0 {
				buf.WriteBit(0)
			}
		}
	}

	data := buf.Bytes()
	if err := reedsolomon.Decode(data, qrCapacity.MaxError); err != nil {
		return nil, err
	}
	data = data[:qrCapacity.Data]
	buf0 := bitstream.NewBuffer(data)

	switch version {
	case 1:
		return decodeVersion1(buf0, mask, level)
	case 2:
		return decodeVersion2(buf0, mask, level)
	case 3:
		return decodeVersion3(buf0, mask, level)
	case 4:
		return decodeVersion4(buf0, mask, level)
	default:
		panic("invalid version: " + strconv.Itoa(int(version)))
	}
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

func decodeVersion1(buf *bitstream.Buffer, mask Mask, level Level) (*QRCode, error) {
	segments := make([]Segment, 0)
	for {
		length, err := buf.ReadBits(3)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if length == 0 { // terminate pattern
			break
		}
		data := make([]byte, length)
		if err := bitstream.DecodeNumeric(buf, data); err != nil {
			return nil, err
		}
		segments = append(segments, Segment{
			Mode: ModeNumeric,
			Data: data,
		})
	}
	return &QRCode{
		Version:  1,
		Level:    level,
		Mask:     mask,
		Segments: segments,
	}, nil
}

func decodeVersion2(buf *bitstream.Buffer, mask Mask, level Level) (*QRCode, error) {
	segments := make([]Segment, 0)

LOOP:
	for {
		mode, err := buf.ReadBits(1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		var data []byte
		switch Mode(mode) {
		case ModeNumeric:
			length, err := buf.ReadBits(4)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			if length == 0 { // terminate pattern
				break LOOP
			}
			data = make([]byte, length)
			if err := bitstream.DecodeNumeric(buf, data); err != nil {
				return nil, err
			}
		case ModeAlphanumeric:
			length, err := buf.ReadBits(3)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			data = make([]byte, length)
			if err := bitstream.DecodeAlphanumeric(buf, data); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("qrcode: unknown mode: " + strconv.Itoa(int(mode)))
		}
		segments = append(segments, Segment{
			Mode: Mode(mode),
			Data: data,
		})
	}
	return &QRCode{
		Version:  2,
		Level:    level,
		Mask:     mask,
		Segments: segments,
	}, nil
}

func decodeVersion3(buf *bitstream.Buffer, mask Mask, level Level) (*QRCode, error) {
	segments := make([]Segment, 0)

LOOP:
	for {
		mode, err := buf.ReadBits(2)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		var data []byte
		switch Mode(mode) {
		case ModeNumeric:
			length, err := buf.ReadBits(5)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			if length == 0 { // terminate pattern
				break LOOP
			}
			data = make([]byte, length)
			if err := bitstream.DecodeNumeric(buf, data); err != nil {
				return nil, err
			}
		case ModeAlphanumeric:
			length, err := buf.ReadBits(4)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			data = make([]byte, length)
			if err := bitstream.DecodeAlphanumeric(buf, data); err != nil {
				return nil, err
			}
		case ModeBytes:
			length, err := buf.ReadBits(4)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			data = make([]byte, length)
			if err := bitstream.DecodeBytes(buf, data); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("qrcode: unknown mode: " + strconv.Itoa(int(mode)))
		}
		segments = append(segments, Segment{
			Mode: Mode(mode),
			Data: data,
		})
	}
	return &QRCode{
		Version:  3,
		Level:    level,
		Mask:     mask,
		Segments: segments,
	}, nil
}

func decodeVersion4(buf *bitstream.Buffer, mask Mask, level Level) (*QRCode, error) {
	segments := make([]Segment, 0)

LOOP:
	for {
		mode, err := buf.ReadBits(3)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		var data []byte
		switch Mode(mode) {
		case ModeNumeric:
			length, err := buf.ReadBits(6)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			if length == 0 { // terminate pattern
				break LOOP
			}
			data = make([]byte, length)
			if err := bitstream.DecodeNumeric(buf, data); err != nil {
				return nil, err
			}
		case ModeAlphanumeric:
			length, err := buf.ReadBits(5)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			data = make([]byte, length)
			if err := bitstream.DecodeAlphanumeric(buf, data); err != nil {
				return nil, err
			}
		case ModeBytes:
			length, err := buf.ReadBits(5)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break LOOP
				}
				return nil, err
			}
			data = make([]byte, length)
			if err := bitstream.DecodeBytes(buf, data); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("qrcode: unknown mode: " + strconv.Itoa(int(mode)))
		}
		if len(data) == 0 {
			continue
		}
		segments = append(segments, Segment{
			Mode: Mode(mode),
			Data: data,
		})
	}
	return &QRCode{
		Version:  4,
		Level:    level,
		Mask:     mask,
		Segments: segments,
	}, nil
}
