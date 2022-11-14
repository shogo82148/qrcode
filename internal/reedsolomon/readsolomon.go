package reedsolomon

import (
	"io"
)

type Coder interface {
	io.Writer
	Code() []byte
}
