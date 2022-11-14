package bitstream

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestReadBit(t *testing.T) {
	buf := NewBuffer([]byte{
		0b1000_0000, 0b1000_0001,
	})

	got := make([]byte, 0, 16)
	for {
		bit, err := buf.ReadBit()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, bit)
	}

	want := []byte{
		1, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 1,
	}

	if !bytes.Equal(got, want) {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestWriteBit(t *testing.T) {
	in := []byte{
		1, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 1,
	}
	var buf Buffer
	for _, b := range in {
		if err := buf.WriteBit(b); err != nil {
			t.Fatal(err)
		}
	}
	got := buf.Bytes()
	want := []byte{
		0b1000_0000, 0b1000_0001,
	}

	if !bytes.Equal(got, want) {
		t.Errorf("got %x, want %x", got, want)
	}
}
