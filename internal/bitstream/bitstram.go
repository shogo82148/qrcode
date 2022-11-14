package bitstream

import "io"

type Buffer struct {
	buf    []byte
	offset int // read at buf[off], write at buf[len(buf)]
	read   int // read number of bit at buf[off]
	wrote  int
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{buf: data}
}

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

func (b *Buffer) WriteBit(bit uint8) error {
	if b.wrote == 0 {
		b.buf = append(b.buf, byte(bit<<7))
		b.wrote = 1
		return nil
	}
	b.buf[len(b.buf)-1] |= bit << (7 - b.wrote)
	b.wrote = (b.wrote + 1) % 8
	return nil
}
