package microqr

import (
	"bytes"
	"image/png"
	"testing"
)

func TestNew1(t *testing.T) {
	qr, err := New([]byte("MICROQR"), WithLevel(LevelL))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 3 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 3)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeAlphanumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeAlphanumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("MICROQR")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "MICROQR")
	}
}

func TestNew2(t *testing.T) {
	// Maximum size for Version M4
	qr, err := New([]byte("12345678901234567890123456789012345"), WithLevel(LevelL))
	if err != nil {
		t.Fatal(err)
	}

	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 4 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 4)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeNumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeNumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("12345678901234567890123456789012345")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "12345678901234567890123456789012345")
	}
}

func TestNew3(t *testing.T) {
	qr, err := New([]byte("123456789012345678901234"), WithLevel(LevelL))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 4 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 4)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeNumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeNumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("123456789012345678901234")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "123456789012345678901234")
	}
}

func TestNew4(t *testing.T) {
	// maximum size for Version M3
	qr, err := New([]byte("12345678901234567890123"), WithLevel(LevelL))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 3 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 3)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeNumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeNumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("12345678901234567890123")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "12345678901234567890123")
	}
}

func TestNew5(t *testing.T) {
	_, err := New([]byte("123456789012345678901234567890123456"), WithLevel(LevelL))
	if err == nil {
		t.Fatal("want error, but not")
	}
}

func TestNew6(t *testing.T) {
	// maximum size for Version M1
	qr, err := New([]byte("12345"), WithLevel(LevelCheck))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 1 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 1)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeNumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeNumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("12345")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "12345")
	}
}

func TestNewFromKanji1(t *testing.T) {
	qr, err := New([]byte("MICROQR"), WithLevel(LevelL), WithKanji(true))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 3 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 3)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeAlphanumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeAlphanumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("MICROQR")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "MICROQR")
	}
}

func TestNewFromKanji2(t *testing.T) {
	qr, err := New([]byte("点"), WithLevel(LevelL), WithKanji(true))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 3 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 3)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeKanji {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeKanji)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("点")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "点")
	}
}

func TestNewFromKanji3(t *testing.T) {
	qr, err := New([]byte("免"), WithLevel(LevelL), WithKanji(true))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 3 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 3)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeKanji {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeKanji)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("免")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "免")
	}
}

func TestEncode1(t *testing.T) {
	qr := &QRCode{
		Version: 2,
		Level:   LevelL,
		Mask:    Mask1,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("01234567"),
			},
		},
	}
	img, err := qr.Encode()
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
}

func TestEncodeToBitmap1(t *testing.T) {
	qr := &QRCode{
		Version: 2,
		Level:   LevelL,
		Mask:    Mask1,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("01234567"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	// X 0510 : 2018
	// 附属書I (参考)
	// シンボルの符号化例
	want := []byte{
		0b11111110, 0b10101000,
		0b10000010, 0b11101000,
		0b10111010, 0b01101000,
		0b10111010, 0b01111000,
		0b10111010, 0b11100000,
		0b10000010, 0b10001000,
		0b11111110, 0b01111000,
		0b00000000, 0b01100000,
		0b11010000, 0b10001000,
		0b01101010, 0b10101000,
		0b11100111, 0b11110000,
		0b00010100, 0b00110000,
		0b11101001, 0b10111000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap2(t *testing.T) {
	qr := &QRCode{
		Version: 3,
		Level:   LevelM,
		Mask:    Mask3,
		Segments: []Segment{
			{
				Mode: ModeAlphanumeric,
				Data: []byte("MICROQR"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010,
		0b10000010, 0b10011110,
		0b10111010, 0b10111100,
		0b10111010, 0b11111010,
		0b10111010, 0b00000110,
		0b10000010, 0b01110100,
		0b11111110, 0b01001010,
		0b00000000, 0b01001000,
		0b10001001, 0b10101010,
		0b00111100, 0b00011010,
		0b10100000, 0b10001000,
		0b00101111, 0b01000010,
		0b10110001, 0b01001100,
		0b01110010, 0b11010100,
		0b10110111, 0b11011100,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap3(t *testing.T) {
	qr := &QRCode{
		Version: 2,
		Level:   LevelM,
		Mask:    Mask0,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("12345"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101000,
		0b10000010, 0b10000000,
		0b10111010, 0b11101000,
		0b10111010, 0b00011000,
		0b10111010, 0b01110000,
		0b10000010, 0b10001000,
		0b11111110, 0b00101000,
		0b00000000, 0b01011000,
		0b11100111, 0b10000000,
		0b01111001, 0b01100000,
		0b11100000, 0b01110000,
		0b01001000, 0b10101000,
		0b11111110, 0b10011000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap4(t *testing.T) {
	qr := &QRCode{
		Version: 3,
		Level:   LevelM,
		Mask:    Mask2,
		Segments: []Segment{
			{
				Mode: ModeBytes,
				Data: []byte("1haicso"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	// from https://www.qrcode.com/img/rmqr/gra2.jpg
	want := []byte{
		0b11111110, 0b10101010,
		0b10000010, 0b00111110,
		0b10111010, 0b00011110,
		0b10111010, 0b01100110,
		0b10111010, 0b00010010,
		0b10000010, 0b11101010,
		0b11111110, 0b11001100,
		0b00000000, 0b01000010,
		0b10001100, 0b11101000,
		0b00011100, 0b00110010,
		0b10111000, 0b01001100,
		0b00010001, 0b00100000,
		0b11001001, 0b10011000,
		0b00100000, 0b01110010,
		0b11011110, 0b01011110,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap5(t *testing.T) {
	qr := &QRCode{
		Version: 3,
		Level:   LevelL,
		Mask:    Mask1,
		Segments: []Segment{
			{
				Mode: ModeAlphanumeric,
				Data: []byte("AINIX"),
			},
			{
				Mode: ModeNumeric,
				Data: []byte("12345"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010,
		0b10000010, 0b11000010,
		0b10111010, 0b11110110,
		0b10111010, 0b11100100,
		0b10111010, 0b11101110,
		0b10000010, 0b00101110,
		0b11111110, 0b00011100,
		0b00000000, 0b11111010,
		0b11110011, 0b00001100,
		0b00110010, 0b00010100,
		0b11101010, 0b01101110,
		0b00110001, 0b11100000,
		0b10011111, 0b10000010,
		0b01101101, 0b00000010,
		0b10011111, 0b00101100,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap6(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelCheck,
		Mask:    Mask2,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("0000"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10100000,
		0b10000010, 0b10000000,
		0b10111010, 0b11100000,
		0b10111010, 0b00100000,
		0b10111010, 0b11000000,
		0b10000010, 0b00100000,
		0b11111110, 0b11100000,
		0b00000000, 0b00000000,
		0b11001110, 0b01100000,
		0b01011101, 0b00100000,
		0b10100111, 0b11100000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap7(t *testing.T) {
	qr := &QRCode{
		Version: 3,
		Level:   LevelM,
		Mask:    Mask2,
		Segments: []Segment{
			{
				Mode: ModeKanji,
				Data: []byte("点"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010,
		0b10000010, 0b01101110,
		0b10111010, 0b00001100,
		0b10111010, 0b00010010,
		0b10111010, 0b00010010,
		0b10000010, 0b10000000,
		0b11111110, 0b10011000,
		0b00000000, 0b01001000,
		0b10001100, 0b10011000,
		0b01000010, 0b10001000,
		0b11011010, 0b10101110,
		0b01000110, 0b01100110,
		0b10111110, 0b11110100,
		0b00010011, 0b01010110,
		0b11010101, 0b11101010,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}
