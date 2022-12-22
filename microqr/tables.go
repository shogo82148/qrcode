package microqr

var formatTable = [5][4]int{
	{}, // dummy

	// M1
	{
		LevelCheck: 0b000,
		LevelL:     -1,
		LevelM:     -1,
		LevelQ:     -1,
	},

	// M2
	{
		LevelCheck: -1,
		LevelL:     0b001,
		LevelM:     0b010,
		LevelQ:     -1,
	},
	// M3
	{
		LevelCheck: -1,
		LevelL:     0b011,
		LevelM:     0b100,
		LevelQ:     -1,
	},
	// M4
	{
		LevelCheck: -1,
		LevelL:     0b101,
		LevelM:     0b110,
		LevelQ:     0b111,
	},
}

var rawFormatTable = [8]struct {
	version Version
	level   Level
}{
	{1, LevelCheck},
	{2, LevelL},
	{2, LevelM},
	{3, LevelL},
	{3, LevelM},
	{4, LevelL},
	{4, LevelM},
	{4, LevelQ},
}

type capacity struct {
	Total      int // number of total code words.
	Data       int // number of data code words.
	DataBits   int // number of bit for data.
	Correction int // number of correction code words.
	MaxError   int // maximum number of code word errors.
	Reserved   int // number of code words reserved for error detection.
}

// X 0510 : 2018
// 表9 マイクロORコード及びORコードの誤り訂正特性
var capacityTable = [5][4]capacity{
	{}, // dummy

	// version 1
	{
		LevelCheck: {
			Total:      5,
			Data:       3,
			DataBits:   20,
			Correction: 2,
			Reserved:   2,
			MaxError:   0,
		},
	},

	// version 2
	{
		LevelL: {
			Total:      10,
			Data:       5,
			DataBits:   40,
			Correction: 5,
			Reserved:   3,
			MaxError:   1,
		},
		LevelM: {
			Total:      10,
			Data:       4,
			DataBits:   32,
			Correction: 6,
			Reserved:   2,
			MaxError:   2,
		},
	},

	// version 3
	{
		LevelL: {
			Total:      17,
			Data:       11,
			DataBits:   84,
			Correction: 6,
			Reserved:   2,
			MaxError:   2,
		},
		LevelM: {
			Total:      17,
			Data:       9,
			DataBits:   68,
			Correction: 8,
			Reserved:   2,
			MaxError:   4,
		},
	},

	// version 4
	{
		LevelL: {
			Total:      24,
			Data:       16,
			DataBits:   128,
			Correction: 8,
			Reserved:   2,
			MaxError:   3,
		},
		LevelM: {
			Total:      24,
			Data:       14,
			DataBits:   112,
			Correction: 10,
			Reserved:   0,
			MaxError:   5,
		},
		LevelQ: {
			Total:      24,
			Data:       10,
			DataBits:   80,
			Correction: 14,
			Reserved:   0,
			MaxError:   7,
		},
	},
}
