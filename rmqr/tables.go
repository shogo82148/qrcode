package rmqr

type capacity struct {
	Total      int    // number of total code words
	Data       int    // number of data code words
	Correction int    // number of correction code words
	BitLength  [5]int // length of character count indicator
	Blocks     []blockCapacity
}

type blockCapacity struct {
	Num      int // number of blocks
	Total    int // number of total code words
	Data     int // number of data code words
	MaxError int // maximum number of code word errors
	Reserved int // number of code words reserved for error detection
}

var capacityOrder = []Version{
	R7x43,
	R11x27,
	R7x59,
	R13x27,
	R9x43,
	R7x77,
	R9x59,
	R11x43,
	R13x43,
	R7x99,
	R15x43,
	R11x59,
	R9x77,
	R13x59,
	R17x43,
	R9x99,
	R11x77,
	R7x139,
	R15x59,
	R17x59,
	R11x99,
	R13x77,
	R15x77,
	R9x139,
	R13x99,
	R17x77,
	R11x139,
	R15x99,
	R13x139,
	R17x99,
	R15x139,
	R17x139,
}

var capacityTable = [32][2]capacity{
	// R7x43
	{
		LevelM: {
			Total:      13,
			Data:       6,
			Correction: 7,
			BitLength: [5]int{
				ModeNumeric:      4,
				ModeAlphanumeric: 3,
				ModeBytes:        3,
				ModeKanji:        2,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    13,
					Data:     6,
					MaxError: 3,
				},
			},
		},
		LevelH: {
			Total:      13,
			Data:       3,
			Correction: 10,
			BitLength: [5]int{
				ModeNumeric:      4,
				ModeAlphanumeric: 3,
				ModeBytes:        3,
				ModeKanji:        2,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    13,
					Data:     3,
					MaxError: 4,
				},
			},
		},
	},

	// R7x59
	{
		LevelM: {
			Total:      21,
			Data:       12,
			Correction: 9,
			BitLength: [5]int{
				ModeNumeric:      5,
				ModeAlphanumeric: 5,
				ModeBytes:        4,
				ModeKanji:        3,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    21,
					Data:     12,
					MaxError: 4,
				},
			},
		},
		LevelH: {
			Total:      21,
			Data:       7,
			Correction: 14,
			BitLength: [5]int{
				ModeNumeric:      5,
				ModeAlphanumeric: 5,
				ModeBytes:        4,
				ModeKanji:        3,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    21,
					Data:     7,
					MaxError: 6,
				},
			},
		},
	},

	// R7x77
	{
		LevelM: {
			Total:      32,
			Data:       20,
			Correction: 12,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 5,
				ModeBytes:        5,
				ModeKanji:        4,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    32,
					Data:     20,
					MaxError: 5,
				},
			},
		},
		LevelH: {
			Total:      32,
			Data:       10,
			Correction: 22,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 5,
				ModeBytes:        5,
				ModeKanji:        4,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    32,
					Data:     10,
					MaxError: 10,
				},
			},
		},
	},

	// R7x99
	{
		LevelM: {
			Total:      44,
			Data:       28,
			Correction: 16,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    44,
					Data:     28,
					MaxError: 7,
				},
			},
		},
		LevelH: {
			Total:      44,
			Data:       14,
			Correction: 30,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    44,
					Data:     14,
					MaxError: 14,
				},
			},
		},
	},

	// R7x139
	{
		LevelM: {
			Total:      68,
			Data:       44,
			Correction: 24,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    68,
					Data:     44,
					MaxError: 11,
				},
			},
		},
		LevelH: {
			Total:      68,
			Data:       24,
			Correction: 44,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    34,
					Data:     12,
					MaxError: 10,
				},
			},
		},
	},

	// R9x43
	{
		LevelM: {
			Total:      21,
			Data:       12,
			Correction: 9,
			BitLength: [5]int{
				ModeNumeric:      5,
				ModeAlphanumeric: 5,
				ModeBytes:        4,
				ModeKanji:        3,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    21,
					Data:     12,
					MaxError: 4,
				},
			},
		},
		LevelH: {
			Total:      21,
			Data:       7,
			Correction: 14,
			BitLength: [5]int{
				ModeNumeric:      5,
				ModeAlphanumeric: 5,
				ModeBytes:        4,
				ModeKanji:        3,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    21,
					Data:     7,
					MaxError: 6,
				},
			},
		},
	},

	// R9x59
	{
		LevelM: {
			Total:      33,
			Data:       21,
			Correction: 12,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 5,
				ModeBytes:        5,
				ModeKanji:        4,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    33,
					Data:     21,
					MaxError: 5,
				},
			},
		},
		LevelH: {
			Total:      33,
			Data:       11,
			Correction: 22,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 5,
				ModeBytes:        5,
				ModeKanji:        4,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    33,
					Data:     11,
					MaxError: 10,
				},
			},
		},
	},

	// R9x77
	{
		LevelM: {
			Total:      49,
			Data:       31,
			Correction: 18,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    49,
					Data:     31,
					MaxError: 8,
				},
			},
		},
		LevelH: {
			Total:      49,
			Data:       17,
			Correction: 32,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    24,
					Data:     8,
					MaxError: 7,
				},
				{
					Num:      1,
					Total:    25,
					Data:     9,
					MaxError: 7,
				},
			},
		},
	},

	// R9x99
	{
		LevelM: {
			Total:      66,
			Data:       42,
			Correction: 24,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    66,
					Data:     42,
					MaxError: 11,
				},
			},
		},
		LevelH: {
			Total:      66,
			Data:       22,
			Correction: 44,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    33,
					Data:     11,
					MaxError: 10,
				},
			},
		},
	},

	// R9x139
	{
		LevelM: {
			Total:      99,
			Data:       63,
			Correction: 36,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    49,
					Data:     31,
					MaxError: 8,
				},
				{
					Num:      1,
					Total:    50,
					Data:     32,
					MaxError: 8,
				},
			},
		},
		LevelH: {
			Total:      99,
			Data:       33,
			Correction: 66,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      3,
					Total:    33,
					Data:     11,
					MaxError: 10,
				},
			},
		},
	},

	// R11x27
	{
		LevelM: {
			Total:      15,
			Data:       7,
			Correction: 8,
			BitLength: [5]int{
				ModeNumeric:      4,
				ModeAlphanumeric: 4,
				ModeBytes:        3,
				ModeKanji:        2,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    15,
					Data:     7,
					MaxError: 3,
				},
			},
		},
		LevelH: {
			Total:      15,
			Data:       5,
			Correction: 10,
			BitLength: [5]int{
				ModeNumeric:      4,
				ModeAlphanumeric: 4,
				ModeBytes:        3,
				ModeKanji:        2,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    15,
					Data:     5,
					MaxError: 4,
				},
			},
		},
	},

	// R11x43
	{
		LevelM: {
			Total:      31,
			Data:       19,
			Correction: 12,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 5,
				ModeBytes:        5,
				ModeKanji:        4,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    31,
					Data:     19,
					MaxError: 5,
				},
			},
		},
		LevelH: {
			Total:      31,
			Data:       11,
			Correction: 20,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 5,
				ModeBytes:        5,
				ModeKanji:        4,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    31,
					Data:     11,
					MaxError: 9,
				},
			},
		},
	},

	// R11x59
	{
		LevelM: {
			Total:      47,
			Data:       31,
			Correction: 16,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    47,
					Data:     31,
					MaxError: 7,
				},
			},
		},
		LevelH: {
			Total:      47,
			Data:       15,
			Correction: 32,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    23,
					Data:     7,
					MaxError: 7,
				},
				{
					Num:      1,
					Total:    24,
					Data:     8,
					MaxError: 7,
				},
			},
		},
	},

	// R11x77
	{
		LevelM: {
			Total:      67,
			Data:       43,
			Correction: 24,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    67,
					Data:     43,
					MaxError: 11,
				},
			},
		},
		LevelH: {
			Total:      67,
			Data:       23,
			Correction: 44,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    33,
					Data:     11,
					MaxError: 10,
				},
				{
					Num:      1,
					Total:    34,
					Data:     12,
					MaxError: 10,
				},
			},
		},
	},

	// R11x99
	{
		LevelM: {
			Total:      89,
			Data:       57,
			Correction: 32,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    44,
					Data:     28,
					MaxError: 7,
				},
				{
					Num:      1,
					Total:    45,
					Data:     29,
					MaxError: 7,
				},
			},
		},
		LevelH: {
			Total:      89,
			Data:       29,
			Correction: 60,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    44,
					Data:     14,
					MaxError: 14,
				},
				{
					Num:      1,
					Total:    45,
					Data:     15,
					MaxError: 14,
				},
			},
		},
	},

	// R11x139
	{
		LevelM: {
			Total:      132,
			Data:       84,
			Correction: 48,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    66,
					Data:     42,
					MaxError: 11,
				},
			},
		},
		LevelH: {
			Total:      132,
			Data:       42,
			Correction: 90,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      3,
					Total:    44,
					Data:     14,
					MaxError: 14,
				},
			},
		},
	},

	// R13x27
	{
		LevelM: {
			Total:      21,
			Data:       14,
			Correction: 7,
			BitLength: [5]int{
				ModeNumeric:      4,
				ModeAlphanumeric: 3,
				ModeBytes:        3,
				ModeKanji:        2,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    21,
					Data:     14,
					MaxError: 3,
				},
			},
		},
		LevelH: {
			Total:      21,
			Data:       7,
			Correction: 14,
			BitLength: [5]int{
				ModeNumeric:      4,
				ModeAlphanumeric: 3,
				ModeBytes:        3,
				ModeKanji:        2,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    21,
					Data:     7,
					MaxError: 6,
				},
			},
		},
	},

	// R13x43
	{
		LevelM: {
			Total:      41,
			Data:       27,
			Correction: 14,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    41,
					Data:     27,
					MaxError: 6,
				},
			},
		},
		LevelH: {
			Total:      41,
			Data:       13,
			Correction: 28,
			BitLength: [5]int{
				ModeNumeric:      6,
				ModeAlphanumeric: 6,
				ModeBytes:        5,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    41,
					Data:     13,
					MaxError: 13,
				},
			},
		},
	},

	// R13x59
	{
		LevelM: {
			Total:      60,
			Data:       38,
			Correction: 22,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    60,
					Data:     38,
					MaxError: 10,
				},
			},
		},
		LevelH: {
			Total:      60,
			Data:       20,
			Correction: 40,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    30,
					Data:     10,
					MaxError: 9,
				},
			},
		},
	},

	// R13x77
	{
		LevelM: {
			Total:      85,
			Data:       53,
			Correction: 32,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    42,
					Data:     26,
					MaxError: 7,
				},
				{
					Num:      1,
					Total:    43,
					Data:     27,
					MaxError: 7,
				},
			},
		},
		LevelH: {
			Total:      85,
			Data:       29,
			Correction: 56,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    42,
					Data:     14,
					MaxError: 13,
				},
				{
					Num:      1,
					Total:    43,
					Data:     15,
					MaxError: 13,
				},
			},
		},
	},

	// R13x99
	{
		LevelM: {
			Total:      113,
			Data:       73,
			Correction: 40,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    56,
					Data:     36,
					MaxError: 9,
				},
				{
					Num:      1,
					Total:    57,
					Data:     37,
					MaxError: 9,
				},
			},
		},
		LevelH: {
			Total:      113,
			Data:       35,
			Correction: 78,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    37,
					Data:     11,
					MaxError: 12,
				},
				{
					Num:      2,
					Total:    38,
					Data:     12,
					MaxError: 12,
				},
			},
		},
	},

	// R13x139
	{
		LevelM: {
			Total:      166,
			Data:       106,
			Correction: 60,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 8,
				ModeBytes:        7,
				ModeKanji:        7,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    55,
					Data:     35,
					MaxError: 9,
				},
				{
					Num:      1,
					Total:    56,
					Data:     36,
					MaxError: 9,
				},
			},
		},
		LevelH: {
			Total:      166,
			Data:       54,
			Correction: 112,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 8,
				ModeBytes:        7,
				ModeKanji:        7,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    41,
					Data:     13,
					MaxError: 13,
				},
				{
					Num:      2,
					Total:    42,
					Data:     14,
					MaxError: 13,
				},
			},
		},
	},

	// R15x43
	{
		LevelM: {
			Total:      51,
			Data:       33,
			Correction: 18,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    51,
					Data:     33,
					MaxError: 8,
				},
			},
		},
		LevelH: {
			Total:      51,
			Data:       15,
			Correction: 36,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    25,
					Data:     7,
					MaxError: 8,
				},
				{
					Num:      1,
					Total:    26,
					Data:     8,
					MaxError: 8,
				},
			},
		},
	},

	// R15x59
	{
		LevelM: {
			Total:      74,
			Data:       48,
			Correction: 26,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    74,
					Data:     48,
					MaxError: 12,
				},
			},
		},
		LevelH: {
			Total:      74,
			Data:       26,
			Correction: 48,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    37,
					Data:     13,
					MaxError: 11,
				},
			},
		},
	},

	// R15x77
	{
		LevelM: {
			Total:      103,
			Data:       67,
			Correction: 36,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    51,
					Data:     33,
					MaxError: 8,
				},
				{
					Num:      1,
					Total:    52,
					Data:     34,
					MaxError: 8,
				},
			},
		},
		LevelH: {
			Total:      103,
			Data:       31,
			Correction: 72,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    34,
					Data:     10,
					MaxError: 11,
				},
				{
					Num:      1,
					Total:    35,
					Data:     11,
					MaxError: 11,
				},
			},
		},
	},

	// R15x99
	{
		LevelM: {
			Total:      136,
			Data:       88,
			Correction: 48,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    68,
					Data:     44,
					MaxError: 11,
				},
			},
		},
		LevelH: {
			Total:      136,
			Data:       48,
			Correction: 88,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      4,
					Total:    34,
					Data:     12,
					MaxError: 10,
				},
			},
		},
	},

	// R15x139
	{
		LevelM: {
			Total:      199,
			Data:       127,
			Correction: 72,
			BitLength: [5]int{
				ModeNumeric:      9,
				ModeAlphanumeric: 8,
				ModeBytes:        7,
				ModeKanji:        7,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    66,
					Data:     42,
					MaxError: 11,
				},
				{
					Num:      1,
					Total:    67,
					Data:     43,
					MaxError: 11,
				},
			},
		},
		LevelH: {
			Total:      199,
			Data:       69,
			Correction: 130,
			BitLength: [5]int{
				ModeNumeric:      9,
				ModeAlphanumeric: 8,
				ModeBytes:        7,
				ModeKanji:        7,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    39,
					Data:     13,
					MaxError: 12,
				},
				{
					Num:      4,
					Total:    40,
					Data:     14,
					MaxError: 12,
				},
			},
		},
	},

	// R17x43
	{
		LevelM: {
			Total:      61,
			Data:       39,
			Correction: 22,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    61,
					Data:     39,
					MaxError: 10,
				},
			},
		},
		LevelH: {
			Total:      61,
			Data:       21,
			Correction: 40,
			BitLength: [5]int{
				ModeNumeric:      7,
				ModeAlphanumeric: 6,
				ModeBytes:        6,
				ModeKanji:        5,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    30,
					Data:     10,
					MaxError: 9,
				},
				{
					Num:      1,
					Total:    31,
					Data:     11,
					MaxError: 9,
				},
			},
		},
	},

	// R17x59
	{
		LevelM: {
			Total:      88,
			Data:       56,
			Correction: 32,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    44,
					Data:     28,
					MaxError: 7,
				},
			},
		},
		LevelH: {
			Total:      88,
			Data:       28,
			Correction: 60,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        6,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    44,
					Data:     14,
					MaxError: 14,
				},
			},
		},
	},

	// R17x77
	{
		LevelM: {
			Total:      122,
			Data:       78,
			Correction: 44,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    61,
					Data:     39,
					MaxError: 10,
				},
			},
		},
		LevelH: {
			Total:      122,
			Data:       38,
			Correction: 84,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 7,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      1,
					Total:    40,
					Data:     12,
					MaxError: 13,
				},
				{
					Num:      2,
					Total:    41,
					Data:     13,
					MaxError: 13,
				},
			},
		},
	},

	// R17x99
	{
		LevelM: {
			Total:      160,
			Data:       100,
			Correction: 60,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 8,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    53,
					Data:     33,
					MaxError: 9,
				},
				{
					Num:      1,
					Total:    54,
					Data:     34,
					MaxError: 9,
				},
			},
		},
		LevelH: {
			Total:      160,
			Data:       56,
			Correction: 104,
			BitLength: [5]int{
				ModeNumeric:      8,
				ModeAlphanumeric: 8,
				ModeBytes:        7,
				ModeKanji:        6,
			},
			Blocks: []blockCapacity{
				{
					Num:      4,
					Total:    40,
					Data:     14,
					MaxError: 12,
				},
			},
		},
	},

	// R17x139
	{
		LevelM: {
			Total:      232,
			Data:       152,
			Correction: 80,
			BitLength: [5]int{
				ModeNumeric:      9,
				ModeAlphanumeric: 8,
				ModeBytes:        8,
				ModeKanji:        7,
			},
			Blocks: []blockCapacity{
				{
					Num:      4,
					Total:    58,
					Data:     38,
					MaxError: 9,
				},
			},
		},
		LevelH: {
			Total:      232,
			Data:       76,
			Correction: 156,
			BitLength: [5]int{
				ModeNumeric:      9,
				ModeAlphanumeric: 8,
				ModeBytes:        8,
				ModeKanji:        7,
			},
			Blocks: []blockCapacity{
				{
					Num:      2,
					Total:    38,
					Data:     12,
					MaxError: 12,
				},
				{
					Num:      4,
					Total:    39,
					Data:     13,
					MaxError: 12,
				},
			},
		},
	},
}
