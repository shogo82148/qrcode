package poly

import (
	"testing"

	"github.com/shogo82148/qrcode/internal/reedsolomon/element"
)

func TestDegree(t *testing.T) {
	p := Poly{}
	if p.Degree() != 0 {
		t.Errorf("got %d, want %d", p.Degree(), 0)
	}

	p = Poly{element.Zero, element.One}
	if p.Degree() != 0 {
		t.Errorf("got %d, want %d", p.Degree(), 0)
	}

	p = Poly{element.One, element.One}
	if p.Degree() != 1 {
		t.Errorf("got %d, want %d", p.Degree(), 1)
	}
}
