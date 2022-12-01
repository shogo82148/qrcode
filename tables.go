package qrcode

type capacity struct {
	Total      int // number of total code words
	Data       int // number of data code words
	Correction int // number of correction code words
	Blocks     []blockCapacity
}

type blockCapacity struct {
	Num      int // number of blocks
	Total    int // number of total code words
	Data     int // number of data code words
	MaxError int // maximum number of code word errors
	Reserved int // number of code words reserved for error detection
}

// X 0510 : 2018
// 表9 マイクロORコード及びORコードの誤り訂正特性
var capacityTable = [41][4]capacity{
	{}, // dummy

	// version 1
	{
		LevelL: {
			Total:      26,
			Data:       19,
			Correction: 7,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19, MaxError: 2, Reserved: 3},
			},
		},
		LevelM: {
			Total:      26,
			Data:       16,
			Correction: 10,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16, MaxError: 4, Reserved: 2},
			},
		},
		LevelQ: {
			Total:      26,
			Data:       13,
			Correction: 13,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13, MaxError: 6, Reserved: 1},
			},
		},
		LevelH: {
			Total:      26,
			Data:       9,
			Correction: 17,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9, MaxError: 8, Reserved: 1},
			},
		},
	},

	// version 2
	{
		LevelL: {
			Total:      44,
			Data:       34,
			Correction: 10,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 34, MaxError: 4, Reserved: 2},
			},
		},
		LevelM: {
			Total:      44,
			Data:       28,
			Correction: 16,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 28, MaxError: 8},
			},
		},
		LevelQ: {
			Total:      44,
			Data:       22,
			Correction: 22,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 22, MaxError: 11},
			},
		},
		LevelH: {
			Total:      44,
			Data:       16,
			Correction: 28,
			Blocks: []blockCapacity{
				{Num: 1, Total: 44, Data: 16, MaxError: 14},
			},
		},
	},

	// version 3
	{
		LevelL: {
			Total:      70,
			Data:       55,
			Correction: 15,
			Blocks: []blockCapacity{
				{Num: 1, Total: 70, Data: 55, MaxError: 7},
			},
		},
		LevelM: {
			Total:      70,
			Data:       44,
			Correction: 26,
			Blocks: []blockCapacity{
				{Num: 1, Total: 70, Data: 44, MaxError: 13},
			},
		},
		LevelQ: {
			Total:      70,
			Data:       34,
			Correction: 36,
			Blocks: []blockCapacity{
				{Num: 2, Total: 35, Data: 17, MaxError: 9},
			},
		},
		LevelH: {
			Total:      70,
			Data:       26,
			Correction: 44,
			Blocks: []blockCapacity{
				{Num: 2, Total: 35, Data: 13, MaxError: 11},
			},
		},
	},

	// version 4
	{
		LevelL: {
			Total:      100,
			Data:       80,
			Correction: 20,
			Blocks: []blockCapacity{
				{Num: 1, Total: 100, Data: 80, MaxError: 10},
			},
		},
		LevelM: {
			Total:      100,
			Data:       64,
			Correction: 36,
			Blocks: []blockCapacity{
				{Num: 2, Total: 50, Data: 32, MaxError: 9},
			},
		},
		LevelQ: {
			Total:      100,
			Data:       48,
			Correction: 52,
			Blocks: []blockCapacity{
				{Num: 2, Total: 50, Data: 24, MaxError: 13},
			},
		},
		LevelH: {
			Total:      100,
			Data:       36,
			Correction: 64,
			Blocks: []blockCapacity{
				{Num: 4, Total: 25, Data: 9, MaxError: 8},
			},
		},
	},

	// version 5
	{
		LevelL: {
			Total:      134,
			Data:       108,
			Correction: 26,
			Blocks: []blockCapacity{
				{Num: 1, Total: 134, Data: 108, MaxError: 13},
			},
		},
		LevelM: {
			Total:      134,
			Data:       86,
			Correction: 48,
			Blocks: []blockCapacity{
				{Num: 2, Total: 67, Data: 43, MaxError: 12},
			},
		},
		LevelQ: {
			Total:      134,
			Data:       62,
			Correction: 72,
			Blocks: []blockCapacity{
				{Num: 2, Total: 33, Data: 15, MaxError: 9},
				{Num: 2, Total: 34, Data: 16, MaxError: 9},
			},
		},
		LevelH: {
			Total:      134,
			Data:       46,
			Correction: 88,
			Blocks: []blockCapacity{
				{Num: 2, Total: 33, Data: 11, MaxError: 11},
				{Num: 2, Total: 34, Data: 12, MaxError: 11},
			},
		},
	},

	// version 6
	{
		LevelL: {
			Total:      172,
			Data:       136,
			Correction: 36,
			Blocks: []blockCapacity{
				{Num: 2, Total: 86, Data: 68, MaxError: 9},
			},
		},
		LevelM: {
			Total:      172,
			Data:       108,
			Correction: 64,
			Blocks: []blockCapacity{
				{Num: 4, Total: 43, Data: 27, MaxError: 8},
			},
		},
		LevelQ: {
			Total:      172,
			Data:       76,
			Correction: 96,
			Blocks: []blockCapacity{
				{Num: 4, Total: 43, Data: 19, MaxError: 12},
			},
		},
		LevelH: {
			Total:      172,
			Data:       60,
			Correction: 112,
			Blocks: []blockCapacity{
				{Num: 4, Total: 43, Data: 15, MaxError: 14},
			},
		},
	},

	// version 7
	{
		LevelL: {
			Total:      196,
			Data:       156,
			Correction: 40,
			Blocks: []blockCapacity{
				{Num: 2, Total: 98, Data: 78, MaxError: 10},
			},
		},
		LevelM: {
			Total:      196,
			Data:       124,
			Correction: 72,
			Blocks: []blockCapacity{
				{Num: 4, Total: 49, Data: 31, MaxError: 9},
			},
		},
		LevelQ: {
			Total:      196,
			Data:       88,
			Correction: 108,
			Blocks: []blockCapacity{
				{Num: 2, Total: 32, Data: 14, MaxError: 9},
				{Num: 4, Total: 33, Data: 15, MaxError: 9},
			},
		},
		LevelH: {
			Total:      196,
			Data:       66,
			Correction: 130,
			Blocks: []blockCapacity{
				{Num: 4, Total: 39, Data: 13, MaxError: 13},
				{Num: 1, Total: 40, Data: 14, MaxError: 13},
			},
		},
	},

	// version 8
	{
		LevelL: {
			Total:      242,
			Data:       194,
			Correction: 48,
			Blocks: []blockCapacity{
				{Num: 2, Total: 121, Data: 97, MaxError: 12},
			},
		},
		LevelM: {
			Total:      242,
			Data:       154,
			Correction: 88,
			Blocks: []blockCapacity{
				{Num: 2, Total: 60, Data: 38, MaxError: 11},
				{Num: 2, Total: 61, Data: 39, MaxError: 11},
			},
		},
		LevelQ: {
			Total:      242,
			Data:       110,
			Correction: 132,
			Blocks: []blockCapacity{
				{Num: 4, Total: 40, Data: 18, MaxError: 11},
				{Num: 2, Total: 41, Data: 19, MaxError: 11},
			},
		},
		LevelH: {
			Total:      242,
			Data:       86,
			Correction: 156,
			Blocks: []blockCapacity{
				{Num: 4, Total: 40, Data: 14, MaxError: 13},
				{Num: 2, Total: 41, Data: 15, MaxError: 13},
			},
		},
	},

	// version 9
	{
		LevelL: {
			Total:      292,
			Data:       232,
			Correction: 60,
			Blocks: []blockCapacity{
				{Num: 2, Total: 146, Data: 116, MaxError: 15},
			},
		},
		LevelM: {
			Total:      292,
			Data:       182,
			Correction: 110,
			Blocks: []blockCapacity{
				{Num: 3, Total: 58, Data: 36, MaxError: 11},
				{Num: 2, Total: 59, Data: 37, MaxError: 11},
			},
		},
		LevelQ: {
			Total:      292,
			Data:       132,
			Correction: 160,
			Blocks: []blockCapacity{
				{Num: 4, Total: 36, Data: 16, MaxError: 10},
				{Num: 4, Total: 37, Data: 17, MaxError: 10},
			},
		},
		LevelH: {
			Total:      292,
			Data:       100,
			Correction: 192,
			Blocks: []blockCapacity{
				{Num: 4, Total: 36, Data: 12, MaxError: 12},
				{Num: 4, Total: 37, Data: 13, MaxError: 12},
			},
		},
	},

	// version 10
	{
		LevelL: {
			Total:      346,
			Data:       274,
			Correction: 72,
			Blocks: []blockCapacity{
				{Num: 2, Total: 86, Data: 68, MaxError: 9},
				{Num: 2, Total: 87, Data: 69, MaxError: 9},
			},
		},
		LevelM: {
			Total:      346,
			Data:       216,
			Correction: 130,
			Blocks: []blockCapacity{
				{Num: 4, Total: 69, Data: 43, MaxError: 13},
				{Num: 1, Total: 70, Data: 44, MaxError: 13},
			},
		},
		LevelQ: {
			Total:      346,
			Data:       154,
			Correction: 192,
			Blocks: []blockCapacity{
				{Num: 6, Total: 43, Data: 19, MaxError: 12},
				{Num: 2, Total: 44, Data: 20, MaxError: 12},
			},
		},
		LevelH: {
			Total:      346,
			Data:       122,
			Correction: 224,
			Blocks: []blockCapacity{
				{Num: 6, Total: 43, Data: 15, MaxError: 14},
				{Num: 2, Total: 44, Data: 16, MaxError: 14},
			},
		},
	},

	// version 11
	{
		LevelL: {
			Total:      404,
			Data:       324,
			Correction: 80,
			Blocks: []blockCapacity{
				{Num: 4, Total: 101, Data: 81, MaxError: 10},
			},
		},
		LevelM: {
			Total:      404,
			Data:       254,
			Correction: 150,
			Blocks: []blockCapacity{
				{Num: 1, Total: 80, Data: 50, MaxError: 15},
				{Num: 4, Total: 81, Data: 51, MaxError: 15},
			},
		},
		LevelQ: {
			Total:      404,
			Data:       180,
			Correction: 224,
			Blocks: []blockCapacity{
				{Num: 4, Total: 50, Data: 22, MaxError: 14},
				{Num: 4, Total: 51, Data: 23, MaxError: 14},
			},
		},
		LevelH: {
			Total:      404,
			Data:       140,
			Correction: 264,
			Blocks: []blockCapacity{
				{Num: 3, Total: 36, Data: 12, MaxError: 12},
				{Num: 8, Total: 37, Data: 13, MaxError: 12},
			},
		},
	},

	// version 12
	{
		LevelL: {
			Total:      466,
			Data:       370,
			Correction: 96,
			Blocks: []blockCapacity{
				{Num: 2, Total: 116, Data: 92, MaxError: 12},
				{Num: 2, Total: 117, Data: 93, MaxError: 12},
			},
		},
		LevelM: {
			Total:      466,
			Data:       290,
			Correction: 176,
			Blocks: []blockCapacity{
				{Num: 6, Total: 58, Data: 36, MaxError: 11},
				{Num: 2, Total: 59, Data: 37, MaxError: 11},
			},
		},
		LevelQ: {
			Total:      466,
			Data:       206,
			Correction: 260,
			Blocks: []blockCapacity{
				{Num: 4, Total: 46, Data: 20, MaxError: 13},
				{Num: 6, Total: 47, Data: 21, MaxError: 13},
			},
		},
		LevelH: {
			Total:      466,
			Data:       158,
			Correction: 308,
			Blocks: []blockCapacity{
				{Num: 7, Total: 42, Data: 14, MaxError: 14},
				{Num: 4, Total: 43, Data: 15, MaxError: 14},
			},
		},
	},

	// version 13
	{
		LevelL: {
			Total:      532,
			Data:       428,
			Correction: 104,
			Blocks: []blockCapacity{
				{Num: 4, Total: 133, Data: 107, MaxError: 13},
			},
		},
		LevelM: {
			Total:      532,
			Data:       334,
			Correction: 198,
			Blocks: []blockCapacity{
				{Num: 8, Total: 59, Data: 37, MaxError: 11},
				{Num: 1, Total: 60, Data: 38, MaxError: 11},
			},
		},
		LevelQ: {
			Total:      532,
			Data:       244,
			Correction: 288,
			Blocks: []blockCapacity{
				{Num: 8, Total: 44, Data: 20, MaxError: 12},
				{Num: 4, Total: 45, Data: 21, MaxError: 12},
			},
		},
		LevelH: {
			Total:      532,
			Data:       180,
			Correction: 352,
			Blocks: []blockCapacity{
				{Num: 12, Total: 33, Data: 11, MaxError: 11},
				{Num: 4, Total: 34, Data: 12, MaxError: 11},
			},
		},
	},

	// version 14
	{
		LevelL: {
			Total:      581,
			Data:       461,
			Correction: 120,
			Blocks: []blockCapacity{
				{Num: 3, Total: 145, Data: 115, MaxError: 15},
				{Num: 1, Total: 146, Data: 116, MaxError: 15},
			},
		},
		LevelM: {
			Total:      581,
			Data:       365,
			Correction: 216,
			Blocks: []blockCapacity{
				{Num: 4, Total: 64, Data: 40, MaxError: 12},
				{Num: 5, Total: 65, Data: 41, MaxError: 12},
			},
		},
		LevelQ: {
			Total:      581,
			Data:       261,
			Correction: 320,
			Blocks: []blockCapacity{
				{Num: 11, Total: 36, Data: 16, MaxError: 10},
				{Num: 5, Total: 37, Data: 17, MaxError: 10},
			},
		},
		LevelH: {
			Total:      581,
			Data:       197,
			Correction: 384,
			Blocks: []blockCapacity{
				{Num: 11, Total: 36, Data: 12, MaxError: 12},
				{Num: 5, Total: 37, Data: 13, MaxError: 12},
			},
		},
	},

	// version 15
	{
		LevelL: {
			Total:      655,
			Data:       523,
			Correction: 132,
			Blocks: []blockCapacity{
				{Num: 5, Total: 109, Data: 87, MaxError: 11},
				{Num: 1, Total: 110, Data: 88, MaxError: 11},
			},
		},
		LevelM: {
			Total:      655,
			Data:       415,
			Correction: 240,
			Blocks: []blockCapacity{
				{Num: 5, Total: 65, Data: 41, MaxError: 12},
				{Num: 5, Total: 66, Data: 42, MaxError: 12},
			},
		},
		LevelQ: {
			Total:      655,
			Data:       295,
			Correction: 360,
			Blocks: []blockCapacity{
				{Num: 5, Total: 54, Data: 24, MaxError: 15},
				{Num: 7, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      655,
			Data:       223,
			Correction: 432,
			Blocks: []blockCapacity{
				{Num: 11, Total: 36, Data: 12, MaxError: 12},
				{Num: 7, Total: 37, Data: 13, MaxError: 12},
			},
		},
	},

	// version 16
	{
		LevelL: {
			Total:      733,
			Data:       589,
			Correction: 144,
			Blocks: []blockCapacity{
				{Num: 5, Total: 122, Data: 98, MaxError: 12},
				{Num: 1, Total: 123, Data: 99, MaxError: 12},
			},
		},
		LevelM: {
			Total:      733,
			Data:       453,
			Correction: 280,
			Blocks: []blockCapacity{
				{Num: 7, Total: 73, Data: 45, MaxError: 14},
				{Num: 3, Total: 74, Data: 46, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      733,
			Data:       325,
			Correction: 408,
			Blocks: []blockCapacity{
				{Num: 15, Total: 43, Data: 19, MaxError: 12},
				{Num: 2, Total: 44, Data: 20, MaxError: 12},
			},
		},
		LevelH: {
			Total:      733,
			Data:       253,
			Correction: 480,
			Blocks: []blockCapacity{
				{Num: 3, Total: 45, Data: 15, MaxError: 15},
				{Num: 13, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 17
	{
		LevelL: {
			Total:      815,
			Data:       647,
			Correction: 168,
			Blocks: []blockCapacity{
				{Num: 1, Total: 135, Data: 107, MaxError: 14},
				{Num: 5, Total: 136, Data: 108, MaxError: 14},
			},
		},
		LevelM: {
			Total:      815,
			Data:       507,
			Correction: 308,
			Blocks: []blockCapacity{
				{Num: 10, Total: 74, Data: 46, MaxError: 14},
				{Num: 1, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      815,
			Data:       367,
			Correction: 448,
			Blocks: []blockCapacity{
				{Num: 1, Total: 50, Data: 22, MaxError: 14},
				{Num: 15, Total: 51, Data: 23, MaxError: 14},
			},
		},
		LevelH: {
			Total:      815,
			Data:       283,
			Correction: 532,
			Blocks: []blockCapacity{
				{Num: 2, Total: 42, Data: 14, MaxError: 14},
				{Num: 17, Total: 43, Data: 15, MaxError: 14},
			},
		},
	},

	// version 18
	{
		LevelL: {
			Total:      901,
			Data:       721,
			Correction: 180,
			Blocks: []blockCapacity{
				{Num: 5, Total: 150, Data: 120, MaxError: 15},
				{Num: 1, Total: 151, Data: 121, MaxError: 15},
			},
		},
		LevelM: {
			Total:      901,
			Data:       563,
			Correction: 338,
			Blocks: []blockCapacity{
				{Num: 9, Total: 69, Data: 43, MaxError: 13},
				{Num: 4, Total: 70, Data: 44, MaxError: 13},
			},
		},
		LevelQ: {
			Total:      901,
			Data:       397,
			Correction: 504,
			Blocks: []blockCapacity{
				{Num: 17, Total: 50, Data: 22, MaxError: 14},
				{Num: 1, Total: 51, Data: 23, MaxError: 14},
			},
		},
		LevelH: {
			Total:      901,
			Data:       313,
			Correction: 588,
			Blocks: []blockCapacity{
				{Num: 2, Total: 42, Data: 14, MaxError: 14},
				{Num: 19, Total: 43, Data: 15, MaxError: 14},
			},
		},
	},

	// version 19
	{
		LevelL: {
			Total:      991,
			Data:       795,
			Correction: 196,
			Blocks: []blockCapacity{
				{Num: 3, Total: 141, Data: 113, MaxError: 14},
				{Num: 4, Total: 142, Data: 114, MaxError: 14},
			},
		},
		LevelM: {
			Total:      991,
			Data:       627,
			Correction: 364,
			Blocks: []blockCapacity{
				{Num: 3, Total: 70, Data: 44, MaxError: 13},
				{Num: 11, Total: 71, Data: 45, MaxError: 13},
			},
		},
		LevelQ: {
			Total:      991,
			Data:       445,
			Correction: 546,
			Blocks: []blockCapacity{
				{Num: 17, Total: 47, Data: 21, MaxError: 13},
				{Num: 4, Total: 48, Data: 22, MaxError: 13},
			},
		},
		LevelH: {
			Total:      991,
			Data:       341,
			Correction: 650,
			Blocks: []blockCapacity{
				{Num: 9, Total: 39, Data: 13, MaxError: 13},
				{Num: 16, Total: 40, Data: 14, MaxError: 13},
			},
		},
	},

	// version 20
	{
		LevelL: {
			Total:      1085,
			Data:       861,
			Correction: 224,
			Blocks: []blockCapacity{
				{Num: 3, Total: 135, Data: 107, MaxError: 14},
				{Num: 5, Total: 136, Data: 108, MaxError: 14},
			},
		},
		LevelM: {
			Total:      1085,
			Data:       669,
			Correction: 416,
			Blocks: []blockCapacity{
				{Num: 3, Total: 67, Data: 41, MaxError: 13},
				{Num: 13, Total: 68, Data: 42, MaxError: 13},
			},
		},
		LevelQ: {
			Total:      1085,
			Data:       485,
			Correction: 600,
			Blocks: []blockCapacity{
				{Num: 15, Total: 54, Data: 24, MaxError: 15},
				{Num: 5, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      1085,
			Data:       385,
			Correction: 700,
			Blocks: []blockCapacity{
				{Num: 15, Total: 43, Data: 15, MaxError: 14},
				{Num: 10, Total: 44, Data: 16, MaxError: 14},
			},
		},
	},

	// version 21
	{
		LevelL: {
			Total:      1156,
			Data:       932,
			Correction: 224,
			Blocks: []blockCapacity{
				{Num: 4, Total: 144, Data: 116, MaxError: 14},
				{Num: 4, Total: 145, Data: 117, MaxError: 14},
			},
		},
		LevelM: {
			Total:      1156,
			Data:       714,
			Correction: 442,
			Blocks: []blockCapacity{
				{Num: 17, Total: 68, Data: 42, MaxError: 13},
			},
		},
		LevelQ: {
			Total:      1156,
			Data:       512,
			Correction: 644,
			Blocks: []blockCapacity{
				{Num: 17, Total: 50, Data: 22, MaxError: 14},
				{Num: 6, Total: 51, Data: 23, MaxError: 14},
			},
		},
		LevelH: {
			Total:      1156,
			Data:       406,
			Correction: 750,
			Blocks: []blockCapacity{
				{Num: 19, Total: 46, Data: 16, MaxError: 15},
				{Num: 6, Total: 47, Data: 17, MaxError: 15},
			},
		},
	},

	// version 22
	{
		LevelL: {
			Total:      1258,
			Data:       1006,
			Correction: 252,
			Blocks: []blockCapacity{
				{Num: 2, Total: 139, Data: 111, MaxError: 14},
				{Num: 7, Total: 140, Data: 112, MaxError: 14},
			},
		},
		LevelM: {
			Total:      1258,
			Data:       782,
			Correction: 476,
			Blocks: []blockCapacity{
				{Num: 17, Total: 74, Data: 46, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      1258,
			Data:       568,
			Correction: 690,
			Blocks: []blockCapacity{
				{Num: 7, Total: 54, Data: 24, MaxError: 15},
				{Num: 16, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      1258,
			Data:       442,
			Correction: 816,
			Blocks: []blockCapacity{
				{Num: 34, Total: 37, Data: 13, MaxError: 12},
			},
		},
	},

	// version 23
	{
		LevelL: {
			Total:      1364,
			Data:       1094,
			Correction: 270,
			Blocks: []blockCapacity{
				{Num: 4, Total: 151, Data: 121, MaxError: 15},
				{Num: 5, Total: 152, Data: 122, MaxError: 15},
			},
		},
		LevelM: {
			Total:      1364,
			Data:       860,
			Correction: 504,
			Blocks: []blockCapacity{
				{Num: 4, Total: 75, Data: 47, MaxError: 14},
				{Num: 14, Total: 76, Data: 48, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      1364,
			Data:       614,
			Correction: 750,
			Blocks: []blockCapacity{
				{Num: 11, Total: 54, Data: 24, MaxError: 15},
				{Num: 14, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      1364,
			Data:       464,
			Correction: 900,
			Blocks: []blockCapacity{
				{Num: 16, Total: 45, Data: 15, MaxError: 15},
				{Num: 14, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 24
	{
		LevelL: {
			Total:      1474,
			Data:       1174,
			Correction: 300,
			Blocks: []blockCapacity{
				{Num: 6, Total: 147, Data: 117, MaxError: 15},
				{Num: 4, Total: 148, Data: 118, MaxError: 15},
			},
		},
		LevelM: {
			Total:      1474,
			Data:       914,
			Correction: 560,
			Blocks: []blockCapacity{
				{Num: 6, Total: 73, Data: 45, MaxError: 14},
				{Num: 14, Total: 74, Data: 46, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      1474,
			Data:       664,
			Correction: 810,
			Blocks: []blockCapacity{
				{Num: 11, Total: 54, Data: 24, MaxError: 15},
				{Num: 16, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      1474,
			Data:       514,
			Correction: 960,
			Blocks: []blockCapacity{
				{Num: 30, Total: 46, Data: 16, MaxError: 15},
				{Num: 2, Total: 47, Data: 17, MaxError: 15},
			},
		},
	},

	// version 25
	{
		LevelL: {
			Total:      1588,
			Data:       1276,
			Correction: 312,
			Blocks: []blockCapacity{
				{Num: 8, Total: 132, Data: 106, MaxError: 13},
				{Num: 4, Total: 133, Data: 107, MaxError: 13},
			},
		},
		LevelM: {
			Total:      1588,
			Data:       1000,
			Correction: 588,
			Blocks: []blockCapacity{
				{Num: 8, Total: 75, Data: 47, MaxError: 14},
				{Num: 13, Total: 76, Data: 48, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      1588,
			Data:       718,
			Correction: 870,
			Blocks: []blockCapacity{
				{Num: 7, Total: 54, Data: 24, MaxError: 15},
				{Num: 22, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      1588,
			Data:       538,
			Correction: 1050,
			Blocks: []blockCapacity{
				{Num: 22, Total: 45, Data: 15, MaxError: 15},
				{Num: 13, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 26
	{
		LevelL: {
			Total:      1706,
			Data:       1370,
			Correction: 336,
			Blocks: []blockCapacity{
				{Num: 10, Total: 142, Data: 114, MaxError: 14},
				{Num: 2, Total: 143, Data: 115, MaxError: 14},
			},
		},
		LevelM: {
			Total:      1706,
			Data:       1062,
			Correction: 644,
			Blocks: []blockCapacity{
				{Num: 19, Total: 74, Data: 46, MaxError: 14},
				{Num: 4, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      1706,
			Data:       754,
			Correction: 952,
			Blocks: []blockCapacity{
				{Num: 28, Total: 50, Data: 22, MaxError: 14},
				{Num: 6, Total: 51, Data: 23, MaxError: 14},
			},
		},
		LevelH: {
			Total:      1706,
			Data:       596,
			Correction: 1110,
			Blocks: []blockCapacity{
				{Num: 33, Total: 46, Data: 16, MaxError: 15},
				{Num: 4, Total: 47, Data: 17, MaxError: 15},
			},
		},
	},

	// version 27
	{
		LevelL: {
			Total:      1828,
			Data:       1468,
			Correction: 360,
			Blocks: []blockCapacity{
				{Num: 8, Total: 152, Data: 122, MaxError: 15},
				{Num: 4, Total: 153, Data: 123, MaxError: 15},
			},
		},
		LevelM: {
			Total:      1828,
			Data:       1128,
			Correction: 700,
			Blocks: []blockCapacity{
				{Num: 22, Total: 73, Data: 45, MaxError: 14},
				{Num: 3, Total: 74, Data: 46, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      1828,
			Data:       808,
			Correction: 1020,
			Blocks: []blockCapacity{
				{Num: 8, Total: 53, Data: 23, MaxError: 15},
				{Num: 26, Total: 54, Data: 24, MaxError: 15},
			},
		},
		LevelH: {
			Total:      1828,
			Data:       628,
			Correction: 1200,
			Blocks: []blockCapacity{
				{Num: 12, Total: 45, Data: 15, MaxError: 15},
				{Num: 28, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 28
	{
		LevelL: {
			Total:      1921,
			Data:       1531,
			Correction: 390,
			Blocks: []blockCapacity{
				{Num: 3, Total: 147, Data: 117, MaxError: 15},
				{Num: 10, Total: 148, Data: 118, MaxError: 15},
			},
		},
		LevelM: {
			Total:      1921,
			Data:       1193,
			Correction: 728,
			Blocks: []blockCapacity{
				{Num: 3, Total: 73, Data: 45, MaxError: 14},
				{Num: 23, Total: 74, Data: 46, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      1921,
			Data:       871,
			Correction: 1050,
			Blocks: []blockCapacity{
				{Num: 4, Total: 54, Data: 24, MaxError: 15},
				{Num: 31, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      1921,
			Data:       661,
			Correction: 1260,
			Blocks: []blockCapacity{
				{Num: 11, Total: 45, Data: 15, MaxError: 15},
				{Num: 31, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 29
	{
		LevelL: {
			Total:      2051,
			Data:       1631,
			Correction: 420,
			Blocks: []blockCapacity{
				{Num: 7, Total: 146, Data: 116, MaxError: 15},
				{Num: 7, Total: 147, Data: 117, MaxError: 15},
			},
		},
		LevelM: {
			Total:      2051,
			Data:       1267,
			Correction: 784,
			Blocks: []blockCapacity{
				{Num: 21, Total: 73, Data: 45, MaxError: 14},
				{Num: 7, Total: 74, Data: 46, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      2051,
			Data:       911,
			Correction: 1140,
			Blocks: []blockCapacity{
				{Num: 1, Total: 53, Data: 23, MaxError: 15},
				{Num: 37, Total: 54, Data: 24, MaxError: 15},
			},
		},
		LevelH: {
			Total:      2051,
			Data:       701,
			Correction: 1350,
			Blocks: []blockCapacity{
				{Num: 19, Total: 45, Data: 15, MaxError: 15},
				{Num: 26, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 30
	{
		LevelL: {
			Total:      2185,
			Data:       1735,
			Correction: 450,
			Blocks: []blockCapacity{
				{Num: 5, Total: 145, Data: 115, MaxError: 15},
				{Num: 10, Total: 146, Data: 116, MaxError: 15},
			},
		},
		LevelM: {
			Total:      2185,
			Data:       1373,
			Correction: 812,
			Blocks: []blockCapacity{
				{Num: 19, Total: 75, Data: 47, MaxError: 14},
				{Num: 10, Total: 76, Data: 48, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      2185,
			Data:       985,
			Correction: 1200,
			Blocks: []blockCapacity{
				{Num: 15, Total: 54, Data: 24, MaxError: 15},
				{Num: 25, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      2185,
			Data:       745,
			Correction: 1440,
			Blocks: []blockCapacity{
				{Num: 23, Total: 45, Data: 15, MaxError: 15},
				{Num: 25, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 31
	{
		LevelL: {
			Total:      2323,
			Data:       1843,
			Correction: 480,
			Blocks: []blockCapacity{
				{Num: 13, Total: 145, Data: 115, MaxError: 15},
				{Num: 3, Total: 146, Data: 116, MaxError: 15},
			},
		},
		LevelM: {
			Total:      2323,
			Data:       1455,
			Correction: 868,
			Blocks: []blockCapacity{
				{Num: 2, Total: 74, Data: 46, MaxError: 14},
				{Num: 29, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      2323,
			Data:       1033,
			Correction: 1290,
			Blocks: []blockCapacity{
				{Num: 42, Total: 54, Data: 24, MaxError: 15},
				{Num: 1, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      2323,
			Data:       793,
			Correction: 1530,
			Blocks: []blockCapacity{
				{Num: 23, Total: 45, Data: 15, MaxError: 15},
				{Num: 28, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 32
	{
		LevelL: {
			Total:      2465,
			Data:       1955,
			Correction: 510,
			Blocks: []blockCapacity{
				{Num: 17, Total: 145, Data: 115, MaxError: 15},
			},
		},
		LevelM: {
			Total:      2465,
			Data:       1541,
			Correction: 924,
			Blocks: []blockCapacity{
				{Num: 10, Total: 74, Data: 46, MaxError: 14},
				{Num: 23, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      2465,
			Data:       1115,
			Correction: 1350,
			Blocks: []blockCapacity{
				{Num: 10, Total: 54, Data: 24, MaxError: 15},
				{Num: 35, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      2465,
			Data:       845,
			Correction: 1620,
			Blocks: []blockCapacity{
				{Num: 19, Total: 45, Data: 15, MaxError: 15},
				{Num: 35, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 33
	{
		LevelL: {
			Total:      2611,
			Data:       2071,
			Correction: 540,
			Blocks: []blockCapacity{
				{Num: 17, Total: 145, Data: 115, MaxError: 15},
				{Num: 1, Total: 146, Data: 116, MaxError: 15},
			},
		},
		LevelM: {
			Total:      2611,
			Data:       1631,
			Correction: 980,
			Blocks: []blockCapacity{
				{Num: 14, Total: 74, Data: 46, MaxError: 14},
				{Num: 21, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      2611,
			Data:       1171,
			Correction: 1440,
			Blocks: []blockCapacity{
				{Num: 29, Total: 54, Data: 24, MaxError: 15},
				{Num: 19, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      2611,
			Data:       901,
			Correction: 1710,
			Blocks: []blockCapacity{
				{Num: 11, Total: 45, Data: 15, MaxError: 15},
				{Num: 46, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 34
	{
		LevelL: {
			Total:      2761,
			Data:       2191,
			Correction: 570,
			Blocks: []blockCapacity{
				{Num: 13, Total: 145, Data: 115, MaxError: 15},
				{Num: 6, Total: 146, Data: 116, MaxError: 15},
			},
		},
		LevelM: {
			Total:      2761,
			Data:       1725,
			Correction: 1036,
			Blocks: []blockCapacity{
				{Num: 14, Total: 74, Data: 46, MaxError: 14},
				{Num: 23, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      2761,
			Data:       1231,
			Correction: 1530,
			Blocks: []blockCapacity{
				{Num: 44, Total: 54, Data: 24, MaxError: 15},
				{Num: 7, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      2761,
			Data:       961,
			Correction: 1800,
			Blocks: []blockCapacity{
				{Num: 59, Total: 46, Data: 16, MaxError: 15},
				{Num: 1, Total: 47, Data: 17, MaxError: 15},
			},
		},
	},

	// version 35
	{
		LevelL: {
			Total:      2876,
			Data:       2306,
			Correction: 570,
			Blocks: []blockCapacity{
				{Num: 12, Total: 151, Data: 121, MaxError: 15},
				{Num: 7, Total: 152, Data: 122, MaxError: 15},
			},
		},
		LevelM: {
			Total:      2876,
			Data:       1812,
			Correction: 1064,
			Blocks: []blockCapacity{
				{Num: 12, Total: 75, Data: 47, MaxError: 14},
				{Num: 26, Total: 76, Data: 48, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      2876,
			Data:       1286,
			Correction: 1590,
			Blocks: []blockCapacity{
				{Num: 39, Total: 54, Data: 24, MaxError: 15},
				{Num: 14, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      2876,
			Data:       986,
			Correction: 1890,
			Blocks: []blockCapacity{
				{Num: 22, Total: 45, Data: 15, MaxError: 15},
				{Num: 41, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 36
	{
		LevelL: {
			Total:      3034,
			Data:       2434,
			Correction: 600,
			Blocks: []blockCapacity{
				{Num: 6, Total: 151, Data: 121, MaxError: 15},
				{Num: 14, Total: 152, Data: 122, MaxError: 15},
			},
		},
		LevelM: {
			Total:      3034,
			Data:       1914,
			Correction: 1120,
			Blocks: []blockCapacity{
				{Num: 6, Total: 75, Data: 47, MaxError: 14},
				{Num: 34, Total: 76, Data: 48, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      3034,
			Data:       1354,
			Correction: 1680,
			Blocks: []blockCapacity{
				{Num: 46, Total: 54, Data: 24, MaxError: 15},
				{Num: 10, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      3034,
			Data:       1054,
			Correction: 1980,
			Blocks: []blockCapacity{
				{Num: 2, Total: 45, Data: 15, MaxError: 15},
				{Num: 64, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 37
	{
		LevelL: {
			Total:      3196,
			Data:       2566,
			Correction: 630,
			Blocks: []blockCapacity{
				{Num: 17, Total: 152, Data: 122, MaxError: 15},
				{Num: 4, Total: 153, Data: 123, MaxError: 15},
			},
		},
		LevelM: {
			Total:      3196,
			Data:       1992,
			Correction: 1204,
			Blocks: []blockCapacity{
				{Num: 29, Total: 74, Data: 46, MaxError: 14},
				{Num: 14, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      3196,
			Data:       1426,
			Correction: 1770,
			Blocks: []blockCapacity{
				{Num: 49, Total: 54, Data: 24, MaxError: 15},
				{Num: 10, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      3196,
			Data:       1096,
			Correction: 2100,
			Blocks: []blockCapacity{
				{Num: 24, Total: 45, Data: 15, MaxError: 15},
				{Num: 46, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 38
	{
		LevelL: {
			Total:      3362,
			Data:       2702,
			Correction: 660,
			Blocks: []blockCapacity{
				{Num: 4, Total: 152, Data: 122, MaxError: 15},
				{Num: 18, Total: 153, Data: 123, MaxError: 15},
			},
		},
		LevelM: {
			Total:      3362,
			Data:       2102,
			Correction: 1260,
			Blocks: []blockCapacity{
				{Num: 13, Total: 74, Data: 46, MaxError: 14},
				{Num: 32, Total: 75, Data: 47, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      3362,
			Data:       1502,
			Correction: 1860,
			Blocks: []blockCapacity{
				{Num: 48, Total: 54, Data: 24, MaxError: 15},
				{Num: 14, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      3362,
			Data:       1142,
			Correction: 2220,
			Blocks: []blockCapacity{
				{Num: 42, Total: 45, Data: 15, MaxError: 15},
				{Num: 32, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 39
	{
		LevelL: {
			Total:      3532,
			Data:       2812,
			Correction: 720,
			Blocks: []blockCapacity{
				{Num: 20, Total: 147, Data: 117, MaxError: 15},
				{Num: 4, Total: 148, Data: 118, MaxError: 15},
			},
		},
		LevelM: {
			Total:      3532,
			Data:       2216,
			Correction: 1316,
			Blocks: []blockCapacity{
				{Num: 40, Total: 75, Data: 47, MaxError: 14},
				{Num: 7, Total: 76, Data: 48, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      3532,
			Data:       1582,
			Correction: 1950,
			Blocks: []blockCapacity{
				{Num: 43, Total: 54, Data: 24, MaxError: 15},
				{Num: 22, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      3532,
			Data:       1222,
			Correction: 2310,
			Blocks: []blockCapacity{
				{Num: 10, Total: 45, Data: 15, MaxError: 15},
				{Num: 67, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},

	// version 40
	{
		LevelL: {
			Total:      3706,
			Data:       2956,
			Correction: 750,
			Blocks: []blockCapacity{
				{Num: 19, Total: 148, Data: 118, MaxError: 15},
				{Num: 6, Total: 149, Data: 119, MaxError: 15},
			},
		},
		LevelM: {
			Total:      3706,
			Data:       2334,
			Correction: 1372,
			Blocks: []blockCapacity{
				{Num: 18, Total: 75, Data: 47, MaxError: 14},
				{Num: 31, Total: 76, Data: 48, MaxError: 14},
			},
		},
		LevelQ: {
			Total:      3706,
			Data:       1666,
			Correction: 2040,
			Blocks: []blockCapacity{
				{Num: 34, Total: 54, Data: 24, MaxError: 15},
				{Num: 34, Total: 55, Data: 25, MaxError: 15},
			},
		},
		LevelH: {
			Total:      3706,
			Data:       1276,
			Correction: 2430,
			Blocks: []blockCapacity{
				{Num: 20, Total: 45, Data: 15, MaxError: 15},
				{Num: 61, Total: 46, Data: 16, MaxError: 15},
			},
		},
	},
}
