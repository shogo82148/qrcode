package microqr

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
	if bitstream.IsNumeric(data[0]) {
		states[0][modeNumeric] = state{
			cost:     (4+6)*6 + 20,
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
			cost:     (4+5)*6 + 33,
			lastMode: modeInit,
		}
	} else {
		states[0][modeAlphanumeric] = state{
			cost:     inf,
			lastMode: modeInit,
		}
	}
	states[0][modeBytes] = state{
		cost:     (4 + 5 + 8) * 6,
		lastMode: modeInit,
	}

	for i := 1; i < len(data); i++ {
		if bitstream.IsNumeric(data[i]) {
			// numeric -> numeric
			minCost := states[i-1][modeNumeric].cost + 20
			lastMode := modeNumeric

			// alphanumeric -> numeric
			cost := states[i-1][modeAlphanumeric].cost + (4+6)*6 + 20
			if cost < minCost {
				minCost = cost
				lastMode = modeAlphanumeric
			}

			// bytes -> numeric
			cost = states[i-1][modeBytes].cost + (4+6)*6 + 20
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
			minCost := states[i-1][modeNumeric].cost + (4+5)*6 + 33
			lastMode := modeNumeric

			// alphanumeric -> numeric
			cost := states[i-1][modeAlphanumeric].cost + 33
			if cost < minCost {
				minCost = cost
				lastMode = modeAlphanumeric
			}

			// bytes -> numeric
			cost = states[i-1][modeBytes].cost + (4+5)*6 + 33
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
		minCost := states[i-1][modeNumeric].cost + (4+5+8)*6
		lastMode := modeNumeric

		// alphanumeric -> bytes
		cost := states[i-1][modeAlphanumeric].cost + (4+5+8)*6
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

	version := calcVersion(level, segments)
	if version == 0 {
		return nil, errors.New("microqr: data too large")
	}

	return &QRCode{
		Version:  version,
		Level:    level,
		Mask:     MaskAuto,
		Segments: segments,
	}, nil
}

func calcVersion(level Level, segments []Segment) Version {
LOOP:
	for version := Version(1); version <= 4; version++ {
		capacity := capacityTable[version][level].DataBits
		length := 0
		for _, s := range segments {
			l, ok := s.length(version)
			if !ok {
				continue LOOP
			}
			length += l
			if length > capacity {
				continue LOOP
			}
		}
		if length <= capacity {
			return version
		}
	}
	return 0
}

type EncodeOptions interface {
	apply(opts *encodeOptions)
}

type encodeOptions struct {
	QuiteZone  int
	ModuleSize float64
	Level      Level
}

func newEncodeOptions(opts ...EncodeOptions) encodeOptions {
	myopts := encodeOptions{
		QuiteZone:  4,
		ModuleSize: 1,
		Level:      LevelQ,
	}
	for _, o := range opts {
		o.apply(&myopts)
	}
	return myopts
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

type withLevel Level

func (opt withLevel) apply(opts *encodeOptions) {
	opts.Level = Level(opt)
}

func WithLevel(lv Level) EncodeOptions {
	return withLevel(lv)
}

func Encode(data []byte, opts ...EncodeOptions) (image.Image, error) {
	myopts := newEncodeOptions(opts...)

	qr, err := New(myopts.Level, data)
	if err != nil {
		return nil, err
	}
	return qr.Encode(opts...)
}

func (qr *QRCode) Encode(opts ...EncodeOptions) (image.Image, error) {
	myopts := encodeOptions{
		QuiteZone:  2,
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
	readBits := 0
	capacity := capacityTable[qr.Version][qr.Level]
	for {
		if !used.BinaryAt(x, y) {
			readBits++
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
			readBits++
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
		if readBits == capacity.DataBits {
			for readBits%8 != 0 {
				readBits++
				buf.ReadBit()
			}
		}
	}

	mask := qr.Mask
	if mask == MaskAuto {
		var tmp internalbitmap.Image
		var minPoint int
		mask = Mask0
		for i := Mask0; i < maskMax; i++ {
			tmp.Mask(img, used, maskList[i])
			point := tmp.Point()
			if point < minPoint {
				minPoint = point
				mask = i
			}
		}
	}
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

	capacity := capacityTable[qr.Version][qr.Level]
	if buf.Len() > capacity.DataBits {
		return errors.New("qrcode: data too large")
	}

	// terminate pattern
	left := capacity.DataBits - buf.Len()
	var terminate int
	switch qr.Version {
	case 1:
		terminate = 3
	case 2:
		terminate = 5
	case 3:
		terminate = 7
	case 4:
		terminate = 8
	}
	if terminate < left {
		left = terminate
	}
	buf.WriteBitsLSB(0x00, left)

	// add padding.
	if buf.Len() < capacity.DataBits {
		buf.WriteBitsLSB(0x00, int(8-buf.Len()%8))
		for i := 0; buf.Len() < capacity.DataBits; i++ {
			switch i % 4 {
			case 0:
				buf.WriteBitsLSB(0b1110, 4)
			case 1:
				buf.WriteBitsLSB(0b1100, 4)
			case 2:
				buf.WriteBitsLSB(0b0001, 4)
			case 3:
				buf.WriteBitsLSB(0b0001, 4)
			}
		}
		buf.WriteBitsLSB(0, capacity.Data*8-buf.Len())
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

// length returns the length of s in bits.
func (s *Segment) length(version Version) (int, bool) {
	var n int

	// mode indicator
	switch version {
	case 2:
		n = 1
	case 3:
		n = 2
	case 4:
		n = 3
	}

	switch s.Mode {
	case ModeNumeric:
		switch version {
		case 1:
			n += 3
		case 2:
			n += 4
		case 3:
			n += 5
		case 4:
			n += 6
		default:
			return 0, false
		}
		n += 10 * (len(s.Data) / 3)
		switch len(s.Data) % 3 {
		case 1:
			n += 4
		case 2:
			n += 7
		}
		return n, true
	case ModeAlphanumeric:
		switch version {
		case 2:
			n += 3
		case 3:
			n += 4
		case 4:
			n += 5
		default:
			return 0, false
		}
		n += 11 * (len(s.Data) / 2)
		if len(s.Data)%2 != 0 {
			n += 6
		}
		return n, true
	case ModeBytes:
		switch version {
		case 3:
			n += 4
		case 4:
			n += 5
		default:
			return 0, false
		}
		n += len(s.Data) * 8
		return n, true
	default:
		return 0, false
	}
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
	return bitstream.EncodeBytes(buf, data)
}
