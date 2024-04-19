package rmqr

import (
	"errors"
	"fmt"
	"image"
	"math"
	"unicode/utf8"

	"github.com/shogo82148/go-imaging/bitmap"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/resize"
	"github.com/shogo82148/go-imaging/srgb"
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
		return newFromKanji(lv, myopts.Priority, data)
	} else {
		return newQR(lv, myopts.Priority, data)
	}
}

func newQR(level Level, priority Priority, data []byte) (*QRCode, error) {
	if len(data) == 0 {
		return &QRCode{
			Version: R7x43,
			Level:   level,
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
					cost += (3 + 9) * 6
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
					cost += (3 + 8) * 6
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
					cost += (3 + 8) * 6
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

	version, ok := calcVersion(level, priority, segments)
	if !ok {
		return nil, errors.New("qrcode: data too large")
	}

	return &QRCode{
		Version:  version,
		Level:    level,
		Segments: segments,
	}, nil
}

func newFromKanji(level Level, priority Priority, data []byte) (*QRCode, error) {
	if len(data) == 0 {
		return &QRCode{
			Version: R7x43,
			Level:   level,
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

	version, ok := calcVersion(level, priority, segments)
	if !ok {
		return nil, errors.New("qrcode: data too large")
	}

	return &QRCode{
		Version:  version,
		Level:    level,
		Segments: segments,
	}, nil
}

func calcVersion(level Level, priority Priority, segments []Segment) (Version, bool) {
	if !level.IsValid() {
		return 0, false
	}

	var order []Version
	switch priority {
	case PriorityArea:
		order = capacityOrderArea
	case PriorityHeight:
		order = capacityOrderHeight
	case PriorityWidth:
		order = capacityOrderWidth
	default:
		return 0, false
	}

LOOP:
	for _, version := range order {
		capacity := capacityTable[version][level].Data * 8
		length := 0
		for _, s := range segments {
			l, ok := s.length(version, level)
			if !ok {
				continue LOOP
			}
			length += l
			if length > capacity {
				continue LOOP
			}
		}
		if length <= capacity {
			return version, true
		}
	}
	return 0, false
}

type EncodeOptions func(opts *encodeOptions)

type encodeOptions struct {
	QuietZone  int
	ModuleSize float64
	Level      Level
	Kanji      bool
	Priority   Priority
	Width      int
}

func newEncodeOptions(opts ...EncodeOptions) encodeOptions {
	myopts := encodeOptions{
		QuietZone:  2,
		ModuleSize: 1,
		Level:      LevelM,
		Kanji:      true,
		Priority:   PriorityArea,
		Width:      0,
	}
	for _, o := range opts {
		o(&myopts)
	}
	return myopts
}

// WithModuleSize sets the module size.
// The default value is 1.
func WithModuleSize(size float64) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.ModuleSize = size
	}
}

// WithQuietZone sets the quiet zone size.
// The default value is 2.
func WithQuietZone(n int) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.QuietZone = n
	}
}

// WithLevel sets the error correction level.
// The default value is LevelM.
func WithLevel(lv Level) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.Level = lv
	}
}

// WithKanji enables the kanji mode.
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
// The image height is calculated from the image width.
func WithWidth(width int) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.Width = width
	}
}

// Priority is a priority for selecting the version.
type Priority int

const (
	// PriorityArea selects the version that minimizes the area.
	PriorityArea Priority = iota

	// PriorityHeight selects the version that minimizes the height.
	PriorityHeight

	// PriorityWidth selects the version that minimizes the width.
	PriorityWidth
)

