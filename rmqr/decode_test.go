package rmqr

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/shogo82148/qrcode/bitmap"
)

func TestDecode1(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/rmqr2.png
	r, err := os.Open("testdata/rmqr2.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 59, 15))
	for y := 0; y <= 15; y++ {
		for x := 0; x <= 59; x++ {
			X := float64(x)*(164/59.0) + 9
			Y := float64(y)*(46/15.0) + 9
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
}

func TestDecode2(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r7x43.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 43, 7))
	for y := 0; y < 7; y++ {
		for x := 0; x < 43; x++ {
			X := float64(x)*(335/43.0) + 29
			Y := float64(y)*(50/7.0) + 30
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode3(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r7x139.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 139, 7))
	for y := 0; y < 7; y++ {
		for x := 0; x < 139; x++ {
			X := float64(x)*(386/139.0) + 5
			Y := float64(y)*(20/7.0) + 7
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode4(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r9x43.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 43, 9))
	for y := 0; y < 9; y++ {
		for x := 0; x < 43; x++ {
			X := float64(x)*(120/43.0) + 7
			Y := float64(y)*(25/9.0) + 7
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode5(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r9x139.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 139, 9))
	for y := 0; y < 9; y++ {
		for x := 0; x < 139; x++ {
			X := float64(x)*(385.5/139.0) + 7
			Y := float64(y)*(25/9.0) + 7
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode6(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r11x27.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 27, 11))
	for y := 0; y < 11; y++ {
		for x := 0; x < 27; x++ {
			X := float64(x)*(75/27.0) + 6
			Y := float64(y)*(30/11.0) + 7
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode7(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r11x139.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 139, 11))
	for y := 0; y < 11; y++ {
		for x := 0; x < 139; x++ {
			X := float64(x)*(75/27.0) + 7
			Y := float64(y)*(30/11.0) + 7
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode8(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r13x27.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 27, 13))
	for y := 0; y < 13; y++ {
		for x := 0; x < 27; x++ {
			X := float64(x)*(74/27.0) + 9
			Y := float64(y)*(35/13.0) + 9
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode9(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r15x43.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 43, 15))
	for y := 0; y < 15; y++ {
		for x := 0; x < 43; x++ {
			X := float64(x)*(120/43.0) + 9
			Y := float64(y)*(40/15.0) + 9
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode10(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r15x139.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 139, 15))
	for y := 0; y < 15; y++ {
		for x := 0; x < 139; x++ {
			X := float64(x)*(390/139.0) + 7
			Y := float64(y)*(40/15.0) + 8
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode11(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r17x43.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 43, 17))
	for y := 0; y < 17; y++ {
		for x := 0; x < 43; x++ {
			X := float64(x)*(120/43.0) + 7
			Y := float64(y)*(46/17.0) + 8
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode12(t *testing.T) {
	// from https://www.qrcode.com/img/rmqr/graph.jpg
	r, err := os.Open("testdata/r17x139.png")
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	binimg := bitmap.New(image.Rect(0, 0, 139, 17))
	for y := 0; y < 17; y++ {
		for x := 0; x < 139; x++ {
			X := float64(x)*(386/139.0) + 7
			Y := float64(y)*(45/17.0) + 9
			binimg.Set(x, y, img.At(round(X), round(Y)))
		}
	}
	_, err = DecodeBitmap(binimg)
	if err != nil {
		t.Fatal(err)
	}
}
