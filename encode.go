package qrcode

//go:generate go run genbch/main.go

import (
	"errors"
	"fmt"
	"image"

	"github.com/shogo82148/qrcode/internal/binimage"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func New(data []byte) (*QRCode, error) {
	return &QRCode{}, nil
}

const timingPatternOffset = 6

func skipTimingPattern(n int) int {
	if n < timingPatternOffset {
		return n
	}
	return n + 1
}

func (qr *QRCode) Encode() (image.Image, error) {
	var buf bitstream.Buffer
	if err := qr.encodeToBits(&buf); err != nil {
		return nil, err
	}

	w := 16 + 4*int(qr.Version)
	img := baseList[qr.Version].Clone()
	used := usedList[qr.Version]

	dy := -1
	x, y := w, w
	for {
		if x == timingPatternOffset {
			// skip timing pattern
			x--
			continue
		}
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

	// format
	format := encodedFormat[int(qr.Level)<<3+int(qr.Mask)]
	for i := 0; i < 8; i++ {
		img.SetBinary(8, skipTimingPattern(i), (format>>i)&1 != 0)
		img.SetBinary(skipTimingPattern(i), 8, (format>>(14-i))&1 != 0)

		img.SetBinary(w-i, 8, (format>>i)&1 != 0)
		img.SetBinary(8, w-i, (format>>(14-i))&1 != 0)
	}
	img.SetBinary(8, w-7, binimage.Black)

	// version
	if qr.Version >= 7 {
		version := encodedVersion[qr.Version]
		for i := 0; i < 18; i++ {
			img.SetBinary(i/3, w-10+i%3, (version>>i)&1 != 0)
			img.SetBinary(w-10+i%3, i/3, (version>>i)&1 != 0)
		}
	}

	// mask
	var f func(i, j int) int
	switch qr.Mask {
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
			img.XorBinary(j, i, !used.BinaryAt(j, i) && f(i, j) == 0)
		}
	}

	return img, nil
}

type block struct {
	data       []byte
	correction []byte
}

func (qr *QRCode) encodeToBits(ret *bitstream.Buffer) error {
	var buf bitstream.Buffer
	if err := qr.encodeSegments(&buf); err != nil {
		return err
	}

	// split to block and calculate error correction code.
	capacity := capacityTable[qr.Version][qr.Level]
	data := buf.Bytes()
	blocks := []block{}
	for _, blockCapacity := range capacity.Blocks {
		for i := 0; i < blockCapacity.Num; i++ {
			n := blockCapacity.Total - blockCapacity.Data
			rs := reedsolomon.New(n)
			rs.Write(data[:blockCapacity.Data])
			correction := rs.Sum(make([]byte, 0, n))
			blocks = append(blocks, block{
				data:       data[:blockCapacity.Data],
				correction: correction,
			})
			data = data[blockCapacity.Data:]
		}
	}

	// Interleave the blocks.
	for i := 0; ; i++ {
		var wrote int
		for _, b := range blocks {
			if i < len(b.data) {
				ret.WriteBitsLSB(uint64(b.data[i]), 8)
				wrote++
			}
		}
		if wrote == 0 {
			break
		}
	}
	for i := 0; ; i++ {
		var wrote int
		for _, b := range blocks {
			if i < len(b.correction) {
				ret.WriteBitsLSB(uint64(b.correction[i]), 8)
				wrote++
			}
		}
		if wrote == 0 {
			break
		}
	}
	return nil
}

func (qr *QRCode) encodeSegments(buf *bitstream.Buffer) error {
	for _, s := range qr.Segments {
		if err := s.encode(qr.Version, buf); err != nil {
			return err
		}
	}
	capacity := capacityTable[qr.Version][qr.Level]
	if buf.Len() > capacity.Data*8 {
		return errors.New("qrcode: data is too large")
	}
	l := buf.Len()
	buf.WriteBitsLSB(0x00, int(8-l%8))

	// add padding.
	for i := 0; buf.Len() < capacity.Data*8; i++ {
		if i%2 == 0 {
			buf.WriteBitsLSB(0b1110_1100, 8)
		} else {
			buf.WriteBitsLSB(0b0001_0001, 8)
		}
	}
	return nil
}

func (s *Segment) encode(version Version, buf *bitstream.Buffer) error {
	switch s.Mode {
	case ModeNumber:
		return s.encodeNumber(version, buf)
	case ModeAlphanumeric:
		return s.encodeAlphabet(version, buf)
	case ModeBytes:
		return s.encodeBytes(version, buf)
	default:
		return errors.New("qrcode: unknown mode")
	}
}

func (s *Segment) encodeNumber(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch {
	case version <= 0 || version > 40:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	case version < 10:
		n = 10
	case version < 27:
		n = 12
	default:
		n = 14
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	for _, ch := range data {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("qrcode: invalid character in number mode: %02x", ch)
		}
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeNumber), 4)

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	for i := 0; i+2 < len(data); i += 3 {
		n1 := int(data[i+0] - '0')
		n2 := int(data[i+1] - '0')
		n3 := int(data[i+2] - '0')
		bits := n1*100 + n2*10 + n3
		buf.WriteBitsLSB(uint64(bits), 10)
	}

	switch len(data) % 3 {
	case 1:
		bits := data[len(data)-1] - '0'
		buf.WriteBitsLSB(uint64(bits), 4)
	case 2:
		n1 := int(data[len(data)-2] - '0')
		n2 := int(data[len(data)-1] - '0')
		bits := n1*10 + n2
		buf.WriteBitsLSB(uint64(bits), 7)
	}
	return nil
}

var alphabets [256]int

func init() {
	for i := range alphabets {
		alphabets[i] = -1
	}
	for i, ch := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:" {
		alphabets[ch] = i
	}
}

func (s *Segment) encodeAlphabet(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch {
	case version <= 0 || version > 40:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	case version < 10:
		n = 9
	case version < 27:
		n = 11
	default:
		n = 13
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	for _, ch := range data {
		if alphabets[ch] < 0 {
			return fmt.Errorf("qrcode: invalid character in alphabet mode: %02x", ch)
		}
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeAlphanumeric), 4)

	// data length
	buf.WriteBitsLSB(uint64(len(data)), n)

	// data
	for i := 0; i+1 < len(data); i += 2 {
		n1 := alphabets[data[i+0]]
		n2 := alphabets[data[i+1]]
		bits := n1*45 + n2
		buf.WriteBitsLSB(uint64(bits), 11)
	}

	if len(data)%2 != 0 {
		bits := alphabets[data[len(data)-1]]
		buf.WriteBitsLSB(uint64(bits), 6)
	}
	return nil
}

func (s *Segment) encodeBytes(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch {
	case version <= 0 || version > 40:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	case version < 10:
		n = 8
	default:
		n = 16
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeBytes), 4)

	// data length
	buf.WriteBitsLSB(uint64(len(data)), n)

	// data
	for _, bits := range data {
		buf.WriteBitsLSB(uint64(bits), 8)
	}
	return nil
}
