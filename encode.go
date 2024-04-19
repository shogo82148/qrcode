package qrcode

//go:generate go run genbch/main.go

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"unicode/utf8"

	bitmap "github.com/shogo82148/qrcode/bitmap"
	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func New(data []byte, opts ...EncodeOptions) (*QRCode, error) {
	myopts := newEncodeOptions(opts...)
	lv := myopts.Level
	if !lv.IsValid() {
		return nil, fmt.Errorf("qrcode: invalid level: %d", lv)
	}
	if myopts.Kanji {
		return newFromKanji(lv, data)
	} else {
		return newQR(lv, data)
	}
}

func newQR(level Level, data []byte) (*QRCode, error) {
	if len(data) == 0 {
		return &QRCode{
			Version: 1,
			Level:   level,
			Mask:    MaskAuto,
		}, nil
	}

	const inf = math.MaxInt - 1<<18 // 1<<18 is for avoiding overflow
	const (
		modeInit = iota
		modeNumeric
		modeAlphanumeric
		modeBytes
		modeMax
	)
	type state struct {
		cost     int // = bit length * 6
		lastMode int
	}
	states := make([][4]state, len(data)+1)
	states[0][modeNumeric].cost = inf
	states[0][modeAlphanumeric].cost = inf
	states[0][modeBytes].cost = inf

	for i := 0; i < len(data); i++ {
		if i != 0 {
			states[i][modeInit].cost = inf
		}
		// numeric
		if bitstream.IsNumeric(data[i]) {
			minCost := inf
			lastMode := modeInit
			for mode := modeInit; mode < modeMax; mode++ {
				cost := states[i][mode].cost + 20
				if mode != modeNumeric {
					cost += (4 + 14) * 6
				}
				if cost < minCost {
					minCost = cost
					lastMode = mode
				}
			}

			states[i+1][modeNumeric] = state{
				cost:     minCost,
				lastMode: lastMode,
			}
		} else {
			states[i+1][modeNumeric] = state{
				cost:     inf,
				lastMode: modeInit,
			}
		}

		// alphanumeric
		if bitstream.IsAlphanumeric(data[i]) {
			minCost := inf
			lastMode := modeInit
			for mode := modeInit; mode < modeMax; mode++ {
				cost := states[i][mode].cost + 33
				if mode != modeAlphanumeric {
					cost += (4 + 13) * 6
				}
				if cost < minCost {
					minCost = cost
					lastMode = mode
				}
			}

			states[i+1][modeAlphanumeric] = state{
				cost:     minCost,
				lastMode: lastMode,
			}
		} else {
			states[i+1][modeAlphanumeric] = state{
				cost:     inf,
				lastMode: modeInit,
			}
		}

		// bytes
		{
			minCost := inf
			lastMode := modeInit
			for mode := modeInit; mode < modeMax; mode++ {
				cost := states[i][mode].cost + 8*6
				if mode != modeBytes {
					cost += (4 + 16) * 6
				}
				if cost < minCost {
					minCost = cost
					lastMode = mode
				}
			}
			states[i+1][modeBytes] = state{
				cost:     minCost,
				lastMode: lastMode,
			}
		}
	}

	best := make([]int, len(data))
	minCost := states[len(data)][modeNumeric].cost
	bestMode := modeNumeric
	if cost := states[len(data)][modeAlphanumeric].cost; cost < minCost {
		minCost = cost
		bestMode = modeAlphanumeric
	}
	if cost := states[len(data)][modeBytes].cost; cost < minCost {
		bestMode = modeBytes
	}
	best[len(data)-1] = bestMode
	for i := len(data) - 1; i > 0; i-- {
		bestMode = states[i+1][bestMode].lastMode
		best[i-1] = bestMode
	}

	modeList := [...]Mode{0, ModeNumeric, ModeAlphanumeric, ModeBytes}
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
		return nil, errors.New("qrcode: data too large")
	}

	return &QRCode{
		Version:  version,
		Level:    level,
		Mask:     MaskAuto,
		Segments: segments,
	}, nil
}

