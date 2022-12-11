// Package rmqr handles rMRQ Codes.
package rmqr

import "github.com/shogo82148/qrcode/bitmap"

func DecodeBitmap(img *bitmap.Image) (*QRCode, error) {
	return &QRCode{}, nil
}
