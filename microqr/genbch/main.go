package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
)

func main() {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, "// Code generated by genbch/main.go; DO NOT EDIT.")
	fmt.Fprintln(&buf, "")
	fmt.Fprintln(&buf, "package microqr")

	// model number
	fmt.Fprintln(&buf, "var encodedFormat = [...]uint{")
	p := parameter{
		in:    5,
		check: 10,
		mask:  0b100_0100_0100_0101,
		g:     1<<10 + 1<<8 + 1<<5 + 1<<4 + 1<<2 + 1<<1 + 1,
	}
	for i := 0; i < 32; i++ {
		fmt.Fprintf(&buf, "0x%04x, // %05b\n", p.generate(uint(i)), i)
	}
	fmt.Fprintln(&buf, "}")
	fmt.Fprintln(&buf, "")

	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("bch_gen.go", out, 0o644); err != nil {
		log.Fatal(err)
	}
}

type parameter struct {
	in    uint // bit length of input
	check uint // bit length of check bits
	mask  uint // mask
	g     uint // generating polynomial
}

func (p parameter) generate(a uint) uint {
	ret := uint(0)
	b := a
	for i := uint(0); i < p.in+p.check; i++ {
		b <<= 1
		ret = ret<<1 + (b>>p.in)&1
		c := (ret >> p.check) & 1 // carry
		ret ^= c * p.g
	}
	return (a<<p.check + ret) ^ p.mask
}