// WithPriority sets the priority for selecting the version.
func WithPriority(priority Priority) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.Priority = priority
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
	h := binimg.Bounds().Dy() + myopts.QuietZone*2

	// convert bitmap to image
	src := fp16.NewNRGBAh(image.Rect(0, 0, w, h))
	black := fp16color.NewNRGBAh(0, 0, 0, 1)
	white := fp16color.NewNRGBAh(1, 1, 1, 1)
	for y := 0; y < w; y++ {
		for x := 0; x < w; x++ {
			c := binimg.BinaryAt(x-myopts.QuietZone, y-myopts.QuietZone)
			if c {
				src.SetNRGBAh(x, y, black)
			} else {
				src.SetNRGBAh(x, y, white)
			}
		}
	}

	// resize
	W := max(
		int(math.Ceil(float64(w)*myopts.ModuleSize)),
		myopts.Width,
	)
	H := int(math.Ceil(float64(h) * float64(W) / float64(w)))
	dst := fp16.NewNRGBAh(image.Rect(0, 0, W, H))
	resize.AreaAverage(dst, src)

	return srgb.EncodeTone(dst), nil
}

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
	if l > capacity.Data*8 {
		return errors.New("qrcode: data is too large")
	}

	// terminate pattern
	if capacity.Data*8-l > 3 {
		buf.WriteBitsLSB(uint64(ModeTerminated), 3)
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

// length returns the length of s in bits.
func (s *Segment) length(version Version, level Level) (int, bool) {
	if int(version) >= len(capacityTable) {
		return 0, false
	}
	if int(level) >= len(capacityTable[version]) {
		return 0, false
	}

	capacity := capacityTable[version][level]

	switch s.Mode {
	case ModeNumeric:
		n := capacity.BitLength[ModeNumeric]
		if len(s.Data) >= 1<<n {
			return 0, false
		}
		m := 10 * (len(s.Data) / 3)
		switch len(s.Data) % 3 {
		case 1:
			n += 4
		case 2:
			n += 7
		}
		return 3 + n + m, true
	case ModeAlphanumeric:
		n := capacity.BitLength[ModeAlphanumeric]
		if len(s.Data) >= 1<<n {
			return 0, false
		}
		m := 11 * (len(s.Data) / 2)
		if len(s.Data)%2 != 0 {
			m += 6
		}
		return 3 + n + m, true
	case ModeBytes:
		n := capacity.BitLength[ModeBytes]
		if len(s.Data) >= 1<<n {
			return 0, false
		}
		m := len(s.Data) * 8
		return 3 + n + m, true
	case ModeKanji:
		n := capacity.BitLength[ModeKanji]
		if len(s.Data) >= 1<<n {
			return 0, false
		}
		m := len(s.Data) * 13
		return 3 + n + m, true
	default:
		return 0, false
	}
}

func (s *Segment) encode(bitLength [5]int, buf *bitstream.Buffer) error {
	switch s.Mode {
	case ModeNumeric:
		return s.encodeNumber(bitLength[ModeNumeric], buf)
	case ModeAlphanumeric:
		return s.encodeAlphanumeric(bitLength[ModeAlphanumeric], buf)
	case ModeBytes:
		return s.encodeBytes(bitLength[ModeBytes], buf)
	case ModeKanji:
		return s.encodeKanji(bitLength[ModeKanji], buf)
	default:
		return errors.New("qrcode: unknown mode")
	}
}

func (s *Segment) encodeNumber(n int, buf *bitstream.Buffer) error {
	if len(s.Data) >= 1<<n {
		return fmt.Errorf("rmqr: data is too long for number mode: %d", len(s.Data))
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeNumeric), 3)

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	return bitstream.EncodeNumeric(buf, s.Data)
}

func (s *Segment) encodeAlphanumeric(n int, buf *bitstream.Buffer) error {
	if len(s.Data) >= 1<<n {
		return fmt.Errorf("rmqr: data is too long for alphanumeric mode: %d", len(s.Data))
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeAlphanumeric), 3)

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	return bitstream.EncodeAlphanumeric(buf, s.Data)
}

func (s *Segment) encodeBytes(n int, buf *bitstream.Buffer) error {
	if len(s.Data) >= 1<<n {
		return fmt.Errorf("rmqr: data is too long for bytes: %d", len(s.Data))
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeBytes), 3)

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	return bitstream.EncodeBytes(buf, s.Data)
}

func (s *Segment) encodeKanji(n int, buf *bitstream.Buffer) error {
	if len(s.Data) >= 1<<n {
		return fmt.Errorf("rmqr: data is too long for bytes: %d", len(s.Data))
	}

	// mode
	buf.WriteBitsLSB(uint64(ModeKanji), 3)

	// data length
	count := utf8.RuneCount(s.Data)
	buf.WriteBitsLSB(uint64(count), n)

	// data
	return bitstream.EncodeKanji(buf, s.Data)
}
