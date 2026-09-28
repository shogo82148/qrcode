package tqr

import "testing"

func TestBase(t *testing.T) {
	for y := base.Rect.Min.Y; y < base.Rect.Max.Y; y++ {
		for x := base.Rect.Min.X; x < base.Rect.Max.X; x++ {
			if base.BinaryAt(x, y) {
				if !used.BinaryAt(x, y) {
					t.Errorf("Base binary at (%d, %d) is not used", x, y)
				}
			}
		}
	}
}

func TestUsed(t *testing.T) {
	for _, pos := range dataPos {
		x, y := int(pos.x)+2, int(pos.y)+2
		if used.BinaryAt(x, y) {
			t.Errorf("(%d, %d) is marked as used", x, y)
		}
	}
}
