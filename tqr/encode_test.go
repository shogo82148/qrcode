package tqr

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    []byte("123456789012"),
			wantErr: false,
		},
		{
			name:    "invalid length",
			data:    []byte("123"),
			wantErr: true,
		},
		{
			name:    "non-numeric data",
			data:    []byte("12345678901a"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncodeToBitmap(t *testing.T) {
	data := []byte("000000000000")
	qr, err := New(data)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Errorf("EncodeToBitmap() error = %v", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Errorf("png.Encode() error = %v", err)
	}
	if err := os.WriteFile("tQR.png", buf.Bytes(), 0o644); err != nil {
		t.Errorf("os.WriteFile() error = %v", err)
	}
}
