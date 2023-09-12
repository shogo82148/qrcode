//go:generate go run genbase/main.go
//go:generate go run genbch/main.go

package rmqr

import (
	"fmt"
	"math"
	"strconv"
)

type QRCode struct {
	Version  Version
	Level    Level
	Segments []Segment
}

type Version int

const (
	R7x43   Version = 0b00000
	R7x59   Version = 0b00001
	R7x77   Version = 0b00010
	R7x99   Version = 0b00011
	R7x139  Version = 0b00100
	R9x43   Version = 0b00101
	R9x59   Version = 0b00110
	R9x77   Version = 0b00111
	R9x99   Version = 0b01000
	R9x139  Version = 0b01001
	R11x27  Version = 0b01010
	R11x43  Version = 0b01011
	R11x59  Version = 0b01100
	R11x77  Version = 0b01101
	R11x99  Version = 0b01110
	R11x139 Version = 0b01111
	R13x27  Version = 0b10000
	R13x43  Version = 0b10001
	R13x59  Version = 0b10010
	R13x77  Version = 0b10011
	R13x99  Version = 0b10100
	R13x139 Version = 0b10101
	R15x43  Version = 0b10110
	R15x59  Version = 0b10111
	R15x77  Version = 0b11000
	R15x99  Version = 0b11001
	R15x139 Version = 0b11010
	R17x43  Version = 0b11011
	R17x59  Version = 0b11100
	R17x77  Version = 0b11101
	R17x99  Version = 0b11110
	R17x139 Version = 0b11111
)

func (version Version) String() string {
	base := baseList[version]
	return fmt.Sprintf("R%dx%d", base.Rect.Dy(), base.Rect.Dx())
}

// Width returns the width of version.
func (version Version) Width() int {
	base := baseList[version]
	return base.Rect.Dx()
}

// Height returns the width of version.
func (version Version) Height() int {
	base := baseList[version]
	return base.Rect.Dy()
}

type Level int

const (
	LevelM   Level = 0b0
	LevelH   Level = 0b1
	levelMax Level = LevelH + 1
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
