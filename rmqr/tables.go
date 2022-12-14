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
}

var capacityTable = [32][2]capacity{
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
			// TODO
		},
	},
}