func newFromKanji(level Level, data []byte) (*QRCode, error) {
	if len(data) == 0 {
		return &QRCode{
			Version: 1,
			Level:   level,
			Mask:    MaskAuto,
		}, nil
	}

	const inf = math.MaxInt - 1<<18 // 1<<18 is for avoiding overflow
	const (
		modeInit = iota
		modeNumeric
		modeAlphanumeric
		modeBytes
		modeKanji
		modeMax
	)
	type state struct {
		cost     int // = bit length * 6
		lastMode int
		data     []byte
	}
	states := make([][5]state, len(data)+1)
	states[0][modeNumeric].cost = inf
	states[0][modeAlphanumeric].cost = inf
	states[0][modeBytes].cost = inf
	states[0][modeKanji].cost = inf

	for i := 0; i < len(data); i++ {
		if i != 0 {
			states[i][modeInit].cost = inf
		}
		// numeric
		if bitstream.IsNumeric(data[i]) {
			minCost := states[i][modeInit].cost + (4+14)*6 + 20
			lastMode := modeInit
			for mode := modeInit + 1; mode < modeMax; mode++ {
				cost := states[i][mode].cost + 20
				if mode != modeNumeric {
					cost += (4 + 14) * 6
				}
				if cost < minCost {
					minCost = cost
					lastMode = mode
				}
			}
			states[i+1][modeNumeric] = state{
				cost:     minCost,
				lastMode: lastMode,
				data:     data[i : i+1],
			}
		} else {
			states[i+1][modeNumeric] = state{
				cost:     inf,
				lastMode: modeInit,
				data:     []byte{},
			}
		}

		// alphanumeric
		if bitstream.IsAlphanumeric(data[i]) {
			minCost := states[i][modeInit].cost + (4+13)*6 + 33
			lastMode := modeInit
			for mode := modeInit + 1; mode < modeMax; mode++ {
				cost := states[i][mode].cost + 33
				if mode != modeAlphanumeric {
					cost += (4 + 13) * 6
				}
				if cost < minCost {
					minCost = cost
					lastMode = mode
				}
			}
			states[i+1][modeAlphanumeric] = state{
				cost:     minCost,
				lastMode: lastMode,
				data:     data[i : i+1],
			}
		} else {
			states[i+1][modeAlphanumeric] = state{
				cost:     inf,
				lastMode: modeInit,
				data:     []byte{},
			}
		}

		// bytes
		{
			minCost := states[i][modeInit].cost + (4+16+8)*6
			lastMode := modeInit
			for mode := modeInit + 1; mode < modeMax; mode++ {
				cost := states[i][mode].cost + 8*6
				if mode != modeBytes {
					cost += (4 + 16) * 6
				}
				if cost < minCost {
					minCost = cost
					lastMode = mode
				}
			}
			states[i+1][modeBytes] = state{
				cost:     minCost,
				lastMode: lastMode,
				data:     data[i : i+1],
			}
		}

		// kanji
		r, size := utf8.DecodeRune(data[i:])
		if r != utf8.RuneError && bitstream.IsKanji(r) && i+size < len(states) {
			minCost := states[i][modeInit].cost + (4+12+13)*6
			lastMode := modeInit
			for mode := modeInit + 1; mode < modeMax; mode++ {
				cost := states[i][mode].cost + 13*6
				if mode != modeKanji {
					cost += (4 + 12) * 6
				}
				if cost < minCost {
					minCost = cost
					lastMode = mode
				}
			}

			for j := 0; j < size; j++ {
				states[i+j][modeKanji].cost = inf
			}

			states[i+size][modeKanji] = state{
				cost:     minCost,
				lastMode: lastMode,
				data:     data[i : i+size],
			}
		} else if states[i+1][modeKanji].data == nil {
			states[i+1][modeKanji] = state{
				cost:     inf,
				lastMode: modeInit,
				data:     []byte{},
			}
		}
	}

	// find the best path
	type elem struct {
		mode int
		data []byte
	}
	best := make([]elem, 0, len(data))
	minCost := states[len(data)][modeNumeric].cost
	bestMode := modeNumeric
	for mode := modeNumeric + 1; mode < modeMax; mode++ {
		if cost := states[len(data)][mode].cost; cost < minCost {
			minCost = cost
			bestMode = mode
		}
	}
	best = append(best, elem{
		mode: bestMode,
		data: states[len(data)][bestMode].data,
	})
	for i := len(data); ; {
		size := len(states[i][bestMode].data)
		bestMode = states[i][bestMode].lastMode
		if bestMode == modeInit {
			break
		}
		i -= size
		best = append(best, elem{
			mode: bestMode,
			data: states[i][bestMode].data,
		})
	}

	modeList := [...]Mode{0, ModeNumeric, ModeAlphanumeric, ModeBytes, ModeKanji}
	segments := []Segment{
		{
			Mode: modeList[best[len(best)-1].mode],
			Data: best[len(best)-1].data,
		},
	}
	for i := len(best) - 2; i >= 0; i-- {
		mode := modeList[best[i].mode]
		if segments[len(segments)-1].Mode == mode {
			segments[len(segments)-1].Data = append(segments[len(segments)-1].Data, best[i].data...)
		} else {
			segments = append(segments, Segment{
				Mode: mode,
				Data: best[i].data,
			})
		}
	}

	version := calcVersion(level, segments)
	if version == 0 {
		return nil, errors.New("qrcode: data too large")
	}

	return &QRCode{
		Version:  version,
		Level:    level,
		Mask:     MaskAuto,
		Segments: segments,
	}, nil
}

