package qrcode

import (
	"bytes"
	"testing"

	"github.com/shogo82148/qrcode/internal/binimage"
	"github.com/shogo82148/qrcode/internal/bitstream"
)

func TestQRCode_Encode(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelM,
		Mask:    0b010,
		Segments: []Segment{
			{
				Mode: ModeNumber,
				Data: []byte("01234567"),
			},
		},
	}
	img, err := qr.Encode()
	if err != nil {
		t.Fatal(err)
	}

	got := img.(*binimage.Binary).Pix

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

func TestQRCode_encodeToBits(t *testing.T) {
	qr := &QRCode{
		Version: 1,
		Level:   LevelM,
		Segments: []Segment{
			{
				Mode: ModeNumber,
				Data: []byte("01234567"),
			},
		},
	}

	var buf bitstream.Buffer
	qr.encodeToBits(&buf)
	got := buf.Bytes()
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

func TestSegment_encodeNumber(t *testing.T) {
	s := &Segment{
		Mode: ModeNumber,
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
