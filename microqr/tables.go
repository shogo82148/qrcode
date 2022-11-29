package microqr

type capacity struct {
	Total int // number of total code words
	Data  int // number of data code words
}

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

// X 0510 : 2018
// 表9 マイクロORコード及びORコードの誤り訂正特性
var capacityTable = [5][4]capacity{
	{}, // dummy

	// version 1
	{
		LevelCheck: {
			Total: 5,
			Data:  3,
		},
	},

	// version 2
	{
		LevelL: {
			Total: 10,
			Data:  5,
		},
		LevelM: {
			Total: 10,
			Data:  4,
		},
	},

	// version 3
	{
		LevelL: {
			Total: 17,
			Data:  11,
		},
		LevelM: {
			Total: 17,
			Data:  9,
		},
	},

	// version 4
	{
		LevelL: {
			Total: 24,
			Data:  16,
		},
		LevelM: {
			Total: 24,
			Data:  14,
		},
		LevelQ: {
			Total: 24,
			Data:  10,
		},
	},
}
