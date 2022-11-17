package qrcode

//go:generate go run genbch/main.go

import (
	"errors"
	"fmt"
	"image"
	"strconv"

	"github.com/shogo82148/qrcode/internal/bitstream"
	binimage "github.com/shogo82148/qrcode/internal/image"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

type QRCode struct {
	Version  Version
	Level    Level
	Mask     Mask
	Segments []Segment
}

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
	qr.encodeToBits(&buf)

	w := 16 + 4*int(qr.Version)
	img := baseList[qr.Version]
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

	for i := 0; i <= w; i++ {
		for j := 0; j <= w; j++ {
			img.XorBinary(i, j, !used.BinaryAt(i, j) && i%3 == 0)
		}
	}

	return img, nil
}

type block struct {
	data       []byte
	correction []byte
}

func (qr *QRCode) encodeToBits(ret *bitstream.Buffer) {
	var buf bitstream.Buffer
	for _, s := range qr.Segments {
		s.encode(qr.Version, &buf)
	}
	l := buf.Len()
	buf.WriteBitsLSB(0x00, int(8-l%8))

	// add padding.
	capacity := capacityTable[qr.Version][qr.Level]
	for i := 0; buf.Len() < capacity.Data*8; i++ {
		if i%2 == 0 {
			buf.WriteBitsLSB(0b1110_1100, 8)
		} else {
			buf.WriteBitsLSB(0b0001_0001, 8)
		}
	}

	// split to block and calculate error correction code.
	data := buf.Bytes()
	blocks := []block{}
	for _, blockCapacity := range capacity.Blocks {
		for i := 0; i < blockCapacity.Num; i++ {
			n := blockCapacity.Total - blockCapacity.Data
			rs := reedsolomon.New(n)
			rs.Write(data[:capacity.Data])
			correction := rs.Sum(make([]byte, 0, n))
			blocks = append(blocks, block{
				data:       data[:capacity.Data],
				correction: correction,
			})
		}
	}

	// assemble
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
}

type Version int

type Level int

const (
	LevelL Level = 0b01
	LevelM Level = 0b00
	LevelQ Level = 0b11
	LevelH Level = 0b10
)

func (lv Level) String() string {
	switch lv {
	case LevelL:
		return "L"
	case LevelM:
		return "M"
	case LevelQ:
		return "Q"
	case LevelH:
		return "H"
	}
	return "invalid(" + strconv.Itoa(int(lv)) + ")"
}

type Mask int

type Mode uint8

const (
	// ModeECI is ECI(Extended Channel Interpretation) mode.
	ModeECI Mode = 0b0111

	// ModeNumber is number mode.
	// The Data must be ascii characters [0-9].
	ModeNumber Mode = 0b0001

	// ModeNumber is alphabet and number mode.
	// The Data must be ascii characters [0-9A-Z $%*+\-./:].
	ModeAlphabet Mode = 0b0010

	// ModeBytes is 8-bit bytes mode.
	// The Data can include any bytes.
	ModeBytes Mode = 0b0100

	// ModeKanji is Japanese Kanji mode.
	ModeKanji Mode = 0b1000

	// ModeConnected is connected structure mode.
	ModeConnected Mode = 0b0011

	ModeFNC1_1 Mode = 0b0101
	ModeFNC1_2 Mode = 0b1001

	ModeTerminated Mode = 0b0000
)

type Segment struct {
	Mode Mode
	Data []byte
}

func (s *Segment) encode(version Version, buf *bitstream.Buffer) error {
	switch s.Mode {
	case ModeNumber:
		return s.encodeNumber(version, buf)
	case ModeAlphabet:
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
	buf.WriteBitsLSB(uint64(ModeAlphabet), 4)

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
