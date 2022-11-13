package reedsolomon

import (
	"io"
	"log"
)

type Coder interface {
	io.Writer
	Code() []byte
}

// func New5() Coder {
// 	return &coder5{}
// }

// // x^5 + a^113*x^4 + a^164*x^3 + a^166*x^2 + a^119*x + a^10
// type coder5 [5]element

// func (c *coder5) WriteByte(b byte) error {
// 	x := c[4]
// 	c[4].AddMulExp(c[3], exp(113), x)
// 	c[3].AddMulExp(c[2], exp(164), x)
// 	c[2].AddMulExp(c[1], exp(166), x)
// 	c[1].AddMulExp(c[0], exp(119), x)
// 	c[0].AddMulExp(element(b), exp(10), x)
// 	return nil
// }

// func (c *coder5) Write(p []byte) (int, error) {
// 	for _, b := range p {
// 		c.WriteByte(b)
// 	}
// 	return len(p), nil
// }

// func (c coder5) Code() []byte {
// 	for i := 0; i < 5; i++ {
// 		c.WriteByte(0)
// 	}

// 	var buf [5]byte
// 	for i, b := range c {
// 		buf[i] = byte(b)
// 	}
// 	return buf[:]
// }

// func New10() Coder {
// 	return &coder10{}
// }

// // x^10 + a^251*x^9 + a^67*x^8 + a^46*x^7 + a^61*x^6 + a^118*x^5 + a^70*x^4 + a^64*x^3 + a^96*x^2 + a^32*x + a^45
// type coder10 [10]element

// func (c *coder10) Write(p []byte) (int, error) {
// 	for _, b := range p {
// 		x := c[9]
// 		c[9].AddMulExp(c[8], 251, x)
// 		c[8].AddMulExp(c[7], 67, x)
// 		c[7].AddMulExp(c[6], 46, x)
// 		c[6].AddMulExp(c[5], 61, x)
// 		c[5].AddMulExp(c[4], 118, x)
// 		c[4].AddMulExp(c[3], 70, x)
// 		c[3].AddMulExp(c[2], 64, x)
// 		c[2].AddMulExp(c[1], 96, x)
// 		c[1].AddMulExp(c[0], 32, x)
// 		c[0].AddMulExp(element(b), 45, x)
// 	}
// 	return len(p), nil
// }

// func (c coder10) Code() []byte {
// 	var buf [10]byte
// 	for i, b := range c {
// 		buf[i] = byte(b)
// 	}
// 	return buf[:]
// }

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
		c[i-0].AddMulExp(c[i-0], x, 0)
		c[i-1].AddMulExp(c[i-1], x, 43)
		c[i-2].AddMulExp(c[i-2], x, 139)
		c[i-3].AddMulExp(c[i-3], x, 206)
		c[i-4].AddMulExp(c[i-4], x, 78)
		c[i-5].AddMulExp(c[i-5], x, 43)
		c[i-6].AddMulExp(c[i-6], x, 239)
		c[i-7].AddMulExp(c[i-7], x, 123)
		c[i-8].AddMulExp(c[i-8], x, 206)
		c[i-9].AddMulExp(c[i-9], x, 214)
		c[i-10].AddMulExp(c[i-10], x, 147)
		c[i-11].AddMulExp(c[i-11], x, 24)
		c[i-12].AddMulExp(c[i-12], x, 99)
		c[i-13].AddMulExp(c[i-13], x, 150)
		c[i-14].AddMulExp(c[i-14], x, 39)
		c[i-15].AddMulExp(c[i-15], x, 243)
		c[i-16].AddMulExp(c[i-16], x, 163)
		c[i-17].AddMulExp(c[i-17], x, 136)
		log.Println(c)
	}

	var buf [17]byte
	for i := range buf {
		buf[i] = byte(c[16-i])
	}
	return buf[:]
}