func calcVersion(level Level, segments []Segment) Version {
	if !level.IsValid() {
		return 0
	}

LOOP:
	for version := Version(1); version <= 40; version++ {
		capacity := capacityTable[version][level].Data * 8
		length := 0
		for _, s := range segments {
			l := s.length(version)
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

const timingPatternOffset = 6

func skipTimingPattern(n int) int {
	if n < timingPatternOffset {
		return n
	}
	return n + 1
}

type EncodeOptions func(opts *encodeOptions)

type encodeOptions struct {
	QuietZone  int
	ModuleSize float64
	Level      Level
	Kanji      bool
	Width      int
}

func newEncodeOptions(opts ...EncodeOptions) encodeOptions {
	myopts := encodeOptions{
		QuietZone:  4,
		ModuleSize: 1,
		Level:      LevelM,
		Kanji:      true,
		Width:      0,
	}
	for _, o := range opts {
		o(&myopts)
	}
	return myopts
}

// WithModuleSize sets the module size.
// The default size is 1.
func WithModuleSize(size float64) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.ModuleSize = size
	}
}

// WithQuietZone sets the quiet zone size.
// The default size is 4.
func WithQuietZone(n int) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.QuietZone = n
	}
}

// WithLevel sets the error correction level.
// The default level is LevelM.
// If the level is invalid, it panics.
func WithLevel(lv Level) EncodeOptions {
	if !lv.IsValid() {
		panic(fmt.Sprintf("qrcode: invalid level: %d", lv))
	}
	return func(opts *encodeOptions) {
		opts.Level = lv
	}
}

// WithKanji sets the kanji mode.
// The default mode is true.
// If it's enabled, Shift-JIS encoding is used for kanji mode.
func WithKanji(use bool) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.Kanji = use
	}
}

// WithWidth sets the width of the image.
// The larger of the image width calculated from [WithModuleSize]
// and the image width specified with WithWidth is used.
func WithWidth(width int) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.Width = width
	}
}

func Encode(data []byte, opts ...EncodeOptions) (image.Image, error) {
	qr, err := New(data, opts...)
	if err != nil {
		return nil, err
	}
	return qr.Encode(opts...)
}

