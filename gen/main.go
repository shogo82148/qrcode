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
		Version: 3,
		Level:   qrcode.LevelH,
		Mask:    0b001,
		Segments: []qrcode.Segment{
			{
				Mode: qrcode.ModeBytes,
				Data: []byte("Version 3 QR Code"),
			},
		},
	}
	img, err := qr.Encode()
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
