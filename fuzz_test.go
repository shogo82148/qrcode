package qrcode

import (
	"bytes"
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
		if lv != LevelL && lv != LevelM && lv != LevelQ && lv != LevelH {
			return
		}
		qr, err := New(lv, data)
		if err != nil {
			return
		}

		got := make([]byte, 0, len(data))
		for _, s := range qr.Segments {
			got = append(got, s.Data...)
		}
		if !bytes.Equal(data, got) {
			t.Errorf("result not match: in %q, out %q", data, got)
		}

		_, err = qr.EncodeToBitmap()
		if err != nil {
			t.Fatal(err)
		}
	})
}
