package qrcode

//go:generate go run genbase/main.go
//go:generate go run genbch/main.go

import (
	"math"
	"strconv"
)

type QRCode struct {
	Version  Version
	Level    Level
	Mask     Mask
	Segments []Segment
}

// Version is a version of QR code.
type Version int

const versionMin Version = 0 // 0 is special case for auto version
const versionMax Version = 41

func (v Version) IsValid() bool {
	return versionMin <= v && v < versionMax
}

// Level is a error correction level.
type Level int

const (
	levelMin Level = 0
	LevelL   Level = 0b01
	LevelM   Level = 0b00
	LevelQ   Level = 0b11
	LevelH   Level = 0b10
	levelMax Level = 4
)

// IsValid returns true if the level is valid.
func (lv Level) IsValid() bool {
	return levelMin <= lv && lv < levelMax
}

func (lv Level) String() string {
	switch lv {
	case LevelL:
		return "L"
	case LevelM:
		return "M"
	case LevelQ:
		return "Q"
	case LevelH:
		return "H"
	}
	return "invalid(" + strconv.Itoa(int(lv)) + ")"
}

// Mask is a mask pattern.
type Mask int

const (
	maskMin Mask = 0b000
	Mask0   Mask = 0b000
	Mask1   Mask = 0b001
	Mask2   Mask = 0b010
	Mask3   Mask = 0b011
	Mask4   Mask = 0b100
	Mask5   Mask = 0b101
	Mask6   Mask = 0b110
	Mask7   Mask = 0b111
	maskMax Mask = 0b111 + 1

	MaskAuto Mask = -1
)

// IsValid returns true if the mask is valid.
func (m Mask) IsValid() bool {
	return m == MaskAuto || maskMin <= m && m < maskMax
}

type Mode uint8

const (
	// ModeECI is ECI(Extended Channel Interpretation) mode.
	ModeECI Mode = 0b0111

	// ModeNumeric is number mode.
	// The Data must be ascii characters [0-9].
	ModeNumeric Mode = 0b0001

	// ModeAlphanumeric is alphabet and number mode.
	// The Data must be ascii characters [0-9A-Z $%*+\-./:].
	ModeAlphanumeric Mode = 0b0010

	// ModeBytes is 8-bit bytes mode.
	// The Data can include any bytes.
	ModeBytes Mode = 0b0100

	// ModeKanji is Japanese Kanji mode.
	ModeKanji Mode = 0b1000

	// ModeConnected is connected structure mode.
	ModeConnected Mode = 0b0011

	ModeFNC1_1 Mode = 0b0101
	ModeFNC1_2 Mode = 0b1001

	ModeTerminated Mode = 0b0000
)

func (mode Mode) String() string {
	switch mode {
	case ModeECI:
		return "eci"
	case ModeNumeric:
		return "numeric"
	case ModeAlphanumeric:
		return "alphanumeric"
	case ModeBytes:
		return "bytes"
	case ModeKanji:
		return "kanji"
	case ModeConnected:
		return "connected"
	case ModeFNC1_1:
		return "fcn1-1"
	case ModeFNC1_2:
		return "fnc1-2"
	default:
		return "(unknown mode: " + strconv.Itoa(int(mode)) + ")"
	}
}

type Segment struct {
	Mode Mode
	Data []byte
}

func round(x float64) int {
	return int(math.Round(x))
}
