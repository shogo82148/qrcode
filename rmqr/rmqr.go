//go:generate go run genbase/main.go
//go:generate go run genbch/main.go

package rmqr

import (
	"math"
	"strconv"
)

type QRCode struct {
	Version  Version
	Level    Level
	Segments []Segment
}

type Version int

type Level int

const (
	LevelM Level = 0b0
	LevelH Level = 0b1
)

func (lv Level) String() string {
	switch lv {
	case LevelM:
		return "M"
	case LevelH:
		return "H"
	}
	return "invalid(" + strconv.Itoa(int(lv)) + ")"
}

type Mode uint8

const (
	// ModeNumeric is number mode.
	// The Data must be ascii characters [0-9].
	ModeNumeric Mode = 0b001

	// ModeAlphanumeric is alphabet and number mode.
	// The Data must be ascii characters [0-9A-Z $%*+\-./:].
	ModeAlphanumeric Mode = 0b010

	// ModeBytes is 8-bit bytes mode.
	// The Data can include any bytes.
	ModeBytes Mode = 0b011

	// ModeKanji is Japanese Kanji mode.
	ModeKanji Mode = 0b100

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

func round(x float64) int {
	return int(math.Round(x))
}
