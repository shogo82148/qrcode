package rmqr

import (
	"bytes"
	"reflect"
	"testing"
)

func FuzzNew(f *testing.F) {
	f.Add(int(LevelM), []byte("01234567"))
	f.Add(int(LevelH), []byte("01234567"))
	f.Add(int(LevelM), []byte("123456789012"))
	f.Fuzz(func(t *testing.T, level int, data []byte) {
		lv := Level(level)
		if lv != LevelM && lv != LevelH {
			return
		}
		qr0, err := New(lv, data)
		if err != nil {
			return
		}

		// check the result
		got := make([]byte, 0, len(data))
		for _, s := range qr0.Segments {
			got = append(got, s.Data...)
		}
		if !bytes.Equal(data, got) {
			t.Errorf("result not match: in %q, out %q", data, got)
		}

		// encode and decode
		img, err := qr0.EncodeToBitmap()
		if err != nil {
			t.Fatal(err)
		}
		qr1, err := DecodeBitmap(img)
		if err != nil {
			t.Fatal(err)
		}
		if (len(qr0.Segments) != 0 || len(qr1.Segments) != 0) && !reflect.DeepEqual(qr0.Segments, qr1.Segments) {
			t.Errorf("decoded result not match: input %#v, output %#v", qr0.Segments, qr1.Segments)
		}
	})
}

func FuzzNewFromKanji(f *testing.F) {
	f.Add(int(LevelM), []byte("01234567"))
	f.Add(int(LevelH), []byte("01234567"))
	f.Add(int(LevelM), []byte("123456789012"))
	f.Add(int(LevelM), []byte("点"))
	f.Fuzz(func(t *testing.T, level int, data []byte) {
		lv := Level(level)
		if lv != LevelM && lv != LevelH {
			return
		}
		qr0, err := NewFromKanji(lv, data)
		if err != nil {
			return
		}

		// check the result
		got := make([]byte, 0, len(data))
		for _, s := range qr0.Segments {
			got = append(got, s.Data...)
		}
		if !bytes.Equal(data, got) {
			t.Errorf("result not match: in %q, out %q", data, got)
		}

		// encode and decode
		img, err := qr0.EncodeToBitmap()
		if err != nil {
			t.Fatal(err)
		}
		_, err = DecodeBitmap(img)
		if err != nil {
			t.Fatal(err)
		}
	})
}
