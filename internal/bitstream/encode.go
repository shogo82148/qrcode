package bitstream

import (
	"fmt"
)

func IsNumeric(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func IsAlphanumeric(ch byte) bool {
	return alphabets[ch] >= 0
}

func EncodeNumeric(buf *Buffer, data []byte) error {
	// validate
	for _, ch := range data {
		if !IsNumeric(ch) {
			return fmt.Errorf("qrcode: invalid character in number mode: %02x", ch)
		}
	}

	// encode
	for i := 0; i+2 < len(data); i += 3 {
		n1 := int(data[i+0] - '0')
		n2 := int(data[i+1] - '0')
		n3 := int(data[i+2] - '0')
		bits := n1*100 + n2*10 + n3
		buf.WriteBitsLSB(uint64(bits), 10)
	}

	switch len(data) % 3 {
	case 1:
		bits := data[len(data)-1] - '0'
		buf.WriteBitsLSB(uint64(bits), 4)
	case 2:
		n1 := int(data[len(data)-2] - '0')
		n2 := int(data[len(data)-1] - '0')
		bits := n1*10 + n2
		buf.WriteBitsLSB(uint64(bits), 7)
	}

	return nil
}

var bitToAlphanumeric = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:")
var alphabets [256]int

func init() {
	for i := range alphabets {
		alphabets[i] = -1
	}
	for i, ch := range bitToAlphanumeric {
		alphabets[ch] = i
	}
}

func EncodeAlphanumeric(buf *Buffer, data []byte) error {
	for _, ch := range data {
		if !IsAlphanumeric(ch) {
			return fmt.Errorf("qrcode: invalid character in alphabet mode: %02x", ch)
		}
	}

	// data
	for i := 0; i+1 < len(data); i += 2 {
		n1 := alphabets[data[i+0]]
		n2 := alphabets[data[i+1]]
		bits := n1*45 + n2
		buf.WriteBitsLSB(uint64(bits), 11)
	}

	if len(data)%2 != 0 {
		bits := alphabets[data[len(data)-1]]
		buf.WriteBitsLSB(uint64(bits), 6)
	}

	return nil
}

func EncodeBytes(buf *Buffer, data []byte) error {
	for _, bits := range data {
		buf.WriteBitsLSB(uint64(bits), 8)
	}

	return nil
}
