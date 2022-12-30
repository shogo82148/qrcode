package qrcode

import (
	"bytes"
	"testing"

	"github.com/shogo82148/qrcode/internal/bitstream"
)

func TestNew1(t *testing.T) {
	qr, err := New([]byte("01234567"), WithLevel(LevelH))
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
	if !bytes.Equal(qr.Segments[0].Data, []byte("01234567")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "01234567")
	}
}

func TestNew2(t *testing.T) {
	qr, err := New([]byte("Ver1"), WithLevel(LevelH))
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
	if qr.Segments[0].Mode != ModeBytes {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeBytes)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("Ver1")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "Ver1")
	}
}

func TestNew3(t *testing.T) {
	qr, err := New([]byte("VERSION 10 QR CODE, UP TO 174 CHAR AT H LEVEL, WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND. NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"), WithLevel(LevelH))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 10 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 1)
	}
	if len(qr.Segments) != 5 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 5)
	}
	if qr.Segments[0].Mode != ModeAlphanumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeAlphanumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("VERSION 10 QR CODE")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "VERSION 10 QR CODE")
	}
}

func TestNewFromKanji1(t *testing.T) {
	qr, err := New([]byte("点茗"), WithLevel(LevelH), WithKanji(true))
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
	if qr.Segments[0].Mode != ModeKanji {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeKanji)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("点茗")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "点茗")
	}
}

func TestNewFromKanji2(t *testing.T) {
	qr, err := New([]byte("点"), WithLevel(LevelH), WithKanji(true))
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
	if qr.Segments[0].Mode != ModeKanji {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeKanji)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("点")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "点")
	}
}

func TestNewFromKanji3(t *testing.T) {
	qr, err := New([]byte("VERSION 10 QR CODE, UP TO 174 CHAR AT H LEVEL, WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND. NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"), WithLevel(LevelH), WithKanji(true))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Mask != MaskAuto {
		t.Errorf("unexpected mask: got %v, want %v", qr.Mask, MaskAuto)
	}
	if qr.Version != 10 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, 10)
	}
	if len(qr.Segments) != 5 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 5)
	}
	if qr.Segments[0].Mode != ModeAlphanumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeBytes)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("VERSION 10 QR CODE")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "VERSION 10 QR CODE")
	}
}

func TestNewFromKanji4(t *testing.T) {
	qr, err := New([]byte("1234567"), WithLevel(LevelH), WithKanji(true))
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
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 2)
	}
	if qr.Segments[0].Mode != ModeNumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeNumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("1234567")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "1234567")
	}
}

func TestSegment_Length(t *testing.T) {
	test := func(version Version, s Segment) {
		t.Helper()
		l := s.length(version)
		var buf bitstream.Buffer
		if err := s.encode(version, &buf); err != nil {
			t.Fatal(err)
		}
		if l != buf.Len() {
			t.Errorf("length mismatch: want %d, got %d", buf.Len(), l)
		}
	}

	test(1, Segment{
		Mode: ModeNumeric,
		Data: []byte{'0'},
	})
	test(1, Segment{
		Mode: ModeNumeric,
		Data: []byte{'0', '1'},
	})
	test(1, Segment{
		Mode: ModeNumeric,
		Data: []byte{'0', '1', '2'},
	})
	test(1, Segment{
		Mode: ModeNumeric,
		Data: []byte{'0', '1', '2', '3'},
	})
	test(1, Segment{
		Mode: ModeAlphanumeric,
		Data: []byte{'A'},
	})
	test(1, Segment{
		Mode: ModeAlphanumeric,
		Data: []byte{'A', 'B'},
	})
	test(1, Segment{
		Mode: ModeBytes,
		Data: []byte{'a'},
	})
}