func (qr *QRCode) Encode(opts ...EncodeOptions) (image.Image, error) {
	if !qr.Version.IsValid() {
		return nil, errors.New("qrcode: invalid version")
	}
	if !qr.Level.IsValid() {
		return nil, errors.New("qrcode: invalid level")
	}

	myopts := newEncodeOptions(opts...)

	binimg, err := qr.EncodeToBitmap()
	if err != nil {
		return nil, err
	}

	w := binimg.Bounds().Dx() + myopts.QuietZone*2
	W := int(math.Ceil(float64(w) * myopts.ModuleSize))
	scale := myopts.ModuleSize
	dX := float64(myopts.QuietZone) * scale
	dY := float64(myopts.QuietZone) * scale

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
	if !qr.Version.IsValid() {
		return nil, errors.New("qrcode: invalid version")
	}
	if !qr.Level.IsValid() {
		return nil, errors.New("qrcode: invalid level")
	}

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
	for _, s := range qr.Segments {
		if err := s.encode(qr.Version, buf); err != nil {
			return err
		}
	}
	l := buf.Len()
	capacity := capacityTable[qr.Version][qr.Level]
	if l > capacity.Data*8 {
		return errors.New("qrcode: data is too large")
	}

	// terminate pattern
	if capacity.Data*8-l > 4 {
		buf.WriteBitsLSB(uint64(ModeTerminated), 4)
	}

	// align to bytes
	if mod := buf.Len() % 8; mod != 0 {
		buf.WriteBitsLSB(0x00, int(8-mod))
	}

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
	case ModeKanji:
		return s.encodeKanji(version, buf)
	default:
		return errors.New("qrcode: unknown mode")
	}
}

// length returns the length of s in bits.
func (s *Segment) length(version Version) int {
	var n int = 4 // mode indicator
	switch s.Mode {
	case ModeNumeric:
		switch {
		case version <= 0 || version > 40:
			panic(fmt.Errorf("qrcode: invalid version: %d", version))
		case version < 10:
			n += 10
		case version < 27:
			n += 12
		default:
			n += 14
		}
		n += 10 * (len(s.Data) / 3)
		switch len(s.Data) % 3 {
		case 1:
			n += 4
		case 2:
			n += 7
		}
		return n
	case ModeAlphanumeric:
		switch {
		case version <= 0 || version > 40:
			panic(fmt.Errorf("qrcode: invalid version: %d", version))
		case version < 10:
			n += 9
		case version < 27:
			n += 11
		default:
			n += 13
		}
		n += 11 * (len(s.Data) / 2)
		if len(s.Data)%2 != 0 {
			n += 6
		}
		return n
	case ModeBytes:
		switch {
		case version <= 0 || version > 40:
			panic(fmt.Errorf("qrcode: invalid version: %d", version))
		case version < 10:
			n += 8
		default:
			n += 16
		}
		n += len(s.Data) * 8
		return n
	case ModeKanji:
		switch {
		case version <= 0 || version > 40:
			panic(fmt.Errorf("qrcode: invalid version: %d", version))
		case version < 10:
			n += 8
		case version < 27:
			n += 10
		default:
			n += 12
		}
		n += utf8.RuneCount(s.Data) * 13
		return n
	default:
		panic(errors.New("qrcode: unknown mode"))
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

	// mode
	buf.WriteBitsLSB(uint64(ModeNumeric), 4)

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	return bitstream.EncodeNumeric(buf, data)
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

	// mode
	buf.WriteBitsLSB(uint64(ModeAlphanumeric), 4)

	// data length
	buf.WriteBitsLSB(uint64(len(data)), n)

	return bitstream.EncodeAlphanumeric(buf, data)
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
	return bitstream.EncodeBytes(buf, data)
}

func (s *Segment) encodeKanji(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch {
	case version <= 0 || version > 40:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	case version < 10:
		n = 8
	case version < 27:
		n = 10
	default:
		n = 12
	}
	count := utf8.RuneCount(data)
	if count >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeKanji), 4)

	// data length
	buf.WriteBitsLSB(uint64(count), n)

	// data
	return bitstream.EncodeKanji(buf, data)
}
