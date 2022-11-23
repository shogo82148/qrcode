package qrcode

//go:generate go run genbase/main.go
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
	LevelL Level = 0b01
	LevelM Level = 0b00
	LevelQ Level = 0b11
	LevelH Level = 0b10
)

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

type Mask int

const (
	Mask0 Mask = iota
	Mask1
	Mask2
	Mask3
	Mask4
	Mask5
	Mask6
	Mask7
	maskMax

	MaskAuto Mask = -1
)

type Mode uint8

const (
	// ModeECI is ECI(Extended Channel Interpretation) mode.
	ModeECI Mode = 0b0111

	// ModeNumber is number mode.
	// The Data must be ascii characters [0-9].
	ModeNumber Mode = 0b0001

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

type Segment struct {
	Mode Mode
	Data []byte
}
