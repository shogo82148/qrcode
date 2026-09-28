package tqr

import (
	"fmt"
	"math/bits"

	"github.com/shogo82148/go-imaging/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

// https://github.com/kikuchan/libqrean/blob/b1dd0bdc72c3644593c0a503bd5fe24588f10d76/src/code_tqr.c#L77-L241
var dataPos = [...]struct{ x, y uint8 }{
	{16, 18}, // block #1
	{17, 17}, // block #1
	{16, 17}, // block #1
	{18, 16}, // block #1
	{17, 16}, // block #1
	{16, 16}, // block #1
	{18, 15}, // block #1
	{17, 15}, // block #1
	{16, 15}, // block #1
	{18, 14}, // block #1
	{17, 14}, // block #2
	{16, 14}, // block #2
	{18, 13}, // block #2
	{17, 13}, // block #2
	{16, 13}, // block #2
	{18, 12}, // block #2
	{17, 12}, // block #2
	{16, 12}, // block #2
	{18, 11}, // block #2
	{17, 11}, // block #2
	{16, 11}, // block #3
	{18, 10}, // block #3
	{17, 10}, // block #3
	{16, 10}, // block #3
	{18, 9},  // block #3
	{17, 9},  // block #3
	{16, 9},  // block #3
	{18, 8},  // block #3
	{17, 8},  // block #3
	{16, 8},  // block #3
	{15, 9},  // block #4
	{14, 9},  // block #4
	{13, 9},  // block #4
	{12, 9},  // block #4
	{11, 9},  // block #4
	{15, 8},  // block #4
	{14, 8},  // block #4
	{13, 8},  // block #4
	{12, 8},  // block #4
	{11, 8},  // block #4
	{15, 11}, // block #5
	{14, 11}, // block #5
	{13, 11}, // block #5
	{12, 11}, // block #5
	{11, 11}, // block #5
	{15, 10}, // block #5
	{14, 10}, // block #5
	{13, 10}, // block #5
	{12, 10}, // block #5
	{11, 10}, // block #5
	{14, 15}, // block #6
	{13, 15}, // block #6
	{15, 14}, // block #6
	{14, 14}, // block #6
	{13, 14}, // block #6
	{15, 13}, // block #6
	{14, 13}, // block #6
	{13, 13}, // block #6
	{15, 12}, // block #6
	{14, 12}, // block #6
	{11, 15}, // block #7
	{10, 15}, // block #7
	{12, 14}, // block #7
	{11, 14}, // block #7
	{10, 14}, // block #7
	{12, 13}, // block #7
	{11, 13}, // block #7
	{13, 12}, // block #7
	{12, 12}, // block #7
	{11, 12}, // block #7
	{15, 18}, // block #8
	{14, 18}, // block #8
	{13, 18}, // block #8
	{15, 17}, // block #8
	{14, 17}, // block #8
	{13, 17}, // block #8
	{15, 16}, // block #8
	{14, 16}, // block #8
	{13, 16}, // block #8
	{15, 15}, // block #8
	{12, 18}, // block #9
	{11, 18}, // block #9
	{10, 18}, // block #9
	{12, 17}, // block #9
	{11, 17}, // block #9
	{10, 17}, // block #9
	{12, 16}, // block #9
	{11, 16}, // block #9
	{10, 16}, // block #9
	{12, 15}, // block #9
	{9, 18},  // block #10
	{8, 18},  // block #10
	{9, 17},  // block #10
	{8, 17},  // block #10
	{9, 16},  // block #10
	{8, 16},  // block #10
	{9, 15},  // block #10
	{8, 15},  // block #10
	{9, 14},  // block #10
	{8, 14},  // block #10
	{10, 13}, // block #11
	{9, 13},  // block #11
	{8, 13},  // block #11
	{10, 12}, // block #11
	{9, 12},  // block #11
	{8, 12},  // block #11
	{10, 11}, // block #11
	{9, 11},  // block #11
	{8, 11},  // block #11
	{10, 10}, // block #11
	{9, 10},  // block #12
	{8, 10},  // block #12
	{10, 9},  // block #12
	{9, 9},   // block #12
	{8, 9},   // block #12
	{10, 8},  // block #12
	{9, 8},   // block #12
	{8, 8},   // block #12
	{7, 8},   // block #12
	{8, 7},   // block #12
	{10, 7},  // block #13
	{9, 7},   // block #13
	{10, 5},  // block #13
	{9, 5},   // block #13
	{8, 5},   // block #13
	{10, 4},  // block #13
	{9, 4},   // block #13
	{8, 4},   // block #13
	{10, 3},  // block #13
	{9, 3},   // block #13
	{8, 3},   // block #14
	{10, 2},  // block #14
	{9, 2},   // block #14
	{8, 2},   // block #14
	{10, 1},  // block #14
	{9, 1},   // block #14
	{8, 1},   // block #14
	{10, 0},  // block #14
	{9, 0},   // block #14
	{8, 0},   // block #14
	{7, 10},  // block #15
	{5, 10},  // block #15
	{4, 10},  // block #15
	{7, 9},   // block #15
	{5, 9},   // block #15
	{4, 9},   // block #15
	{3, 9},   // block #15
	{5, 8},   // block #15
	{4, 8},   // block #15
	{3, 8},   // block #15
	{3, 10},  // block #16
	{2, 10},  // block #16
	{1, 10},  // block #16
	{0, 10},  // block #16
	{2, 9},   // block #16
	{1, 9},   // block #16
	{0, 9},   // block #16
	{2, 8},   // block #16
	{1, 8},   // block #16
	{0, 8},   // block #16
}

func New(data []byte, opts ...EncodeOptions) (*QRCode, error) {
	qr := &QRCode{
		Data: data,
	}
	if err := qr.validate(); err != nil {
		return nil, err
	}

	return qr, nil
}

func (qr *QRCode) validate() error {
	if len(qr.Data) != 12 {
		return fmt.Errorf("tqr: data length must be 12, got %d", len(qr.Data))
	}
	for _, b := range qr.Data {
		if b < '0' || b > '9' {
			return fmt.Errorf("tqr: data must be numeric, got %q", b)
		}
	}
	return nil
}

type EncodeOptions func(opts *encodeOptions)

type encodeOptions struct {
}

func (qr *QRCode) EncodeToBitmap(opts ...EncodeOptions) (*bitmap.Image, error) {
	if err := qr.validate(); err != nil {
		return nil, err
	}

	// encode the numeric data into a bitstream
	var buf bitstream.Buffer
	if err := bitstream.EncodeNumeric(&buf, qr.Data); err != nil {
		return nil, err
	}

	// generate error correction code using Reed-Solomon
	rs := reedsolomon.New(11)
	rs.Write(buf.Bytes())
	correction := rs.Sum(make([]byte, 0, 11))

	// add parity bits
	var buf2 bitstream.Buffer
	for _, b := range buf.Bytes() {
		_ = buf2.WriteBitsLSB(uint64(b), 8)
		p := 1 + bits.OnesCount8(b)
		_ = buf2.WriteBitsLSB(uint64(p%2), 1)
		_ = buf2.WriteBitsLSB(uint64(p/2), 1)
	}
	for _, b := range correction {
		_ = buf2.WriteBitsLSB(uint64(b), 8)
		p := 1 + bits.OnesCount8(b)
		_ = buf2.WriteBitsLSB(uint64(p%2), 1)
		_ = buf2.WriteBitsLSB(uint64(p/2), 1)
	}

	// fill the QR code matrix
	img := base.Clone()
	for _, pos := range dataPos {
		bit, err := buf2.ReadBit()
		if err != nil {
			return nil, err
		}
		img.SetBinary(int(pos.x)+2, int(pos.y)+2, bit != 0)
	}

	// mask
	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			if used.BinaryAt(x+2, y+2) {
				continue
			}
			if (y/2+x/3)%2 == 0 {
				img.SetBinary(x+2, y+2, !img.BinaryAt(x+2, y+2))
			}
		}
	}

	return img.Export(), nil
}
