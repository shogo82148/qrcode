package main

import (
	"bytes"
	"image/png"
	"log"
	"os"

	"github.com/shogo82148/qrcode"
)

func main() {
	qr := &qrcode.QRCode{
		Version: 10,
		Level:   qrcode.LevelH,
		Mask:    0b100,
		Segments: []qrcode.Segment{
			{
				Mode: qrcode.ModeAlphanumeric,
				Data: []byte("VERSION 10 QR CODE"),
			},
			{
				Mode: qrcode.ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: qrcode.ModeAlphanumeric,
				Data: []byte(" UP TO 174 CHAR AT H LEVEL"),
			},
			{
				Mode: qrcode.ModeBytes,
				Data: []byte(","),
			},
			{
				Mode: qrcode.ModeAlphanumeric,
				Data: []byte(
					" WITH 57X57 MODULES AND PLENTY OF ERROR CORRECTION TO GO AROUND. " +
						"NOTE THAT THERE ARE ADDITIONAL TRACKING BOXES"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("qrcode.png", buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
