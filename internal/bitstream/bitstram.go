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
// if n > 8, it panics.
func (b *Buffer) WriteBitsLSB(bits uint8, n int) error {
	if n > 8 {
		panic("too long bit length" + strconv.Itoa(n))
	}
	if b.wrote == 0 {
		b.buf = append(b.buf, byte(bits<<(8-n)))
		b.wrote = n
		return nil
	}
	if m := b.wrote + n; m > 8 {
		b.buf[len(b.buf)-1] |= bits >> (m - 8)
		b.wrote = m - 8
		b.buf = append(b.buf, byte(bits<<(8-b.wrote)))
		return nil
	}
	b.buf[len(b.buf)-1] |= bits << (8 - (b.wrote + n))
	b.wrote += n
	return nil
}
