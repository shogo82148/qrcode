package bitstream

import (
	"io"
	"strconv"
)

type Buffer struct {
	buf    []byte
	offset int // read at buf[off], write at buf[len(buf)]
	read   int // read number of bit at buf[off]
	wrote  int
}

// NewBuffer returns a buffer for bits sequence data.
// The bits sequence starts from MSB of the first byte.
func NewBuffer(data []byte) *Buffer {
	return &Buffer{buf: data}
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Len returns the number of bits of the buffer.
func (b *Buffer) Len() int {
	l := len(b.buf) * 8
	if b.wrote != 0 {
		l = l - 8 + b.wrote
	}
	return l
}

// ReadBit reads one bit from b.
func (b *Buffer) ReadBit() (uint8, error) {
	if b.offset >= len(b.buf) {
		return 0, io.EOF
	}
	bit := (b.buf[b.offset] >> (7 - b.read)) & 1
	b.read++
	if b.read >= 8 {
		b.offset++
		b.read = 0
	}
	return bit, nil
}

func (b *Buffer) ReadBits(n int) (uint64, error) {
	if n > 64 {
		panic("too long bit length: " + strconv.Itoa(n))
	}
	if b.offset >= len(b.buf) {
		return 0, io.EOF
	}

	var ret uint64
	for i := 0; i < n; i++ {
		bit, err := b.ReadBit()
		if err != nil {
			ret <<= n - i
			return ret, nil
		}
		ret = (ret << 1) | uint64(bit)
	}
	return ret, nil
}

// WriteBit writes one bit to b.
func (b *Buffer) WriteBit(bit uint8) error {
	bit &= 1
	if b.wrote == 0 {
		b.buf = append(b.buf, bit<<7)
		b.wrote = 1
		return nil
	}
	b.buf[len(b.buf)-1] |= bit << (7 - b.wrote)
	b.wrote = (b.wrote + 1) % 8
	return nil
}

// WriteBitsLSB writes n bits of LSB to b.
// if n > 64, it panics.
func (b *Buffer) WriteBitsLSB(bits uint64, n int) error {
	switch {
	case n > 64:
		panic("too long bit length: " + strconv.Itoa(n))
	case n > 56:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits>>56), n-56)
		b.writeByte(uint8(bits >> 48))
		b.writeByte(uint8(bits >> 40))
		b.writeByte(uint8(bits >> 32))
		b.writeByte(uint8(bits >> 24))
		b.writeByte(uint8(bits >> 16))
		b.writeByte(uint8(bits >> 8))
		b.writeByte(uint8(bits))
	case n > 48:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits>>48), n-48)
		b.writeByte(uint8(bits >> 40))
		b.writeByte(uint8(bits >> 32))
		b.writeByte(uint8(bits >> 24))
		b.writeByte(uint8(bits >> 16))
		b.writeByte(uint8(bits >> 8))
		b.writeByte(uint8(bits))
	case n > 40:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits>>40), n-40)
		b.writeByte(uint8(bits >> 32))
		b.writeByte(uint8(bits >> 24))
		b.writeByte(uint8(bits >> 16))
		b.writeByte(uint8(bits >> 8))
		b.writeByte(uint8(bits))
	case n > 32:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits>>32), n-32)
		b.writeByte(uint8(bits >> 24))
		b.writeByte(uint8(bits >> 16))
		b.writeByte(uint8(bits >> 8))
		b.writeByte(uint8(bits))
	case n > 24:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits>>24), n-24)
		b.writeByte(uint8(bits >> 16))
		b.writeByte(uint8(bits >> 8))
		b.writeByte(uint8(bits))
	case n > 16:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits>>16), n-16)
		b.writeByte(uint8(bits >> 8))
		b.writeByte(uint8(bits))
	case n > 8:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits>>8), n-8)
		b.writeByte(uint8(bits))
	case n > 0:
		bits &= (1 << n) - 1
		b.writeBitsLSB(uint8(bits), n)
	case n == 0:
		// nothing to do
	default:
		panic("negative bit length")
	}
	return nil
}

func (b *Buffer) writeBitsLSB(bits uint8, n int) {
	if b.wrote == 0 {
		b.buf = append(b.buf, byte(bits<<(8-n)))
		b.wrote = n % 8
		return
	}
	if m := b.wrote + n; m > 8 {
		b.buf[len(b.buf)-1] |= bits >> (m - 8)
		b.wrote = m - 8
		b.buf = append(b.buf, bits<<(8-b.wrote))
		return
	}
	b.buf[len(b.buf)-1] |= bits << (8 - (b.wrote + n))
	b.wrote += n
}

func (b *Buffer) writeByte(bits uint8) {
	if b.wrote == 0 {
		b.buf = append(b.buf, bits)
		return
	}
	b.buf[len(b.buf)-1] |= bits >> b.wrote
	b.buf = append(b.buf, bits<<(8-b.wrote))
}
