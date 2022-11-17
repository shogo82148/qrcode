package qrcode

type capacity struct {
	Total  int // number of total code words
	Data   int // number of data code words
	Blocks []blockCapacity
}

type blockCapacity struct {
	Num   int // number of blocks
	Total int // number of total code words
	Data  int // number of data code words
}

// X 0510 : 2018
// 表9 マイクロORコード及びORコードの誤り訂正特性
var capacityTable = [41][4]capacity{
	{}, // dummy

	// version 1
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 2
	{
		LevelL: {
			Total: 44,
			Data:  34,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 34},
			},
		},
		LevelM: {
			Total: 44,
			Data:  28,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 28},
			},
		},
		LevelQ: {
			Total: 44,
			Data:  22,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 22},
			},
		},
		LevelH: {
			Total: 44,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 16},
			},
		},
	},

	// version 3
	{
		LevelL: {
			Total: 70,
			Data:  55,
			Blocks: []blockCapacity{
				{Num: 1, Total: 70, Data: 55},
			},
		},
		LevelM: {
			Total: 70,
			Data:  44,
			Blocks: []blockCapacity{
				{Num: 1, Total: 70, Data: 44},
			},
		},
		LevelQ: {
			Total: 70,
			Data:  34,
			Blocks: []blockCapacity{
				{Num: 2, Total: 35, Data: 17},
			},
		},
		LevelH: {
			Total: 70,
			Data:  26,
			Blocks: []blockCapacity{
				{Num: 2, Total: 35, Data: 13},
			},
		},
	},

	// version 4
	{
		LevelL: {
			Total: 100,
			Data:  80,
			Blocks: []blockCapacity{
				{Num: 1, Total: 100, Data: 80},
			},
		},
		LevelM: {
			Total: 100,
			Data:  64,
			Blocks: []blockCapacity{
				{Num: 2, Total: 50, Data: 32},
			},
		},
		LevelQ: {
			Total: 100,
			Data:  48,
			Blocks: []blockCapacity{
				{Num: 2, Total: 50, Data: 24},
			},
		},
		LevelH: {
			Total: 100,
			Data:  36,
			Blocks: []blockCapacity{
				{Num: 4, Total: 25, Data: 9},
			},
		},
	},

	// version 5
	{
		LevelL: {
			Total: 134,
			Data:  108,
			Blocks: []blockCapacity{
				{Num: 1, Total: 134, Data: 108},
			},
		},
		LevelM: {
			Total: 134,
			Data:  86,
			Blocks: []blockCapacity{
				{Num: 2, Total: 67, Data: 43},
			},
		},
		LevelQ: {
			Total: 134,
			Data:  62,
			Blocks: []blockCapacity{
				{Num: 2, Total: 33, Data: 15},
				{Num: 2, Total: 34, Data: 16},
			},
		},
		LevelH: {
			Total: 134,
			Data:  46,
			Blocks: []blockCapacity{
				{Num: 2, Total: 33, Data: 11},
				{Num: 2, Total: 34, Data: 12},
			},
		},
	},

	// version 6
	{
		LevelL: {
			Total: 172,
			Data:  136,
			Blocks: []blockCapacity{
				{Num: 2, Total: 86, Data: 68},
			},
		},
		LevelM: {
			Total: 172,
			Data:  108,
			Blocks: []blockCapacity{
				{Num: 4, Total: 43, Data: 27},
			},
		},
		LevelQ: {
			Total: 172,
			Data:  76,
			Blocks: []blockCapacity{
				{Num: 4, Total: 43, Data: 19},
			},
		},
		LevelH: {
			Total: 172,
			Data:  60,
			Blocks: []blockCapacity{
				{Num: 4, Total: 43, Data: 15},
			},
		},
	},

	// version 7
	{
		LevelL: {
			Total: 196,
			Data:  156,
			Blocks: []blockCapacity{
				{Num: 2, Total: 98, Data: 78},
			},
		},
		LevelM: {
			Total: 196,
			Data:  124,
			Blocks: []blockCapacity{
				{Num: 4, Total: 49, Data: 31},
			},
		},
		LevelQ: {
			Total: 196,
			Data:  88,
			Blocks: []blockCapacity{
				{Num: 2, Total: 32, Data: 14},
				{Num: 4, Total: 33, Data: 15},
			},
		},
		LevelH: {
			Total: 196,
			Data:  66,
			Blocks: []blockCapacity{
				{Num: 4, Total: 39, Data: 13},
				{Num: 1, Total: 40, Data: 14},
			},
		},
	},

	// version 8
	{
		LevelL: {
			Total: 242,
			Data:  194,
			Blocks: []blockCapacity{
				{Num: 2, Total: 121, Data: 97},
			},
		},
		LevelM: {
			Total: 242,
			Data:  154,
			Blocks: []blockCapacity{
				{Num: 2, Total: 60, Data: 38},
				{Num: 2, Total: 61, Data: 39},
			},
		},
		LevelQ: {
			Total: 242,
			Data:  110,
			Blocks: []blockCapacity{
				{Num: 4, Total: 40, Data: 18},
				{Num: 2, Total: 41, Data: 19},
			},
		},
		LevelH: {
			Total: 242,
			Data:  86,
			Blocks: []blockCapacity{
				{Num: 4, Total: 40, Data: 14},
				{Num: 2, Total: 41, Data: 15},
			},
		},
	},

	// version 9
	{
		LevelL: {
			Total: 292,
			Data:  232,
			Blocks: []blockCapacity{
				{Num: 2, Total: 146, Data: 116},
			},
		},
		LevelM: {
			Total: 292,
			Data:  182,
			Blocks: []blockCapacity{
				{Num: 3, Total: 58, Data: 36},
				{Num: 2, Total: 59, Data: 37},
			},
		},
		LevelQ: {
			Total: 292,
			Data:  132,
			Blocks: []blockCapacity{
				{Num: 4, Total: 36, Data: 16},
				{Num: 4, Total: 37, Data: 17},
			},
		},
		LevelH: {
			Total: 292,
			Data:  100,
			Blocks: []blockCapacity{
				{Num: 4, Total: 36, Data: 12},
				{Num: 4, Total: 37, Data: 13},
			},
		},
	},

	// version 10
	{
		LevelL: {
			Total: 346,
			Data:  274,
			Blocks: []blockCapacity{
				{Num: 2, Total: 86, Data: 68},
				{Num: 2, Total: 87, Data: 69},
			},
		},
		LevelM: {
			Total: 346,
			Data:  216,
			Blocks: []blockCapacity{
				{Num: 4, Total: 69, Data: 43},
				{Num: 1, Total: 70, Data: 44},
			},
		},
		LevelQ: {
			Total: 346,
			Data:  154,
			Blocks: []blockCapacity{
				{Num: 6, Total: 43, Data: 19},
				{Num: 2, Total: 44, Data: 20},
			},
		},
		LevelH: {
			Total: 346,
			Data:  122,
			Blocks: []blockCapacity{
				{Num: 6, Total: 43, Data: 15},
				{Num: 2, Total: 44, Data: 16},
			},
		},
	},

	// version 11
	{
		LevelL: {
			Total: 404,
			Data:  324,
			Blocks: []blockCapacity{
				{Num: 4, Total: 101, Data: 81},
			},
		},
		LevelM: {
			Total: 404,
			Data:  254,
			Blocks: []blockCapacity{
				{Num: 1, Total: 80, Data: 50},
				{Num: 4, Total: 81, Data: 51},
			},
		},
		LevelQ: {
			Total: 404,
			Data:  180,
			Blocks: []blockCapacity{
				{Num: 4, Total: 50, Data: 22},
				{Num: 4, Total: 51, Data: 23},
			},
		},
		LevelH: {
			Total: 404,
			Data:  140,
			Blocks: []blockCapacity{
				{Num: 3, Total: 36, Data: 12},
				{Num: 8, Total: 37, Data: 13},
			},
		},
	},

	// version 12
	{
		LevelL: {
			Total: 466,
			Data:  370,
			Blocks: []blockCapacity{
				{Num: 2, Total: 116, Data: 92},
				{Num: 2, Total: 117, Data: 93},
			},
		},
		LevelM: {
			Total: 466,
			Data:  290,
			Blocks: []blockCapacity{
				{Num: 6, Total: 58, Data: 36},
				{Num: 2, Total: 59, Data: 37},
			},
		},
		LevelQ: {
			Total: 466,
			Data:  206,
			Blocks: []blockCapacity{
				{Num: 4, Total: 46, Data: 20},
				{Num: 6, Total: 47, Data: 21},
			},
		},
		LevelH: {
			Total: 466,
			Data:  158,
			Blocks: []blockCapacity{
				{Num: 7, Total: 42, Data: 14},
				{Num: 4, Total: 43, Data: 15},
			},
		},
	},

	// version 13
	{
		LevelL: {
			Total: 532,
			Data:  428,
			Blocks: []blockCapacity{
				{Num: 4, Total: 133, Data: 107},
			},
		},
		LevelM: {
			Total: 532,
			Data:  334,
			Blocks: []blockCapacity{
				{Num: 8, Total: 59, Data: 37},
				{Num: 1, Total: 60, Data: 38},
			},
		},
		LevelQ: {
			Total: 532,
			Data:  244,
			Blocks: []blockCapacity{
				{Num: 8, Total: 44, Data: 20},
				{Num: 4, Total: 45, Data: 21},
			},
		},
		LevelH: {
			Total: 532,
			Data:  180,
			Blocks: []blockCapacity{
				{Num: 12, Total: 33, Data: 11},
				{Num: 4, Total: 34, Data: 12},
			},
		},
	},

	// version 14
	{
		LevelL: {
			Total: 581,
			Data:  461,
			Blocks: []blockCapacity{
				{Num: 3, Total: 145, Data: 115},
				{Num: 1, Total: 146, Data: 116},
			},
		},
		LevelM: {
			Total: 581,
			Data:  365,
			Blocks: []blockCapacity{
				{Num: 4, Total: 64, Data: 40},
				{Num: 5, Total: 65, Data: 41},
			},
		},
		LevelQ: {
			Total: 581,
			Data:  261,
			Blocks: []blockCapacity{
				{Num: 11, Total: 36, Data: 16},
				{Num: 5, Total: 37, Data: 17},
			},
		},
		LevelH: {
			Total: 581,
			Data:  197,
			Blocks: []blockCapacity{
				{Num: 11, Total: 36, Data: 12},
				{Num: 5, Total: 37, Data: 13},
			},
		},
	},

	// version 15
	{
		LevelL: {
			Total: 655,
			Data:  523,
			Blocks: []blockCapacity{
				{Num: 5, Total: 109, Data: 87},
				{Num: 1, Total: 110, Data: 88},
			},
		},
		LevelM: {
			Total: 655,
			Data:  415,
			Blocks: []blockCapacity{
				{Num: 5, Total: 65, Data: 41},
				{Num: 5, Total: 66, Data: 42},
			},
		},
		LevelQ: {
			Total: 655,
			Data:  295,
			Blocks: []blockCapacity{
				{Num: 5, Total: 54, Data: 24},
				{Num: 7, Total: 55, Data: 25},
			},
		},
		LevelH: {
			Total: 655,
			Data:  223,
			Blocks: []blockCapacity{
				{Num: 11, Total: 36, Data: 12},
				{Num: 7, Total: 37, Data: 13},
			},
		},
	},

	// version 16
	{
		LevelL: {
			Total: 733,
			Data:  589,
			Blocks: []blockCapacity{
				{Num: 5, Total: 122, Data: 98},
				{Num: 1, Total: 123, Data: 99},
			},
		},
		LevelM: {
			Total: 733,
			Data:  453,
			Blocks: []blockCapacity{
				{Num: 7, Total: 73, Data: 45},
				{Num: 3, Total: 74, Data: 46},
			},
		},
		LevelQ: {
			Total: 733,
			Data:  325,
			Blocks: []blockCapacity{
				{Num: 15, Total: 43, Data: 19},
				{Num: 2, Total: 44, Data: 20},
			},
		},
		LevelH: {
			Total: 733,
			Data:  253,
			Blocks: []blockCapacity{
				{Num: 3, Total: 45, Data: 15},
				{Num: 13, Total: 46, Data: 16},
			},
		},
	},

	// version 17
	{
		LevelL: {
			Total: 815,
			Data:  647,
			Blocks: []blockCapacity{
				{Num: 1, Total: 135, Data: 107},
				{Num: 5, Total: 136, Data: 108},
			},
		},
		LevelM: {
			Total: 815,
			Data:  507,
			Blocks: []blockCapacity{
				{Num: 10, Total: 74, Data: 46},
				{Num: 1, Total: 75, Data: 47},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},
	// version 2
	{
		LevelL: {
			Total: 26,
			Data:  19,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  9,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 40
	{
		LevelL: {
			Total: 3532,
			Data:  2956,
			Blocks: []blockCapacity{
				{Num: 19, Total: 148, Data: 118},
				{Num: 6, Total: 149, Data: 119},
			},
		},
		LevelM: {
			Total: 3706,
			Data:  2334,
			Blocks: []blockCapacity{
				{Num: 18, Total: 75, Data: 47},
				{Num: 31, Total: 76, Data: 48},
			},
		},
		LevelQ: {
			Total: 3706,
			Data:  1666,
			Blocks: []blockCapacity{
				{Num: 34, Total: 54, Data: 24},
				{Num: 34, Total: 55, Data: 25},
			},
		},
		LevelH: {
			Total: 3706,
			Data:  1276,
			Blocks: []blockCapacity{
				{Num: 20, Total: 45, Data: 15},
				{Num: 61, Total: 46, Data: 16},
			},
		},
	},
}