func TestQRCode_EncodeToBitmap1(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelM,
		Mask:    Mask2,
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
		0b11111110, 0b01011011, 0b11111000,
		0b10000010, 0b01111010, 0b00001000,
		0b10111010, 0b10000010, 0b11101000,
		0b10111010, 0b11000010, 0b11101000,
		0b10111010, 0b10111010, 0b11101000,
		0b10000010, 0b10001010, 0b00001000,
		0b11111110, 0b10101011, 0b11111000,
		0b00000000, 0b10011000, 0b00000000,
		0b10111110, 0b01001011, 0b11100000,
		0b00010101, 0b10101001, 0b01100000,
		0b00100011, 0b01010100, 0b11111000,
		0b00001000, 0b01000001, 0b11100000,
		0b00011111, 0b10010100, 0b10000000,
		0b00000000, 0b10111110, 0b01100000,
		0b11111110, 0b01101011, 0b00000000,
		0b10000010, 0b10111110, 0b00101000,
		0b10111010, 0b10001001, 0b01100000,
		0b10111010, 0b11001001, 0b00000000,
		0b10111010, 0b10110100, 0b10100000,
		0b10000010, 0b00000001, 0b10110000,
		0b11111110, 0b11110100, 0b10100000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestQRCode_EncodeToBitmap2(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelH,
		Mask:    0b001,
		Segments: []Segment{
			{
				Mode: ModeBytes,
				Data: []byte("Ver1"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	// from Wikipedia: https://en.wikipedia.org/wiki/QR_code
	// https://commons.wikimedia.org/wiki/File:Qr-1.png
	want := []byte{
		0b11111110, 0b00110011, 0b11111000,
		0b10000010, 0b10001010, 0b00001000,
		0b10111010, 0b10001010, 0b11101000,
		0b10111010, 0b10001010, 0b11101000,
		0b10111010, 0b10111010, 0b11101000,
		0b10000010, 0b10100010, 0b00001000,
		0b11111110, 0b10101011, 0b11111000,
		0b00000000, 0b00000000, 0b00000000,
		0b00100111, 0b11110101, 0b11110000,
		0b00011101, 0b11111010, 0b01001000,
		0b11011011, 0b00101010, 0b00101000,
		0b01110000, 0b10010100, 0b11001000,
		0b11100011, 0b00111100, 0b00001000,
		0b00000000, 0b11001101, 0b00010000,
		0b11111110, 0b11111101, 0b11001000,
		0b10000010, 0b11001011, 0b00000000,
		0b10111010, 0b00100101, 0b10001000,
		0b10111010, 0b00011110, 0b00000000,
		0b10111010, 0b11000110, 0b00111000,
		0b10000010, 0b01101010, 0b10000000,
		0b11111110, 0b00110010, 0b01101000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestQRCode_EncodeToBitmap3(t *testing.T) {
	qr := &QRCode{
		Version: 2,
		Level:   LevelH,
		Mask:    0b010,
		Segments: []Segment{
			{
				Mode: ModeBytes,
				Data: []byte("Version 2"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}

	got := img.Pix

	// from Wikipedia: https://en.wikipedia.org/wiki/QR_code
	want := []byte{
		0b11111110, 0b10011001, 0b10111111, 0b10000000,
		0b10000010, 0b10011101, 0b00100000, 0b10000000,
		0b10111010, 0b11111110, 0b10101110, 0b10000000,
		0b10111010, 0b00101010, 0b00101110, 0b10000000,
		0b10111010, 0b00111010, 0b00101110, 0b10000000,
		0b10000010, 0b10101000, 0b10100000, 0b10000000,
		0b11111110, 0b10101010, 0b10111111, 0b10000000,
		0b00000000, 0b10101111, 0b10000000, 0b00000000,
		0b00111010, 0b10100111, 0b01110011, 0b10000000,
		0b00000100, 0b10000001, 0b10100101, 0b00000000,
		0b00100011, 0b01110010, 0b01001111, 0b10000000,
		0b00011000, 0b11100101, 0b00000001, 0b10000000,
		0b00010010, 0b10111011, 0b01011011, 0b10000000,
		0b10000100, 0b00000101, 0b10100100, 0b00000000,
		0b10101010, 0b11111000, 0b00100011, 0b10000000,
		0b10110100, 0b00011010, 0b11100000, 0b00000000,
		0b10001110, 0b11001000, 0b11111111, 0b10000000,
		0b00000000, 0b11011000, 0b10001111, 0b10000000,
		0b11111110, 0b00011111, 0b10101001, 0b10000000,
		0b10000010, 0b00000001, 0b10001001, 0b10000000,
		0b10111010, 0b10000111, 0b11111110, 0b00000000,
		0b10111010, 0b11000010, 0b00101010, 0b10000000,
		0b10111010, 0b11101001, 0b01100010, 0b10000000,
		0b10000010, 0b00000111, 0b01001000, 0b10000000,
		0b11111110, 0b01111010, 0b10011111, 0b10000000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestQRCode_EncodeToBitmap4(t *testing.T) {
	qr := &QRCode{
		Version: 3,
		Level:   LevelH,
		Mask:    0b001,
		Segments: []Segment{
			{
				Mode: ModeBytes,
				Data: []byte("Version 3 QR Code"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	// from Wikipedia: https://en.wikipedia.org/wiki/QR_code
	want := []byte{
		0b11111110, 0b00000010, 0b01011011, 0b11111000,
		0b10000010, 0b10110100, 0b10001010, 0b00001000,
		0b10111010, 0b10010011, 0b11011010, 0b11101000,
		0b10111010, 0b10100010, 0b10101010, 0b11101000,
		0b10111010, 0b11110111, 0b10010010, 0b11101000,
		0b10000010, 0b11011010, 0b01100010, 0b00001000,
		0b11111110, 0b10101010, 0b10101011, 0b11111000,
		0b00000000, 0b01010001, 0b11000000, 0b00000000,
		0b00100111, 0b11100111, 0b11011101, 0b11110000,
		0b00010001, 0b10101010, 0b11001001, 0b10001000,
		0b00100110, 0b01101110, 0b11001001, 0b00001000,
		0b01101001, 0b11111001, 0b10011111, 0b11001000,
		0b11111011, 0b01011110, 0b11001011, 0b01001000,
		0b11010001, 0b11101110, 0b11110011, 0b01000000,
		0b11111111, 0b10001100, 0b10110110, 0b10101000,
		0b10110001, 0b11000111, 0b00010100, 0b11000000,
		0b11010111, 0b11010100, 0b00110111, 0b10011000,
		0b01100000, 0b00011001, 0b00000100, 0b01110000,
		0b11111110, 0b00100001, 0b01010000, 0b10001000,
		0b00000001, 0b01000110, 0b10100011, 0b01010000,
		0b11011011, 0b10011010, 0b00011111, 0b11011000,
		0b00000000, 0b10001100, 0b00011000, 0b10001000,
		0b11111110, 0b10110111, 0b00101010, 0b11011000,
		0b10000010, 0b11010010, 0b11001000, 0b11001000,
		0b10111010, 0b01000011, 0b11011111, 0b10011000,
		0b10111010, 0b01011100, 0b00100000, 0b11010000,
		0b10111010, 0b11101010, 0b11111100, 0b11111000,
		0b10000010, 0b01100101, 0b11000111, 0b11000000,
		0b11111110, 0b00001001, 0b10000101, 0b11001000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestQRCode_EncodeToBitmap5(t *testing.T) {
	qr := &QRCode{
		Version: 10,
		Level:   LevelH,
		Mask:    0b100,
		Segments: []Segment{
			{
				Mode: ModeAlphanumeric,
				Data: []byte("VERSION 10 QR CODE"),
			},
			{
				Mode: ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: ModeAlphanumeric,
				Data: []byte(" UP TO 174 CHAR AT H LEVEL"),
			},
			{
				Mode: ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: ModeAlphanumeric,
				Data: []byte(
					" WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND." +
						"  NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	// from Wikipedia: https://en.wikipedia.org/wiki/QR_code
	// https://commons.wikimedia.org/wiki/File:Qr-code-ver-10.png
	want := []byte{
		0b11111110, 0b01101010, 0b11010000, 0b10001101, 0b10100110, 0b10101111, 0b00111111, 0b10000000,
		0b10000010, 0b11011011, 0b10010011, 0b10110101, 0b10110010, 0b01010001, 0b00100000, 0b10000000,
		0b10111010, 0b01111101, 0b00001111, 0b01110111, 0b01000111, 0b11001111, 0b00101110, 0b10000000,
		0b10111010, 0b01000101, 0b11101000, 0b10010110, 0b11110011, 0b01111001, 0b00101110, 0b10000000,
		0b10111010, 0b01110001, 0b00101100, 0b00111111, 0b10101100, 0b11011001, 0b00101110, 0b10000000,
		0b10000010, 0b10101100, 0b11000001, 0b01100011, 0b01101010, 0b00101110, 0b00100000, 0b10000000,
		0b11111110, 0b10101010, 0b10101010, 0b10101010, 0b10101010, 0b10101010, 0b10111111, 0b10000000,
		0b00000000, 0b10110000, 0b10111100, 0b11100011, 0b01000000, 0b10110111, 0b10000000, 0b00000000,
		0b00001111, 0b01110110, 0b00000011, 0b00111110, 0b01011111, 0b10011111, 0b10110001, 0b00000000,
		0b11100001, 0b01000111, 0b11101110, 0b10100100, 0b10001110, 0b00001000, 0b00110110, 0b10000000,
		0b00111010, 0b10001101, 0b10101001, 0b01110110, 0b00001000, 0b11011110, 0b00100100, 0b00000000,
		0b10000001, 0b10000001, 0b00111100, 0b11110001, 0b10110110, 0b10001001, 0b00100000, 0b10000000,
		0b01010110, 0b01010001, 0b01111000, 0b00111111, 0b01111011, 0b10000010, 0b11001001, 0b10000000,
		0b01011001, 0b00010110, 0b10111011, 0b10111010, 0b00100111, 0b00011001, 0b10110010, 0b00000000,
		0b11010010, 0b00111010, 0b00000010, 0b00110011, 0b00001001, 0b01001110, 0b01000001, 0b10000000,
		0b00000001, 0b00100001, 0b10101000, 0b10100110, 0b00000110, 0b10110110, 0b11011101, 0b00000000,
		0b11000010, 0b11001001, 0b00111100, 0b00100101, 0b11000100, 0b01101111, 0b01111100, 0b00000000,
		0b01110100, 0b01000111, 0b10011110, 0b01010110, 0b01001010, 0b01101101, 0b01010101, 0b10000000,
		0b00101010, 0b11011010, 0b10100100, 0b00111011, 0b11001110, 0b01011110, 0b10010101, 0b00000000,
		0b11010000, 0b00110011, 0b01011100, 0b00010100, 0b10000000, 0b10100101, 0b11110010, 0b00000000,
		0b00011010, 0b00101100, 0b01011110, 0b01011000, 0b11101011, 0b00000110, 0b10111001, 0b10000000,
		0b10100000, 0b10101000, 0b00010001, 0b01011001, 0b10110101, 0b00001100, 0b01010100, 0b10000000,
		0b10011110, 0b00101101, 0b00001111, 0b00100011, 0b01101011, 0b01101111, 0b10011011, 0b00000000,
		0b01011000, 0b11100010, 0b01110100, 0b00001001, 0b01111010, 0b01011011, 0b11011100, 0b10000000,
		0b01111011, 0b00011111, 0b11001010, 0b00011100, 0b10100001, 0b00100011, 0b10111101, 0b10000000,
		0b00001001, 0b00101000, 0b10000111, 0b01101101, 0b00111111, 0b10110100, 0b10000111, 0b10000000,
		0b01101111, 0b10011001, 0b10000000, 0b00111111, 0b00001101, 0b00111110, 0b11111111, 0b00000000,
		0b00011000, 0b11001001, 0b11111100, 0b11100011, 0b01010000, 0b11010101, 0b10001100, 0b00000000,
		0b01111010, 0b10010000, 0b11101010, 0b10101010, 0b01001100, 0b11011000, 0b10101010, 0b00000000,
		0b10011000, 0b11011011, 0b11110100, 0b10100010, 0b00101101, 0b11111111, 0b10001000, 0b00000000,
		0b01101111, 0b10110101, 0b10110001, 0b00111110, 0b00000011, 0b11000110, 0b11111011, 0b00000000,
		0b00010001, 0b00110000, 0b10100111, 0b11110100, 0b10110001, 0b00101011, 0b10010100, 0b10000000,
		0b10001010, 0b00101111, 0b00010001, 0b11111001, 0b01111011, 0b11010110, 0b11000011, 0b10000000,
		0b11110000, 0b11110001, 0b11111000, 0b11001000, 0b11111001, 0b01010010, 0b10111100, 0b00000000,
		0b11101010, 0b01010000, 0b01101000, 0b00111011, 0b10110101, 0b01000100, 0b01100101, 0b00000000,
		0b00010001, 0b11011101, 0b01011100, 0b10110010, 0b00010100, 0b11011110, 0b00110011, 0b00000000,
		0b11100010, 0b11011100, 0b00001111, 0b01110000, 0b11011101, 0b10001011, 0b10111010, 0b10000000,
		0b11001000, 0b00110011, 0b10101101, 0b01011101, 0b01010000, 0b10100000, 0b01010010, 0b10000000,
		0b11110011, 0b00100100, 0b11011011, 0b00101111, 0b00001010, 0b11000101, 0b01001000, 0b10000000,
		0b00111100, 0b11101100, 0b01000010, 0b10011010, 0b01101110, 0b11111010, 0b00111010, 0b10000000,
		0b10110111, 0b11101001, 0b10101101, 0b01100010, 0b01100100, 0b10110010, 0b00110110, 0b00000000,
		0b01100000, 0b01011000, 0b11110011, 0b00001101, 0b11100001, 0b00110010, 0b01101100, 0b00000000,
		0b01011010, 0b01110011, 0b01001011, 0b10000100, 0b01010000, 0b00000110, 0b10000000, 0b10000000,
		0b01100000, 0b00000100, 0b10000011, 0b01111001, 0b10010010, 0b00101101, 0b10110101, 0b10000000,
		0b11101010, 0b01110000, 0b11010100, 0b01111001, 0b01100011, 0b11101111, 0b00100101, 0b10000000,
		0b00101001, 0b11010101, 0b00000000, 0b10111001, 0b10111011, 0b00110100, 0b00000100, 0b10000000,
		0b10100111, 0b11110000, 0b01100101, 0b10010010, 0b01001100, 0b01101001, 0b10011010, 0b10000000,
		0b11111000, 0b01110110, 0b00111100, 0b10100110, 0b10000001, 0b10111001, 0b00110101, 0b00000000,
		0b00000011, 0b10101010, 0b11111111, 0b01111110, 0b11001010, 0b00010000, 0b11111101, 0b10000000,
		0b00000000, 0b10001011, 0b01011101, 0b11100011, 0b10000110, 0b01110101, 0b10001011, 0b10000000,
		0b11111110, 0b10000100, 0b11000011, 0b00101010, 0b00001100, 0b00110000, 0b10101110, 0b00000000,
		0b10000010, 0b10010011, 0b11011001, 0b10100011, 0b10010001, 0b11111111, 0b10001010, 0b00000000,
		0b10111010, 0b10011011, 0b11001110, 0b00111110, 0b10110010, 0b11101001, 0b11111011, 0b10000000,
		0b10111010, 0b01110100, 0b01011001, 0b01010000, 0b01100100, 0b10001001, 0b00101001, 0b10000000,
		0b10111010, 0b01001011, 0b11000010, 0b00010101, 0b01111110, 0b11111110, 0b10001110, 0b00000000,
		0b10000010, 0b01111000, 0b01101110, 0b11101010, 0b00100110, 0b10101001, 0b10110010, 0b10000000,
		0b11111110, 0b00001000, 0b11010101, 0b00101100, 0b11111111, 0b11100000, 0b01010101, 0b10000000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestQRCode_EncodeToBitmap6(t *testing.T) {
	qr := &QRCode{
		Version: 1,
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

	// generated by QRQR
	want := []byte{
		0b11111110, 0b00010011, 0b11111000,
		0b10000010, 0b00100010, 0b00001000,
		0b10111010, 0b10100010, 0b11101000,
		0b10111010, 0b11001010, 0b11101000,
		0b10111010, 0b11011010, 0b11101000,
		0b10000010, 0b10101010, 0b00001000,
		0b11111110, 0b10101011, 0b11111000,
		0b00000000, 0b10011000, 0b00000000,
		0b10111110, 0b01001011, 0b11100000,
		0b11100101, 0b11001001, 0b00011000,
		0b11000010, 0b01010100, 0b10111000,
		0b01110001, 0b00100001, 0b10100000,
		0b11001010, 0b10010100, 0b10111000,
		0b00000000, 0b11011110, 0b01001000,
		0b11111110, 0b00101011, 0b00010000,
		0b10000010, 0b10011110, 0b01010000,
		0b10111010, 0b11001001, 0b00100000,
		0b10111010, 0b11101001, 0b00100000,
		0b10111010, 0b11110100, 0b11100000,
		0b10000010, 0b00000001, 0b10100000,
		0b11111110, 0b11010100, 0b11101000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestQRCode_EncodeToBitmap7(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelH,
		Mask:    Mask0,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("000000000"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10000011, 0b11111000,
		0b10000010, 0b00010010, 0b00001000,
		0b10111010, 0b01111010, 0b11101000,
		0b10111010, 0b10111010, 0b11101000,
		0b10111010, 0b01101010, 0b11101000,
		0b10000010, 0b01100010, 0b00001000,
		0b11111110, 0b10101011, 0b11111000,
		0b00000000, 0b01011000, 0b00000000,
		0b00101110, 0b11001100, 0b01001000,
		0b01010100, 0b10111001, 0b01010000,
		0b00110110, 0b01100001, 0b00101000,
		0b01101000, 0b00100111, 0b11010000,
		0b11001111, 0b10111111, 0b00101000,
		0b00000000, 0b10011000, 0b01010000,
		0b11111110, 0b01110010, 0b10111000,
		0b10000010, 0b11111000, 0b01011000,
		0b10111010, 0b11110100, 0b10101000,
		0b10111010, 0b01101011, 0b01010000,
		0b10111010, 0b10110001, 0b00101000,
		0b10000010, 0b01100101, 0b11000000,
		0b11111110, 0b01001101, 0b00101000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestQRCode_EncodeToBitmap_ErrorTooLong(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelH,
		Mask:    0b010,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("123456789012345678"),
			},
		},
	}
	_, err := qr.EncodeToBitmap()
	if err == nil {
		t.Error("want error, got not")
	}
}

func TestQRCode_EncodeToBitmap_ErrorInvalidNumber(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelM,
		Mask:    0b010,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("A"),
			},
		},
	}
	_, err := qr.EncodeToBitmap()
	if err == nil {
		t.Error("want error, got not")
	}
}

func TestQRCode_EncodeToBitmap_ErrorInvalidAlphanumeric(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelM,
		Mask:    0b010,
		Segments: []Segment{
			{
				Mode: ModeAlphanumeric,
				Data: []byte("a"),
			},
		},
	}
	_, err := qr.EncodeToBitmap()
	if err == nil {
		t.Error("want error, got not")
	}
}

func TestQRCode_encodeToBits1(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelM,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("01234567"),
			},
		},
	}

	var buf bitstream.Buffer
	qr.encodeToBits(&buf)
	got := buf.Bytes()

	// JIS X 0510 : 2018
	// 附属書I (参考)
	// シンボルの符号化例
	want := []byte{
		// data
		0b0001_0000, 0b0010_0000, 0b0000_1100, 0b0101_0110,
		0b0110_0001, 0b1000_0000,

		// padding
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,

		// error correction code
		0b1010_0101, 0b0010_0100, 0b1101_0100, 0b1100_0001,
		0b1110_1101, 0b0011_0110, 0b1100_0111, 0b1000_0111,
		0b0010_1100, 0b0101_0101,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("unexpected result:\ngot  %08b,\nwant %08b", got, want)
	}
}

func TestQRCode_encodeToBits2(t *testing.T) {
	qr := &QRCode{
		Version: 10,
		Level:   LevelH,
		Segments: []Segment{
			{
				Mode: ModeAlphanumeric,
				Data: []byte("VERSION 10 QR CODE"),
			},
			{
				Mode: ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: ModeAlphanumeric,
				Data: []byte(" UP TO 174 CHAR AT H LEVEL"),
			},
			{
				Mode: ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: ModeAlphanumeric,
				Data: []byte(
					" WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND." +
						"  NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"),
			},
		},
	}

	var buf bitstream.Buffer
	qr.encodeToBits(&buf)
	got := buf.Bytes()

	want := []byte{
		0b00100000, 0b00000000, 0b00011001, 0b00111011, 0b11101001, 0b01010111, 0b00100111, 0b11000100,
		0b00100101, 0b00000100, 0b11011111, 0b00100001, 0b01000110, 0b00010011, 0b11001110, 0b11110011,
		0b01100000, 0b10110000, 0b11001100, 0b00011101, 0b11010011, 0b11110011, 0b00101100, 0b10101111,
		0b01100110, 0b10000000, 0b10111001, 0b00010111, 0b10111100, 0b10001100, 0b00011110, 0b10011100,
		0b11011011, 0b11010110, 0b10100101, 0b01001000, 0b11011000, 0b01011100, 0b10011110, 0b01100110,
		0b01000010, 0b01110010, 0b01001010, 0b10101111, 0b10110001, 0b01011101, 0b11010010, 0b01001010,
		0b10000101, 0b10010001, 0b10101000, 0b10111111, 0b11010100, 0b00011001, 0b10100101, 0b00110000,
		0b11100000, 0b00110100, 0b10110100, 0b01100100, 0b10001100, 0b01111010, 0b00100010, 0b01101000,
		0b10110111, 0b11000111, 0b00000000, 0b11001110, 0b11010111, 0b01101011, 0b10100110, 0b00101011,
		0b00110111, 0b00101010, 0b00000001, 0b11101111, 0b10011100, 0b11010101, 0b10001110, 0b11010001,
		0b01001110, 0b10010011, 0b00101100, 0b11101000, 0b01101000, 0b11010100, 0b11101010, 0b00000011,
		0b00110100, 0b11111100, 0b00100000, 0b10000001, 0b11010010, 0b11100111, 0b10011010, 0b11011101,
		0b01101000, 0b11000000, 0b11011111, 0b11011001, 0b01101101, 0b10011110, 0b00111001, 0b10110111,
		0b10010101, 0b11000001, 0b10011101, 0b01001101, 0b00101000, 0b00100001, 0b11101001, 0b00000000,
		0b11010000, 0b11100111, 0b00011010, 0b10110001, 0b00101010, 0b00011101, 0b01101110, 0b11101100,
		0b10010101, 0b00010001, 0b00011011, 0b11000101, 0b10101100, 0b11100001, 0b11011000, 0b00101011,
		0b11110100, 0b00010100, 0b01000000, 0b11011010, 0b11110001, 0b11001001, 0b11100100, 0b11000101,
		0b00000001, 0b11101010, 0b10001011, 0b10010011, 0b01111011, 0b11010110, 0b11111000, 0b00101011,
		0b10011011, 0b01111101, 0b11010111, 0b00111110, 0b11111010, 0b11010011, 0b11101100, 0b01110100,
		0b01110111, 0b00000110, 0b00001011, 0b00001100, 0b01100100, 0b11011010, 0b00111101, 0b01100101,
		0b01000110, 0b11110000, 0b01111011, 0b00011101, 0b10001001, 0b00001111, 0b10101110, 0b01010011,
		0b10110111, 0b00101011, 0b01111100, 0b10101111, 0b10111011, 0b00110001, 0b10101110, 0b10001111,
		0b10011001, 0b11101111, 0b00111000, 0b00011101, 0b10001011, 0b10101101, 0b11111110, 0b00010011,
		0b10111110, 0b10001010, 0b01101101, 0b11001011, 0b01111101, 0b01101001, 0b10100010, 0b10100011,
		0b00011101, 0b10010101, 0b00010100, 0b00011001, 0b10001001, 0b00000110, 0b00110111, 0b11011110,
		0b10000111, 0b11100010, 0b00101111, 0b11001111, 0b11010001, 0b10000000, 0b11011110, 0b01001110,
		0b01011010, 0b11010110, 0b11000011, 0b11111011, 0b00010000, 0b10001101, 0b00101011, 0b11100010,
		0b01110111, 0b10101100, 0b00101011, 0b11101001, 0b11000000, 0b11101101, 0b01010110, 0b01010100,
		0b10101101, 0b01111100, 0b01000110, 0b10000001, 0b11000000, 0b11001110, 0b00101001, 0b00110101,
		0b11100010, 0b00011001, 0b11001101, 0b01101000, 0b11000000, 0b10000000, 0b11100111, 0b11110110,
		0b10101011, 0b10011000, 0b11111001, 0b11110100, 0b01110010, 0b11001110, 0b01101110, 0b01100101,
		0b11111001, 0b10001000, 0b00000001, 0b00011110, 0b10110111, 0b00001111, 0b00101110, 0b11000011,
		0b00010110, 0b11110110, 0b01101011, 0b11011100, 0b10111101, 0b01001110, 0b00001011, 0b11100000,
		0b10011101, 0b00111001, 0b10001100, 0b01100100, 0b10111001, 0b01010111, 0b10000000, 0b01001100,
		0b10001100, 0b10101011, 0b01100011, 0b00001101, 0b01000011, 0b00011001, 0b01011111, 0b10011110,
		0b00110101, 0b01100110, 0b10001100, 0b10111101, 0b00010000, 0b00001111, 0b01000100, 0b11011101,
		0b10001100, 0b00000111, 0b11100001, 0b01110000, 0b00101000, 0b01101111, 0b01111011, 0b00001101,
		0b10010011, 0b01001011, 0b10011101, 0b11110100, 0b01010101, 0b10011010, 0b01000010, 0b11100001,
		0b01011011, 0b10100110, 0b01001101, 0b00111101, 0b00001000, 0b11110011, 0b10000000, 0b00110011,
		0b00001100, 0b10011101, 0b11000111, 0b01101010, 0b00011100, 0b10110101, 0b00101111, 0b10100010,
		0b11100111, 0b11001010, 0b11010000, 0b10010111, 0b01001011, 0b00011100, 0b10001011, 0b01010000,
		0b11000010, 0b01010100, 0b00111001, 0b00111110, 0b00110110, 0b00000011, 0b00111001, 0b10001111,
		0b10010000, 0b00010101, 0b11000001, 0b00111110, 0b01100111, 0b10000110, 0b10001000, 0b11000000,
		0b11001001, 0b10100011,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("unexpected result:\ngot  %08b,\nwant %08b", got, want)
	}
}

func TestQRCode_encodeSegments(t *testing.T) {
	qr := &QRCode{
		Version: 10,
		Level:   LevelH,
		Mask:    0b100,
		Segments: []Segment{
			{
				Mode: ModeAlphanumeric,
				Data: []byte("VERSION 10 QR CODE"),
			},
			{
				Mode: ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: ModeAlphanumeric,
				Data: []byte(" UP TO 174 CHAR AT H LEVEL"),
			},
			{
				Mode: ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: ModeAlphanumeric,
				Data: []byte(
					" WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND." +
						"  NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"),
			},
		},
	}

	var buf bitstream.Buffer
	qr.encodeSegments(&buf)
	got := buf.Bytes()
	want := []byte{
		// data
		0b00100000, 0b00100101, 0b01100000, 0b01100110, 0b11011011, 0b01000010, 0b10000101, 0b11100000,
		0b10110111, 0b00110111, 0b01001110, 0b00110100, 0b01101000, 0b10010101, 0b11010000, 0b00000000,
		0b00000100, 0b10110000, 0b10000000, 0b11010110, 0b01110010, 0b10010001, 0b00110100, 0b11000111,
		0b00101010, 0b10010011, 0b11111100, 0b11000000, 0b11000001, 0b11100111, 0b00011001, 0b11011111,
		0b11001100, 0b10111001, 0b10100101, 0b01001010, 0b10101000, 0b10110100, 0b00000000, 0b00000001,
		0b00101100, 0b00100000, 0b11011111, 0b10011101, 0b00011010, 0b00111011, 0b00100001, 0b00011101,
		0b00010111, 0b01001000, 0b10101111, 0b10111111, 0b01100100, 0b11001110, 0b11101111, 0b11101000,
		0b10000001, 0b11011001, 0b01001101, 0b10110001, 0b11101001, 0b01000110, 0b11010011, 0b10111100,
		0b11011000, 0b10110001, 0b11010100, 0b10001100, 0b11010111, 0b10011100, 0b01101000, 0b11010010,
		0b01101101, 0b00101000, 0b00101010, 0b01010111, 0b00010011, 0b11110011, 0b10001100, 0b01011100,
		0b01011101, 0b00011001, 0b01111010, 0b01101011, 0b11010101, 0b11010100, 0b11100111, 0b10011110,
		0b00100001, 0b00011101, 0b00100111, 0b11001110, 0b00101100, 0b00011110, 0b10011110, 0b11010010,
		0b10100101, 0b00100010, 0b10100110, 0b10001110, 0b11101010, 0b10011010, 0b00111001, 0b11101001,
		0b01101110, 0b10010101, 0b11000100, 0b11110011, 0b10101111, 0b10011100, 0b01100110, 0b01001010,
		0b00110000, 0b01101000, 0b00101011, 0b11010001, 0b00000011, 0b11011101, 0b10110111, 0b00000000,
		0b11101100, 0b00010001,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("unexpected result:\ngot  %08b,\nwant %08b", got, want)
	}
}

func TestSegment_encodeNumber(t *testing.T) {
	s := &Segment{
		Mode: ModeNumeric,
		Data: []byte("01234567"),
	}
	var buf bitstream.Buffer

	if err := s.encode(1, &buf); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0b0001_0000, 0b001000_00, 0b00001100, 0b01010110, 0b01_100001, 0b1_0000000,
	}
	got := buf.Bytes()

	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestSegment_encodeAlphabet(t *testing.T) {
	s := &Segment{
		Mode: ModeAlphanumeric,
		Data: []byte("AC-42"),
	}
	var buf bitstream.Buffer

	if err := s.encode(1, &buf); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0b0010_0000, 0b00101_001, 0b11001110, 0b11100111, 0b001_00001, 0b00_000000,
	}
	got := buf.Bytes()

	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestSegment_encodeBytes(t *testing.T) {
	s := &Segment{
		Mode: ModeBytes,
		Data: []byte{0xAA, 0x55},
	}
	var buf bitstream.Buffer

	if err := s.encode(1, &buf); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0b0100_0000, 0b0010_1010, 0b1010_0101, 0b0101_0000,
	}
	got := buf.Bytes()

	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestSegment_encodeKanji(t *testing.T) {
	s := &Segment{
		Mode: ModeKanji,
		Data: []byte("点"),
	}
	var buf bitstream.Buffer

	if err := s.encode(1, &buf); err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0b1000_0000, 0b0001_0110, 0b1100_1111, 0b1000_0000,
	}
	got := buf.Bytes()

	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func BenchmarkEncode(b *testing.B) {
	qr := &QRCode{
		Version: 40,
		Level:   LevelM,
		Mask:    0b010,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: bytes.Repeat([]byte("9"), 3057),
			},
		},
	}

	for i := 0; i < b.N; i++ {
		_, err := qr.EncodeToBitmap()
		if err != nil {
			b.Fatal(err)
		}
	}
}
