package microqr

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
	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

func New(data []byte, opts ...EncodeOptions) (*QRCode, error) {
	myopts := newEncodeOptions(opts...)
	lv := myopts.Level
	if lv < 0 || lv >= 4 {
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
		version := calcVersion(level, nil)
		return &QRCode{
			Version: version,
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
					cost += (4 + 6) * 6
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
					cost += (4 + 5) * 6
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
					cost += (4 + 5) * 6
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
		return nil, errors.New("microqr: data too large")
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
		version := calcVersion(level, nil)
		return &QRCode{
			Version: version,
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
LOOP:
	for version := Version(1); version <= 4; version++ {
		if formatTable[version][level] < 0 {
			continue
		}
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
		Level:      LevelQ,
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

func WithKanji(use bool) EncodeOptions {
	return func(opts *encodeOptions) {
		opts.Kanji = use
	}
}

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
	myopts := newEncodeOptions(opts...)

	binimg, err := qr.EncodeToBitmap()
	if err != nil {
		return nil, err
	}

	w := binimg.Bounds().Dx() + myopts.QuietZone*2

	// convert bitmap to image
	src := fp16.NewNRGBAh(image.Rect(0, 0, w, w))
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
	dst := fp16.NewNRGBAh(image.Rect(0, 0, W, W))
	resize.AreaAverage(dst, src)

	return srgb.EncodeTone(dst), nil
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
	}
	buf.WriteBitsLSB(0, capacity.Data*8-buf.Len())

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
	case ModeKanji:
		switch version {
		case 3:
			n += 3
		case 4:
			n += 4
		default:
			return 0, false
		}
		n += utf8.RuneCount(s.Data) * 13
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
	case ModeKanji:
		return s.encodeKanji(version, buf)
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

func (s *Segment) encodeKanji(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch version {
	case 3:
		n = 3
	case 4:
		n = 4
	default:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	}
	count := utf8.RuneCount(data)
	if count >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", count)
	}

	// mode
	switch version {
	case 3:
		buf.WriteBitsLSB(uint64(ModeKanji), 2)
	case 4:
		buf.WriteBitsLSB(uint64(ModeKanji), 3)
	}

	// data length
	buf.WriteBitsLSB(uint64(count), n)

	// data
	return bitstream.EncodeKanji(buf, data)
}
