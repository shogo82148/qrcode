package microqr

import (
	"errors"
	"fmt"

	bitmap "github.com/shogo82148/qrcode/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

// EncodeToBitmap encodes QR Code into bitmap image.
func (qr *QRCode) EncodeToBitmap() (*bitmap.Image, error) {
	if qr.Version < 1 || qr.Version > 4 {
		return nil, fmt.Errorf("microqr: invalid version: %d", qr.Version)
	}
	if qr.Level < 0 || qr.Level >= 4 {
		return nil, fmt.Errorf("microqr: invalid level: %d", qr.Level)
	}
	format := formatTable[qr.Version][qr.Level]
	if format < 0 {
		return nil, fmt.Errorf("microqr: invalid version-level pair: %d-%s", qr.Version, qr.Level)
	}

	var buf bitstream.Buffer
	if err := qr.encodeSegments(&buf); err != nil {
		return nil, err
	}

	w := 8 + 2*int(qr.Version)
	img := baseList[qr.Version].Clone()
	used := usedList[qr.Version]

	dy := -1
	x, y := w, w
	for {
		if !used.BinaryAt(x, y) {
			bit, err := buf.ReadBit()
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
			bit, err := buf.ReadBit()
			if err != nil {
				break
			}
			img.SetBinary(x, y, bit != 0)
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

	mask := qr.Mask
	encoded := encodedFormat[(format<<2)|int(mask)]
	for i := 0; i < 8; i++ {
		img.SetBinary(8, i+1, (encoded>>i)&1 != 0)
		img.SetBinary(i+1, 8, (encoded>>(14-i))&1 != 0)
	}

	img.Mask(img, used, maskList[mask])

	return img.Export(), nil
}

func (qr *QRCode) encodeSegments(buf *bitstream.Buffer) error {
	for _, s := range qr.Segments {
		if err := s.encode(qr.Version, buf); err != nil {
			return err
		}
	}

	// terminate pattern
	switch qr.Version {
	case 1:
		buf.WriteBitsLSB(0, 3)
	case 2:
		buf.WriteBitsLSB(0, 5)
	case 3:
		buf.WriteBitsLSB(0, 7)
	case 4:
		buf.WriteBitsLSB(0, 8)
	}

	l := buf.Len()
	buf.WriteBitsLSB(0x00, int(8-l%8))
	capacity := capacityTable[qr.Version][qr.Level]

	// add padding.
	for i := 0; buf.Len() < capacity.Data*8; i++ {
		if i%2 == 0 {
			buf.WriteBitsLSB(0b1110_1100, 8)
		} else {
			buf.WriteBitsLSB(0b0001_0001, 8)
		}
	}

	n := capacity.Correction
	rs := reedsolomon.New(n)
	rs.Write(buf.Bytes())
	correction := rs.Sum(make([]byte, 0, n))
	for _, b := range correction {
		buf.WriteBitsLSB(uint64(b), 8)
	}
	return nil
}

func (s *Segment) encode(version Version, buf *bitstream.Buffer) error {
	switch s.Mode {
	case ModeNumeric:
		return s.encodeNumeric(version, buf)
	case ModeAlphanumeric:
		return s.encodeAlphanumeric(version, buf)
	case ModeBytes:
		return s.encodeBytes(version, buf)
	default:
		return errors.New("qrcode: unknown mode")
	}
}

func (s *Segment) encodeNumeric(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch version {
	case 1:
		n = 3
	case 2:
		n = 4
	case 3:
		n = 5
	case 4:
		n = 6
	default:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	// mode
	switch version {
	case 2:
		buf.WriteBitsLSB(uint64(ModeNumeric), 1)
	case 3:
		buf.WriteBitsLSB(uint64(ModeNumeric), 2)
	case 4:
		buf.WriteBitsLSB(uint64(ModeNumeric), 3)
	}

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	return bitstream.EncodeNumeric(buf, data)
}

func (s *Segment) encodeAlphanumeric(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch version {
	case 2:
		n = 3
	case 3:
		n = 4
	case 4:
		n = 5
	default:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	// mode
	switch version {
	case 2:
		buf.WriteBitsLSB(uint64(ModeAlphanumeric), 1)
	case 3:
		buf.WriteBitsLSB(uint64(ModeAlphanumeric), 2)
	case 4:
		buf.WriteBitsLSB(uint64(ModeAlphanumeric), 3)
	}

	// data length
	buf.WriteBitsLSB(uint64(len(data)), n)

	// data
	return bitstream.EncodeAlphanumeric(buf, data)
}

func (s *Segment) encodeBytes(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch version {
	case 3:
		n = 4
	case 4:
		n = 5
	default:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	// mode
	switch version {
	case 3:
		buf.WriteBitsLSB(uint64(ModeBytes), 2)
	case 4:
		buf.WriteBitsLSB(uint64(ModeBytes), 3)
	}

	// data length
	buf.WriteBitsLSB(uint64(len(data)), n)

	// data
	for _, bits := range data {
		buf.WriteBitsLSB(uint64(bits), 8)
	}
	return nil
}
