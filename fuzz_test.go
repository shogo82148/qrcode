package qrcode

import (
	"bytes"
	"reflect"
	"testing"
)

func FuzzNew(f *testing.F) {
	f.Add(int(LevelL), []byte("01234567"))
	f.Add(int(LevelM), []byte("01234567"))
	f.Add(int(LevelH), []byte("01234567"))
	f.Add(int(LevelH), []byte("Ver1"))
	f.Add(int(LevelH), []byte("VERSION 10 QR CODE, UP TO 174 CHAR AT H LEVEL, WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND. NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"))
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
	f.Add(int(LevelL), []byte("01234567"))
	f.Add(int(LevelM), []byte("01234567"))
	f.Add(int(LevelH), []byte("01234567"))
	f.Add(int(LevelH), []byte("Ver1"))
	f.Add(int(LevelH), []byte("VERSION 10 QR CODE, UP TO 174 CHAR AT H LEVEL, WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND. NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"))
	f.Add(int(LevelH), []byte("点茗"))
	f.Add(int(LevelH), []byte("点"))
	f.Fuzz(func(t *testing.T, level int, data []byte) {
		lv := Level(level)
		if lv != LevelL && lv != LevelM && lv != LevelQ && lv != LevelH {
			return
		}
		qr0, err := New(data, WithLevel(lv), WithKanji(true))
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
