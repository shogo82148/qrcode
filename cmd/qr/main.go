package main

import (
	"bytes"
	"flag"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/shogo82148/qrcode"
	"github.com/shogo82148/qrcode/microqr"
	"github.com/shogo82148/qrcode/rmqr"
)

func main() {
	var micro, rmqr bool
	var level string
	flag.BoolVar(&micro, "micro", false, "generates Micro QR Code")
	flag.BoolVar(&rmqr, "rmqr", false, "generates rMQR Code")
	flag.StringVar(&level, "level", "", "error correction level")
	flag.Parse()
	filename := flag.Arg(0)

	if !micro && !rmqr {
		encodeQR(level, filename)
	} else if micro {
		encodeMicroQR(level, filename)
	} else if rmqr {
		encodeRMQR(level, filename)
	}
}

func encodeQR(level, filename string) {
	var lv qrcode.Level
	switch level {
	case "l", "L":
		lv = qrcode.LevelL
	case "m", "M":
		lv = qrcode.LevelM
	case "q", "Q", "":
		lv = qrcode.LevelQ
	case "h", "H":
		lv = qrcode.LevelH
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	img, err := qrcode.Encode(data, qrcode.WithLevel(lv))
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}

func encodeMicroQR(level, filename string) {
	var lv microqr.Level
	switch level {
	case "l", "L":
		lv = microqr.LevelL
	case "m", "M":
		lv = microqr.LevelM
	case "q", "Q", "":
		lv = microqr.LevelQ
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	img, err := microqr.Encode(data, microqr.WithLevel(lv))
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}

func encodeRMQR(level, filename string) {
	var lv rmqr.Level
	switch level {
	case "m", "M":
		lv = rmqr.LevelM
	case "h", "H":
		lv = rmqr.LevelH
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	img, err := rmqr.Encode(data, rmqr.WithLevel(lv))
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
