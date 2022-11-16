package qrcode

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/shogo82148/qrcode/internal/bitstream"
)

func TestQRCode_Encode(t *testing.T) {
	qr := &QRCode{
		Version: 1,
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

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("qrcode.png", buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestQRCode_encodeToBits(t *testing.T) {
	qr := &QRCode{
		Version: 1,
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
		0b0001_0000, 0b0010_0000, 0b0000_1100, 0b0101_0110,
		0b0110_0001, 0b1000_0000,

		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,
		0b1110_1100, 0b0001_0001,

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
		Mode: ModeAlphabet,
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
