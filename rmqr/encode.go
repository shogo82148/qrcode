package rmqr

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/shogo82148/qrcode/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func New(level Level, data []byte) (*QRCode, error) {
	if len(data) == 0 {
		return &QRCode{
			Version: 0,
			Level:   level,
		}, nil
	}

	const inf = math.MaxInt - 1<<18 // 1<<18 is for avoiding overflow
	const (
		modeNumeric byte = iota
		modeAlphanumeric
		modeBytes
		modeInit
	)
	type state struct {
		cost     int // = bit length * 6
		lastMode byte
	}
	states := make([][3]state, len(data))
	if bitstream.IsNumeric(data[0]) {
		states[0][modeNumeric] = state{
			cost:     (3+9)*6 + 20,
			lastMode: modeInit,
		}
	} else {
		states[0][modeNumeric] = state{
			cost:     inf,
			lastMode: modeInit,
		}
	}
	if bitstream.IsAlphanumeric(data[0]) {
		states[0][modeAlphanumeric] = state{
			cost:     (3+8)*6 + 33,
			lastMode: modeInit,
		}
	} else {
		states[0][modeAlphanumeric] = state{
			cost:     inf,
			lastMode: modeInit,
		}
	}
	states[0][modeBytes] = state{
		cost:     (3 + 8 + 8) * 6,
		lastMode: modeInit,
	}

	for i := 1; i < len(data); i++ {
		if bitstream.IsNumeric(data[i]) {
			// numeric -> numeric
			minCost := states[i-1][modeNumeric].cost + 20
			lastMode := modeNumeric

			// alphanumeric -> numeric
			cost := states[i-1][modeAlphanumeric].cost + (3+9)*6 + 20
			if cost < minCost {
				minCost = cost
				lastMode = modeAlphanumeric
			}

			// bytes -> numeric
			cost = states[i-1][modeBytes].cost + (3+9)*6 + 20
			if cost < minCost {
				minCost = cost
				lastMode = modeBytes
			}

			states[i][modeNumeric] = state{
				cost:     minCost,
				lastMode: lastMode,
			}
		} else {
			states[i][modeNumeric] = state{
				cost:     inf,
				lastMode: modeInit,
			}
		}

		if bitstream.IsAlphanumeric(data[i]) {
			// numeric -> alphanumeric
			minCost := states[i-1][modeNumeric].cost + (3+8)*6 + 33
			lastMode := modeNumeric

			// alphanumeric -> numeric
			cost := states[i-1][modeAlphanumeric].cost + 33
			if cost < minCost {
				minCost = cost
				lastMode = modeAlphanumeric
			}

			// bytes -> numeric
			cost = states[i-1][modeBytes].cost + (3+8)*6 + 33
			if cost < minCost {
				minCost = cost
				lastMode = modeBytes
			}

			states[i][modeAlphanumeric] = state{
				cost:     minCost,
				lastMode: lastMode,
			}
		} else {
			states[i][modeAlphanumeric] = state{
				cost:     inf,
				lastMode: modeInit,
			}
		}

		// numeric -> bytes
		minCost := states[i-1][modeNumeric].cost + (3+8+8)*6
		lastMode := modeNumeric

		// alphanumeric -> bytes
		cost := states[i-1][modeAlphanumeric].cost + (3+8+8)*6
		if cost < minCost {
			minCost = cost
			lastMode = modeAlphanumeric
		}

		// bytes -> bytes
		cost = states[i-1][modeBytes].cost + 8*6
		if cost < minCost {
			minCost = cost
			lastMode = modeBytes
		}
		states[i][modeBytes] = state{
			cost:     minCost,
			lastMode: lastMode,
		}
	}

	best := make([]byte, len(data))
	minCost := states[len(data)-1][modeNumeric].cost
	bestMode := modeNumeric
	if cost := states[len(data)-1][modeAlphanumeric].cost; cost < minCost {
		minCost = cost
		bestMode = modeAlphanumeric
	}
	if cost := states[len(data)-1][modeBytes].cost; cost < minCost {
		bestMode = modeBytes
	}
	best[len(data)-1] = bestMode
	for i := len(data) - 1; i > 0; i-- {
		bestMode = states[i][bestMode].lastMode
		best[i-1] = bestMode
	}

	modeList := [...]Mode{ModeNumeric, ModeAlphanumeric, ModeBytes}
	segments := []Segment{
		{
			Mode: modeList[best[0]],
			Data: []byte{data[0]},
		},
	}
	for i := 1; i < len(data); i++ {
		mode := modeList[best[i]]
		if segments[len(segments)-1].Mode == mode {
			segments[len(segments)-1].Data = append(segments[len(segments)-1].Data, data[i])
		} else {
			segments = append(segments, Segment{
				Mode: mode,
				Data: []byte{data[i]},
			})
		}
	}

	// TODO: implement calcVersion
	// version := calcVersion(level, segments)
	// if version == 0 {
	// 	return nil, errors.New("qrcode: data too large")
	// }

	return &QRCode{
		// Version:  version,
		Level:    level,
		Segments: segments,
	}, nil
}

func (qr *QRCode) EncodeToBitmap() (*bitmap.Image, error) {
	var buf bitstream.Buffer
	if err := qr.encodeToBits(&buf); err != nil {
		return nil, err
	}

	used := usedList[qr.Version]
	img := baseList[qr.Version].Clone()
	bounds := img.Rect
	w, h := bounds.Dx()-1, bounds.Dy()-1

	dy := -1
	x, y := w-1, h-5
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

	// fill format information
	encodedFormat := encodedVersion[uint(qr.Version)|uint(qr.Level)<<5]
	for i := 0; i < 18; i++ {
		img.SetBinary(8+i/5, 1+i%5, ((encodedFormat^0b011111101010110010)>>i)&1 != 0)
	}
	for i := 0; i < 15; i++ {
		img.SetBinary(w-7+i/5, h-5+i%5, ((encodedFormat^0b100000101001111011)>>i)&1 != 0)
	}
	img.SetBinary(w-4, h-5, ((encodedFormat^0b100000101001111011)>>15)&1 != 0)
	img.SetBinary(w-3, h-5, ((encodedFormat^0b100000101001111011)>>16)&1 != 0)
	img.SetBinary(w-2, h-5, ((encodedFormat^0b100000101001111011)>>17)&1 != 0)

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
