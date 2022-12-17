package rmqr

import (
	"testing"
)

func TestCapacityTable(t *testing.T) {
	for version, capacity := range capacityTable {
		used := usedList[version]
		totalBits := used.Rect.Dx()*used.Rect.Dy() - used.OnesCount()
		total := totalBits / 8

		if capacity[LevelM].BitLength != capacity[LevelH].BitLength {
			t.Errorf("Version %s: bit length not match", Version(version))
		}
		for level, capacity := range capacity {
			if capacity.Total != total {
				t.Errorf("version %s: level %s: unexpected total capacity: got %d, want %d",
					Version(version), Level(level), capacity.Total, total)
			}
			if capacity.Correction+capacity.Data != capacity.Total {
				t.Errorf("version %s: level %s: total is not match", Version(version), Level(level))
			}

			var c, k int
			for _, block := range capacity.Blocks {
				if 2*block.MaxError > (block.Total-block.Data)-block.Reserved {
					t.Errorf("version %s: level %s: invalid max error", Version(version), Level(level))
				}
				c += block.Total * block.Num
				k += block.Data * block.Num
			}
			if c != capacity.Total {
				t.Errorf("version %s, level %s: number of total code is unmatched, want %d, got %d", Version(version), Level(level), c, capacity.Total)
			}
			if k != capacity.Data {
				t.Errorf("version %s, level %s: number of data code is unmatched, want %d, got %d", Version(version), Level(level), k, capacity.Data)
			}
		}
	}
}
