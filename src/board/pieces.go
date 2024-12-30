package board

const Black, White rune = 'b', 'w'
const Space rune = '0'
const BlackPawn, WhitePawn rune = 'p', 'P'
const BlackKnight, WhiteKnight rune = 'n', 'N'
const BlackBishop, WhiteBishop rune = 'b', 'B'
const BlackRookNC, WhiteRookNC rune = 'r', 'R'
const BlackQueen, WhiteQueen rune = 'q', 'Q'
const BlackKing, WhiteKing rune = 'k', 'K'
const BlackRookC, WhiteRookC rune = 'o', 'O'
const BlackPawnEP, WhitePawnEP rune = 'a', 'A'

// A map that maps the encoded board runes to human readable strings
var humanMap = map[rune]string{
	Space:       "  ",
	BlackPawn:   "bp",
	BlackKnight: "bn",
	BlackBishop: "bb",
	BlackRookNC: "br",
	BlackQueen:  "bq",
	BlackKing:   "bk",
	BlackRookC:  "BR", // a rook with the ability to castle
	BlackPawnEP: "BP", // a pawn which can be taken via en passant
	//--------
	WhitePawn:   "wp",
	WhiteKnight: "wn",
	WhiteBishop: "wb",
	WhiteRookNC: "wr",
	WhiteQueen:  "wq",
	WhiteKing:   "wk",
	WhiteRookC:  "WR", // a rook with the ability to castle
	WhitePawnEP: "WP", // a pawn which can be taken via en passant
}

var blackMap = map[rune]bool{
	BlackPawn:   true,
	BlackKnight: true,
	BlackBishop: true,
	BlackRookNC: true,
	BlackQueen:  true,
	BlackKing:   true,
	BlackRookC:  true,
	BlackPawnEP: true,
}

var whiteMap = map[rune]bool{
	WhitePawn:   true,
	WhiteKnight: true,
	WhiteBishop: true,
	WhiteRookNC: true,
	WhiteQueen:  true,
	WhiteKing:   true,
	WhiteRookC:  true,
	WhitePawnEP: true,
}
