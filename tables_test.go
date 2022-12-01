package qrcode

import (
	"testing"
)

func TestTable(t *testing.T) {
	for version := 1; version <= 40; version++ {
		used := usedList[version]
		totalBits := used.Rect.Dx()*used.Rect.Dx() - used.OnesCount()
		total := totalBits / 8

		for level, capacity := range capacityTable[version] {
			if capacity.Total != total {
				t.Errorf("version %d: level %s: unexpected total capacity: got %d, want %d", version, Level(level), capacity.Total, total)
			}
			if capacity.Data+capacity.Correction != capacity.Total {
				t.Errorf("version %d: level %s: number of code words mismatch", version, Level(level))
			}

			var c, k int
			for _, block := range capacity.Blocks {
				if 2*block.MaxError > (block.Total-block.Data)-block.Reserved {
					t.Errorf("version %d: level %s: invalid max error", version, Level(level))
				}
				c += block.Total * block.Num
				k += block.Data * block.Num
			}

			if c != capacity.Total {
				t.Errorf("version %d, level %s: number of total code is unmatched, want %d, got %d", version, Level(level), c, capacity.Total)
			}
			if k != capacity.Data {
				t.Errorf("version %d, level %s: number of data code is unmatched, want %d, got %d", version, Level(level), k, capacity.Data)
			}
		}
	}
}
