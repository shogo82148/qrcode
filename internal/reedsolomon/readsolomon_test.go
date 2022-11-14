package reedsolomon

import (
	"bytes"
	"testing"
)

// JIS X 0510: 2018
// 附属書1
// シンボルの符号化例
// func TestCoder10(t *testing.T) {
// 	w := New10()
// 	w.Write([]byte{
// 		0b0001_0000, 0b0010_0000, 0b0000_1100, 0b0101_0110,
// 		0b0110_0001, 0b1000_0000,

// 		0b1110_1100, 0b0001_0001,
// 		0b1110_1100, 0b0001_0001,
// 		0b1110_1100, 0b0001_0001,
// 		0b1110_1100, 0b0001_0001,
// 		0b1110_1100, 0b0001_0001,
// 	})
// 	want := []byte{
// 		0b1010_0101, 0b0010_0100, 0b1101_0100, 0b1100_0001,
// 		0b1110_1101, 0b0011_0110, 0b1100_0111, 0b1000_0111,
// 		0b0010_1100, 0b0101_0101,
// 	}
// 	got := w.Code()
// 	if !bytes.Equal(got, want) {
// 		t.Errorf("got %x, want %#x", got, want)
// 	}
// }

// func TestCoder5(t *testing.T) {
// 	w := New5()
// 	w.Write([]byte{
// 		0b0100_0000, 0b0001_1000, 0b1010_1100, 0b1100_0011,
// 		0b0000_0000,
// 	})
// 	want := []byte{
// 		0b1000_0110, 0b0000_1101, 0b0010_0010, 0b1010_1110,
// 		0b0011_0000,
// 	}
// 	got := w.Code()
// 	if !bytes.Equal(got, want) {
// 		t.Errorf("got %x, want %#x", got, want)
// 	}
// }

func TestSample(t *testing.T) {
	// https://www.swetake.com/qrcode/qr3.html
	// https://www.swetake.com/qrcode/qr_ecc_calc_sample.html
	w := &coder17{}
	w.Write([]byte{
		0b0010_0000, 0b0100_0001, 0b1100_1101, 0b0100_0101,
		0b0010_1001, 0b1101_1100, 0b0010_1110, 0b1000_0000,
		0b1110_1100,
	})
	want := []byte{
		42, 159, 74, 221, 244, 169, 239, 150, 138, 70, 237, 85, 224, 96, 74, 219, 61,
	}
	got := w.Code()
	if !bytes.Equal(got, want) {
		t.Errorf("got %x, want %x", got, want)
	}
}
