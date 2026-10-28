package rmqr

import (
	"bytes"
	"runtime"
	"testing"
)

func TestNewDataTooLarge(t *testing.T) {
	for _, kanji := range []bool{false, true} {
		// the largest numeric data that fits in the largest version.
		data := bytes.Repeat([]byte("1"), 361)
		if _, err := New(data, WithLevel(LevelM), WithKanji(kanji)); err != nil {
			t.Errorf("kanji = %v: unexpected error: %v", kanji, err)
		}

		data = bytes.Repeat([]byte("1"), maxDataLength(LevelM)+1)
		if _, err := New(data, WithLevel(LevelM), WithKanji(kanji)); err == nil {
			t.Errorf("kanji = %v: want error, got nil", kanji)
		}
	}
}

func TestNewDataTooLargeAllocation(t *testing.T) {
	data := bytes.Repeat([]byte("a"), 1<<20)
	for _, kanji := range []bool{false, true} {
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		_, err := New(data, WithLevel(LevelM), WithKanji(kanji))
		runtime.ReadMemStats(&after)
		if err == nil {
			t.Errorf("kanji = %v: want error, got nil", kanji)
		}
		if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 1<<16 {
			t.Errorf("kanji = %v: New allocated %d bytes for too large data", kanji, allocated)
		}
	}
}
