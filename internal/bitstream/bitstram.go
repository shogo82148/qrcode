package bitstream

type Buffer struct {
	buf    []byte
	offset int // read at buf[off], write at buf[len(buf)]
	read   int // read number of bit at buf[off]
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{buf: data}
}

func (b *Buffer) ReadBit() uint8 {
	bit := (b.buf[b.offset] >> (7 - b.read)) & 1
	b.read++
	if b.read >= 8 {
		b.offset++
		b.read = 0
	}
	return bit
}
