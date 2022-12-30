package microqr

import (
	"bytes"
	"reflect"
	"testing"
)

func FuzzNew(f *testing.F) {
	f.Add(int(LevelL), []byte("MICROQR"))
	f.Add(int(LevelL), []byte("12345678901234567890123456789012345"))
	f.Add(int(LevelL), []byte("123456789012345678901234"))
	f.Add(int(LevelL), []byte("12345678901234567890123"))
	f.Add(int(LevelM), []byte("1haicso"))
	f.Fuzz(func(t *testing.T, level int, data []byte) {
		lv := Level(level)
		qr0, err := New(data, WithLevel(lv))
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
	f.Add(int(LevelL), []byte("MICROQR"))
	f.Add(int(LevelL), []byte("12345678901234567890123456789012345"))
	f.Add(int(LevelL), []byte("123456789012345678901234"))
	f.Add(int(LevelL), []byte("12345678901234567890123"))
	f.Add(int(LevelM), []byte("1haicso"))
	f.Add(int(LevelM), []byte("点"))
	f.Fuzz(func(t *testing.T, level int, data []byte) {
		lv := Level(level)
		if lv != LevelCheck && lv != LevelL && lv != LevelM && lv != LevelQ {
			return
		}
		qr0, err := newFromKanji(lv, data)
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
