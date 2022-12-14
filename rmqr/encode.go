package rmqr

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/shogo82148/qrcode/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func (qr *QRCode) EncodeToBitmap() (*bitmap.Image, error) {
	var buf bitstream.Buffer
	if err := qr.encodeToBits(&buf); err != nil {
		return nil, err
	}

	h, w := 7-1, 43-1 // TODO: calculate from version.
	used := usedList[qr.Version]
	img := baseList[qr.Version].Clone()

	dy := -1
	x, y := w-1, h-4
	for {
		if !used.BinaryAt(x, y) {
			bit, err := buf.ReadBit()
			if err != nil {
				break
			}
			img.SetBinary(x, y, bit != 0)
		}
		x--
		if x < 1 {
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
		if y < 1 || y > h-1 {
			dy *= -1
			x, y = x-2, y+dy
		}
		if x < 1 {
			break
		}
	}

	encodedFormat := encodedVersion[uint(qr.Version)|uint(qr.Level)<<5]
	for i := 0; i < 18; i++ {
		img.SetBinary(8+i/5, 1+i%5, ((encodedFormat^0b011111101010110010)>>i)&1 != 0)
	}

	img.Mask(img, used, precomputedMask)

	var hoge bytes.Buffer
	img.EncodePBM(&hoge)
	os.WriteFile("rmqr.pbm", hoge.Bytes(), 0o644)

	return img.Export(), nil
}

type block struct {
	data       []byte
	correction []byte
	maxError   int
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
	capacity := capacityTable[qr.Version][qr.Level]
	bitsLength := capacity.BitLength
	for _, s := range qr.Segments {
		if err := s.encode(bitsLength, buf); err != nil {
			return err
		}
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

func (s *Segment) encode(bitLength [5]int, buf *bitstream.Buffer) error {
	switch s.Mode {
	case ModeNumeric:
		return s.encodeNumber(bitLength[ModeNumeric], buf)
	default:
		return errors.New("qrcode: unknown mode")
	}
}

func (s *Segment) encodeNumber(n int, buf *bitstream.Buffer) error {
	if len(s.Data) >= 1<<n {
		return fmt.Errorf("rmqr: data is too long: %d", len(s.Data))
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeNumeric), 3)

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	return bitstream.EncodeNumeric(buf, s.Data)
}
