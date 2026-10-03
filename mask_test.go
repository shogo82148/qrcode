package qrcode

import (
	"testing"

	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
)

func TestEncodeToBitmap_MaskAuto(t *testing.T) {
	inputs := []string{
		"",
		"01234567",
		"HELLO WORLD",
		"https://github.com/shogo82148/qrcode",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	}
	for _, input := range inputs {
		for _, lv := range []Level{LevelL, LevelM, LevelQ, LevelH} {
			qr, err := New([]byte(input), WithLevel(lv))
			if err != nil {
				t.Fatal(err)
			}

			// calculate the penalty scores of all mask patterns.
			want := Mask0
			minPoint := -1
			for mask := Mask0; mask < maskMax; mask++ {
				qr.Mask = mask
				img, err := qr.EncodeToBitmap()
				if err != nil {
					t.Fatal(err)
				}
				point := internalbitmap.Import(img).Point()
				if minPoint < 0 || point < minPoint {
					minPoint = point
					want = mask
				}
			}

			qr.Mask = MaskAuto
			img, err := qr.EncodeToBitmap()
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := DecodeBitmap(img)
			if err != nil {
				t.Fatalf("%q, level %v: failed to decode: %v", input, lv, err)
			}
			if decoded.Mask != want {
				t.Errorf("%q, level %v: got mask %v, want %v", input, lv, decoded.Mask, want)
			}
		}
	}
}
