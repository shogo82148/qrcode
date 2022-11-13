package reedsolomon

import (
	"io"
	"log"
)

type Coder interface {
	io.Writer
	Code() []byte
}

func New5() Coder {
	return &coder5{}
}

// x^5 + a^113*x^4 + a^164*x^3 + a^166*x^2 + a^119*x + a^10
type coder5 [5]element

func (c *coder5) WriteByte(b byte) error {
	x := c[4]
	c[4].AddMul(c[3], exp(113), x)
	c[3].AddMul(c[2], exp(164), x)
	c[2].AddMul(c[1], exp(166), x)
	c[1].AddMul(c[0], exp(119), x)
	c[0].AddMul(element(b), exp(10), x)
	return nil
}

func (c *coder5) Write(p []byte) (int, error) {
	for _, b := range p {
		c.WriteByte(b)
	}
	return len(p), nil
}

func (c coder5) Code() []byte {
	for i := 0; i < 5; i++ {
		c.WriteByte(0)
	}

	var buf [5]byte
	for i, b := range c {
		buf[i] = byte(b)
	}
	return buf[:]
}

func New10() Coder {
	return &coder10{}
}

// x^10 + a^251*x^9 + a^67*x^8 + a^46*x^7 + a^61*x^6 + a^118*x^5 + a^70*x^4 + a^64*x^3 + a^96*x^2 + a^32*x + a^45
type coder10 [10]element

func (c *coder10) Write(p []byte) (int, error) {
	for _, b := range p {
		x := c[9]
		c[9].AddMul(c[8], 251, x)
		c[8].AddMul(c[7], 67, x)
		c[7].AddMul(c[6], 46, x)
		c[6].AddMul(c[5], 61, x)
		c[5].AddMul(c[4], 118, x)
		c[4].AddMul(c[3], 70, x)
		c[3].AddMul(c[2], 64, x)
		c[2].AddMul(c[1], 96, x)
		c[1].AddMul(c[0], 32, x)
		c[0].AddMul(element(b), 45, x)
	}
	return len(p), nil
}

func (c coder10) Code() []byte {
	var buf [10]byte
	for i, b := range c {
		buf[i] = byte(b)
	}
	return buf[:]
}

func New17() Coder {
	return &coder17{}
}

type coder17 [17 + 9]element

// ⁰¹²³⁴⁵⁶⁷⁸⁹

func (c *coder17) Write(p []byte) (int, error) {
	for i, b := range p {
		c[17+9-1-i] = element(b)
	}
	return len(p), nil
}

func (c coder17) Code() []byte {
	for i := 25; i > 16; i-- {
		x := c[i]
		c[i-0].AddMul(c[i-0], one, x)
		c[i-1].AddMul(c[i-1], exp(43), x)
		c[i-2].AddMul(c[i-2], exp(139), x)
		c[i-3].AddMul(c[i-3], exp(206), x)
		c[i-4].AddMul(c[i-4], exp(78), x)
		c[i-5].AddMul(c[i-5], exp(43), x)
		c[i-6].AddMul(c[i-6], exp(239), x)
		c[i-7].AddMul(c[i-7], exp(123), x)
		c[i-8].AddMul(c[i-8], exp(206), x)
		c[i-9].AddMul(c[i-9], exp(214), x)
		c[i-10].AddMul(c[i-10], exp(147), x)
		c[i-11].AddMul(c[i-11], exp(24), x)
		c[i-12].AddMul(c[i-12], exp(99), x)
		c[i-13].AddMul(c[i-13], exp(150), x)
		c[i-14].AddMul(c[i-14], exp(39), x)
		c[i-15].AddMul(c[i-15], exp(243), x)
		c[i-16].AddMul(c[i-16], exp(163), x)
		c[i-17].AddMul(c[i-17], exp(136), x)
		log.Println(c)
	}

	var buf [17]byte
	for i := range buf {
		buf[i] = byte(c[16-i])
	}
	return buf[:]
}
