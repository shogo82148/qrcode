package microqr

import (
	"testing"

	internalbitmap "github.com/shogo82148/qrcode/internal/bitmap"
)

func TestEncodeToBitmap_MaskAuto(t *testing.T) {
	inputs := []string{
		"01234",
		"MICROQR",
		"0123456789012345",
		"HELLO WORLD 12345",
	}
	for _, input := range inputs {
		for _, lv := range []Level{LevelL, LevelM, LevelQ} {
			qr, err := New([]byte(input), WithLevel(lv))
			if err != nil {
				// the data is too large for the level.
				continue
			}

			// calculate the evaluation scores of all mask patterns.
			want := Mask0
			maxPoint := -1
			for mask := Mask0; mask < maskMax; mask++ {
				qr.Mask = mask
				img, err := qr.EncodeToBitmap()
				if err != nil {
					t.Fatal(err)
				}
				point := internalbitmap.Import(img).PointMicro()
				if point > maxPoint {
					maxPoint = point
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
