package bitstream

import (
	"bytes"
	"errors"
)

func DecodeNumeric(buf *Buffer, data []byte) error {
	for i := 0; i+2 < len(data); i += 3 {
		bits, err := buf.ReadBits(10)
		if err != nil {
			return err
		}
		if bits >= 1000 {
			return errors.New("invalid digit")
		}
		n1 := bits / 100
		n2 := bits / 10 % 10
		n3 := bits % 10
		data[i+0] = byte(n1 + '0')
		data[i+1] = byte(n2 + '0')
		data[i+2] = byte(n3 + '0')
	}

	switch len(data) % 3 {
	case 1:
		bits, err := buf.ReadBits(4)
		if err != nil {
			return err
		}
		if bits >= 10 {
			return errors.New("invalid digit")
		}
		data[len(data)-1] = byte(bits + '0')
	case 2:
		bits, err := buf.ReadBits(7)
		if err != nil {
			return err
		}
		if bits >= 100 {
			return errors.New("invalid digit")
		}
		n1 := bits / 10
		n2 := bits % 10
		data[len(data)-2] = byte(n1 + '0')
		data[len(data)-1] = byte(n2 + '0')
	}

	return nil
}

func DecodeAlphanumeric(buf *Buffer, data []byte) error {
	for i := 0; i+1 < len(data); i += 2 {
		bits, err := buf.ReadBits(11)
		if err != nil {
			return err
		}
		n1 := int(bits) / 45
		n2 := int(bits) % 45
		if n1 >= 45 {
			return errors.New("invalid digit")
		}
		data[i+0] = bitToAlphanumeric[n1]
		data[i+1] = bitToAlphanumeric[n2]
	}

	if len(data)%2 != 0 {
		bits, err := buf.ReadBits(6)
		if err != nil {
			return err
		}
		if bits >= 45 {
			return errors.New("invalid digit")
		}
		data[len(data)-1] = bitToAlphanumeric[bits]
	}
	return nil
}

func DecodeBytes(buf *Buffer, data []byte) error {
	for i := range data {
		bits, err := buf.ReadBits(8)
		if err != nil {
			return err
		}
		data[i] = byte(bits)
	}
	return nil
}

func DecodeKanji(buf *Buffer, length int) ([]byte, error) {
	var ret bytes.Buffer
	ret.Grow(length * 3)
	for i := 0; i < length; i++ {
		bits, err := buf.ReadBits(13)
		if err != nil {
			return nil, err
		}
		ret.WriteRune(rune(decode[bits]))
	}
	return ret.Bytes(), nil
}
