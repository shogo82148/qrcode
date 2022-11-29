package microqr

//go:generate go run genbch/main.go

import "strconv"

type QRCode struct {
	Version  Version
	Level    Level
	Mask     Mask
	Segments []Segment
}

type Version int

type Level int

const (
	LevelCheck Level = 0b10
	LevelL     Level = 0b01
	LevelM     Level = 0b00
	LevelQ     Level = 0b11
)

func (lv Level) String() string {
	switch lv {
	case LevelCheck:
		return "Check"
	case LevelL:
		return "L"
	case LevelM:
		return "M"
	case LevelQ:
		return "Q"
	}
	return "invalid(" + strconv.Itoa(int(lv)) + ")"
}

type Mask int

const (
	Mask0 Mask = iota
	Mask1
	Mask2
	Mask3
	maskMax

	MaskAuto Mask = -1
)

type Mode uint8

const (
	// ModeNumeric is number mode.
	// The Data must be ascii characters [0-9].
	ModeNumeric Mode = 0b000

	// ModeAlphanumeric is alphabet and number mode.
	// The Data must be ascii characters [0-9A-Z $%*+\-./:].
	ModeAlphanumeric Mode = 0b001

	// ModeBytes is 8-bit bytes mode.
	// The Data can include any bytes.
	ModeBytes Mode = 0b010

	// ModeKanji is Japanese Kanji mode.
	ModeKanji Mode = 0b011

	ModeTerminated Mode = 0b0000
)

func (mode Mode) String() string {
	switch mode {
	case ModeNumeric:
		return "numeric"
	case ModeAlphanumeric:
		return "alphanumeric"
	case ModeBytes:
		return "bytes"
	case ModeKanji:
		return "kanji"
	default:
		return "(unknown mode: " + strconv.Itoa(int(mode)) + ")"
	}
}

type Segment struct {
	Mode Mode
	Data []byte
}
