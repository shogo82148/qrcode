package qrcode

//go:generate go run genbch/main.go

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"

	bitmap "github.com/shogo82148/qrcode/bitmap"
	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func New(level Level, data []byte) (*QRCode, error) {
	if len(data) == 0 {
		return &QRCode{
			Version: 1,
			Level:   level,
			Mask:    MaskAuto,
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
	if isNumeric(data[0]) {
		states[0][modeNumeric] = state{
			cost:     (4+14)*6 + 20,
			lastMode: modeInit,
		}
	} else {
		states[0][modeNumeric] = state{
			cost:     inf,
			lastMode: modeInit,
		}
	}
	if isAlphanumeric(data[0]) {
		states[0][modeAlphanumeric] = state{
			cost:     (4+13)*6 + 33,
			lastMode: modeInit,
		}
	} else {
		states[0][modeAlphanumeric] = state{
			cost:     inf,
			lastMode: modeInit,
		}
	}
	states[0][modeBytes] = state{
		cost:     (4 + 16 + 8) * 6,
		lastMode: modeInit,
	}

	for i := 1; i < len(data); i++ {
		if isNumeric(data[i]) {
			// numeric -> numeric
			minCost := states[i-1][modeNumeric].cost + 20
			lastMode := modeNumeric

			// alphanumeric -> numeric
			cost := states[i-1][modeAlphanumeric].cost + (4+14)*6 + 20
			if cost < minCost {
				minCost = cost
				lastMode = modeAlphanumeric
			}

			// bytes -> numeric
			cost = states[i-1][modeBytes].cost + (4+14)*6 + 20
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

		if isAlphanumeric(data[i]) {
			// numeric -> alphanumeric
			minCost := states[i-1][modeNumeric].cost + (4+13)*6 + 33
			lastMode := modeNumeric

			// alphanumeric -> numeric
			cost := states[i-1][modeAlphanumeric].cost + 33
			if cost < minCost {
				minCost = cost
				lastMode = modeAlphanumeric
			}

			// bytes -> numeric
			cost = states[i-1][modeBytes].cost + (4+13)*6 + 33
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
		minCost := states[i-1][modeNumeric].cost + (4+16+8)*6
		lastMode := modeNumeric

		// alphanumeric -> bytes
		cost := states[i-1][modeAlphanumeric].cost + (4+16+8)*6
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

	return &QRCode{
		Level:    level,
		Mask:     MaskAuto,
		Segments: segments,
	}, nil
}

func isNumeric(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlphanumeric(ch byte) bool {
	return alphabets[ch] >= 0
}

const timingPatternOffset = 6

func skipTimingPattern(n int) int {
	if n < timingPatternOffset {
		return n
	}
	return n + 1
}

type EncodeOptions interface {
	apply(opts *encodeOptions)
}

type encodeOptions struct {
	QuiteZone  int
	ModuleSize float64
}

type withModuleSize float64

func (opt withModuleSize) apply(opts *encodeOptions) {
	opts.ModuleSize = float64(opt)
}

func WithModuleSize(size float64) EncodeOptions {
	return withModuleSize(size)
}

type withQuiteZone int

func (opt withQuiteZone) apply(opts *encodeOptions) {
	opts.QuiteZone = int(opt)
}

func WithQuiteZone(n int) EncodeOptions {
	return withQuiteZone(n)
}

func (qr *QRCode) Encode(opts ...EncodeOptions) (image.Image, error) {
	myopts := encodeOptions{
		QuiteZone:  4,
		ModuleSize: 1,
	}
	for _, o := range opts {
		o.apply(&myopts)
	}

	binimg, err := qr.EncodeToBitmap()
	if err != nil {
		return nil, err
	}

	w := binimg.Bounds().Dx() + myopts.QuiteZone*2
	W := int(math.Ceil(float64(w) * myopts.ModuleSize))
	scale := myopts.ModuleSize
	dX := float64(myopts.QuiteZone) * scale
	dY := float64(myopts.QuiteZone) * scale

	// create new paletted image
	palette := color.Palette{
		color.White, color.Black,
	}
	white := uint8(0)
	black := uint8(1)
	bounds := image.Rect(0, 0, W, W)
	img := image.NewPaletted(bounds, palette)

	// convert bitmap to image
	for Y := bounds.Min.Y; Y < bounds.Max.Y; Y++ {
		y := int(math.Floor((float64(Y) - dY) / scale))
		for X := bounds.Min.X; X < bounds.Max.X; X++ {
			x := int(math.Floor((float64(X) - dX) / scale))
			c := binimg.BinaryAt(x, y)
			if c {
				img.SetColorIndex(X, Y, black)
			} else {
				img.SetColorIndex(X, Y, white)
			}
		}
	}
	return img, nil
}

// EncodeToBitmap encodes QR Code into bitmap image.
func (qr *QRCode) EncodeToBitmap() (*bitmap.Image, error) {
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

	// version
	if qr.Version >= 7 {
		version := encodedVersion[qr.Version]
		for i := 0; i < 18; i++ {
			img.SetBinary(i/3, w-10+i%3, (version>>i)&1 != 0)
			img.SetBinary(w-10+i%3, i/3, (version>>i)&1 != 0)
		}
	}

	// mask
	mask := qr.Mask
	if mask == MaskAuto {
		var tmp internalbitmap.Image
		var minPoint int
		mask = Mask0
		for i := Mask0; i < maskMax; i++ {
			tmp.Mask(img, used, maskList[i])
			format := encodedFormat[int(qr.Level)<<3+int(i)]
			for i := 0; i < 8; i++ {
				tmp.SetBinary(8, skipTimingPattern(i), (format>>i)&1 != 0)
				tmp.SetBinary(skipTimingPattern(i), 8, (format>>(14-i))&1 != 0)

				tmp.SetBinary(w-i, 8, (format>>i)&1 != 0)
				tmp.SetBinary(8, w-i, (format>>(14-i))&1 != 0)
			}
			tmp.SetBinary(8, w-7, internalbitmap.Black)

			point := tmp.Point()
			if point < minPoint {
				minPoint = point
				mask = i
			}
		}
	}

	// format
	format := encodedFormat[int(qr.Level)<<3+int(mask)]
	for i := 0; i < 8; i++ {
		img.SetBinary(8, skipTimingPattern(i), (format>>i)&1 != 0)
		img.SetBinary(skipTimingPattern(i), 8, (format>>(14-i))&1 != 0)

		img.SetBinary(w-i, 8, (format>>i)&1 != 0)
		img.SetBinary(8, w-i, (format>>(14-i))&1 != 0)
	}
	img.SetBinary(8, w-7, internalbitmap.Black)
	img.Mask(img, used, maskList[mask])

	return img.Export(), nil
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
	case ModeNumeric:
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
	buf.WriteBitsLSB(uint64(ModeNumeric), 4)

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

var bitToAlphanumeric = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:")
var alphabets [256]int

func init() {
	for i := range alphabets {
		alphabets[i] = -1
	}
	for i, ch := range bitToAlphanumeric {
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
