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
			Total: 815,
			Data:  367,
			Blocks: []blockCapacity{
				{Num: 1, Total: 50, Data: 22},
				{Num: 15, Total: 51, Data: 23},
			},
		},
		LevelH: {
			Total: 815,
			Data:  283,
			Blocks: []blockCapacity{
				{Num: 2, Total: 42, Data: 14},
				{Num: 17, Total: 43, Data: 15},
			},
		},
	},

	// version 18
	{
		LevelL: {
			Total: 901,
			Data:  721,
			Blocks: []blockCapacity{
				{Num: 5, Total: 150, Data: 120},
				{Num: 1, Total: 151, Data: 121},
			},
		},
		LevelM: {
			Total: 901,
			Data:  563,
			Blocks: []blockCapacity{
				{Num: 9, Total: 69, Data: 43},
				{Num: 4, Total: 70, Data: 44},
			},
		},
		LevelQ: {
			Total: 901,
			Data:  397,
			Blocks: []blockCapacity{
				{Num: 17, Total: 50, Data: 22},
				{Num: 1, Total: 51, Data: 23},
			},
		},
		LevelH: {
			Total: 901,
			Data:  313,
			Blocks: []blockCapacity{
				{Num: 2, Total: 42, Data: 14},
				{Num: 19, Total: 43, Data: 15},
			},
		},
	},

	// version 19
	{
		LevelL: {
			Total: 991,
			Data:  795,
			Blocks: []blockCapacity{
				{Num: 3, Total: 141, Data: 113},
				{Num: 4, Total: 142, Data: 114},
			},
		},
		LevelM: {
			Total: 991,
			Data:  627,
			Blocks: []blockCapacity{
				{Num: 3, Total: 70, Data: 44},
				{Num: 11, Total: 71, Data: 45},
			},
		},
		LevelQ: {
			Total: 991,
			Data:  445,
			Blocks: []blockCapacity{
				{Num: 17, Total: 47, Data: 21},
				{Num: 4, Total: 48, Data: 22},
			},
		},
		LevelH: {
			Total: 991,
			Data:  341,
			Blocks: []blockCapacity{
				{Num: 9, Total: 39, Data: 13},
				{Num: 16, Total: 40, Data: 14},
			},
		},
	},

	// version 20
	{
		LevelL: {
			Total: 1085,
			Data:  861,
			Blocks: []blockCapacity{
				{Num: 3, Total: 135, Data: 107},
				{Num: 5, Total: 136, Data: 108},
			},
		},
		LevelM: {
			Total: 1085,
			Data:  669,
			Blocks: []blockCapacity{
				{Num: 3, Total: 67, Data: 41},
				{Num: 13, Total: 68, Data: 42},
			},
		},
		LevelQ: {
			Total: 1085,
			Data:  485,
			Blocks: []blockCapacity{
				{Num: 15, Total: 54, Data: 24},
				{Num: 5, Total: 55, Data: 25},
			},
		},
		LevelH: {
			Total: 1085,
			Data:  385,
			Blocks: []blockCapacity{
				{Num: 15, Total: 43, Data: 15},
				{Num: 10, Total: 44, Data: 16},
			},
		},
	},

	// version 21
	{
		LevelL: {
			Total: 1156,
			Data:  932,
			Blocks: []blockCapacity{
				{Num: 4, Total: 144, Data: 116},
				{Num: 4, Total: 145, Data: 117},
			},
		},
		LevelM: {
			Total: 1156,
			Data:  714,
			Blocks: []blockCapacity{
				{Num: 17, Total: 68, Data: 42},
			},
		},
		LevelQ: {
			Total: 1156,
			Data:  512,
			Blocks: []blockCapacity{
				{Num: 17, Total: 50, Data: 22},
				{Num: 6, Total: 51, Data: 23},
			},
		},
		LevelH: {
			Total: 1156,
			Data:  406,
			Blocks: []blockCapacity{
				{Num: 19, Total: 46, Data: 16},
				{Num: 6, Total: 47, Data: 17},
			},
		},
	},

	// version 22
	{
		LevelL: {
			Total: 1258,
			Data:  1006,
			Blocks: []blockCapacity{
				{Num: 2, Total: 139, Data: 111},
				{Num: 7, Total: 140, Data: 112},
			},
		},
		LevelM: {
			Total: 1258,
			Data:  782,
			Blocks: []blockCapacity{
				{Num: 17, Total: 74, Data: 46},
			},
		},
		LevelQ: {
			Total: 1258,
			Data:  568,
			Blocks: []blockCapacity{
				{Num: 7, Total: 54, Data: 24},
				{Num: 16, Total: 55, Data: 25},
			},
		},
		LevelH: {
			Total: 1258,
			Data:  442,
			Blocks: []blockCapacity{
				{Num: 34, Total: 37, Data: 13},
			},
		},
	},

	// version 23
	{
		LevelL: {
			Total: 1364,
			Data:  1094,
			Blocks: []blockCapacity{
				{Num: 4, Total: 151, Data: 121},
				{Num: 5, Total: 152, Data: 122},
			},
		},
		LevelM: {
			Total: 1364,
			Data:  860,
			Blocks: []blockCapacity{
				{Num: 4, Total: 75, Data: 47},
				{Num: 14, Total: 76, Data: 48},
			},
		},
		LevelQ: {
			Total: 1364,
			Data:  614,
			Blocks: []blockCapacity{
				{Num: 11, Total: 54, Data: 24},
				{Num: 14, Total: 55, Data: 25},
			},
		},
		LevelH: {
			Total: 1364,
			Data:  464,
			Blocks: []blockCapacity{
				{Num: 16, Total: 45, Data: 15},
				{Num: 14, Total: 46, Data: 16},
			},
		},
	},

	// version 24
	{
		LevelL: {
			Total: 1474,
			Data:  1174,
			Blocks: []blockCapacity{
				{Num: 6, Total: 147, Data: 117},
				{Num: 4, Total: 148, Data: 118},
			},
		},
		LevelM: {
			Total: 1474,
			Data:  914,
			Blocks: []blockCapacity{
				{Num: 6, Total: 73, Data: 45},
				{Num: 14, Total: 74, Data: 46},
			},
		},
		LevelQ: {
			Total: 1474,
			Data:  664,
			Blocks: []blockCapacity{
				{Num: 11, Total: 54, Data: 24},
				{Num: 16, Total: 55, Data: 25},
			},
		},
		LevelH: {
			Total: 1474,
			Data:  514,
			Blocks: []blockCapacity{
				{Num: 30, Total: 46, Data: 16},
				{Num: 2, Total: 47, Data: 17},
			},
		},
	},

	// version 25
	{
		LevelL: {
			Total: 26,
			Data:  1276,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1000,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  718,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  538,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 26
	{
		LevelL: {
			Total: 26,
			Data:  1370,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1062,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  754,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  596,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 27
	{
		LevelL: {
			Total: 26,
			Data:  1468,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1128,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  808,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  628,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 28
	{
		LevelL: {
			Total: 26,
			Data:  1531,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1193,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  871,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  661,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 29
	{
		LevelL: {
			Total: 26,
			Data:  1631,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1267,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  911,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  701,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 30
	{
		LevelL: {
			Total: 26,
			Data:  1735,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1373,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  985,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  745,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 31
	{
		LevelL: {
			Total: 26,
			Data:  1843,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1455,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1033,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  793,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 32
	{
		LevelL: {
			Total: 26,
			Data:  1955,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1541,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1115,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  845,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 33
	{
		LevelL: {
			Total: 26,
			Data:  2071,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1631,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1171,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  901,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 34
	{
		LevelL: {
			Total: 26,
			Data:  2191,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1725,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1231,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  961,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 35
	{
		LevelL: {
			Total: 26,
			Data:  2306,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1812,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1286,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  986,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 36
	{
		LevelL: {
			Total: 26,
			Data:  2434,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1914,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1354,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  1054,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 37
	{
		LevelL: {
			Total: 26,
			Data:  2566,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  1992,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1426,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  1096,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 38
	{
		LevelL: {
			Total: 26,
			Data:  2702,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  2102,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1502,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  1142,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 9},
			},
		},
	},

	// version 39
	{
		LevelL: {
			Total: 26,
			Data:  2812,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 19},
			},
		},
		LevelM: {
			Total: 26,
			Data:  2102,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 16},
			},
		},
		LevelQ: {
			Total: 26,
			Data:  1502,
			Blocks: []blockCapacity{
				{Num: 1, Total: 26, Data: 13},
			},
		},
		LevelH: {
			Total: 26,
			Data:  1142,
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
