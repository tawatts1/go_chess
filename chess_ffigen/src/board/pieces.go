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

var CastleMap = map[rune]rune{
	WhiteRookC: WhiteRookNC,
	BlackRookC: BlackRookNC,
	WhiteKing:  WhiteRookC,
	BlackKing:  BlackRookC,
}

var EnPassantMap = map[rune]rune{
	White:       BlackPawnEP, //on white's move, check for black pawns which have just two stepped.
	Black:       WhitePawnEP, // and visa versa
	WhitePawnEP: WhitePawn,
	BlackPawnEP: BlackPawn,
	WhitePawn:   WhitePawnEP,
	BlackPawn:   BlackPawnEP,
}

func IsPawn(p rune) bool {
	switch p {
	case BlackPawn:
		return true
	case BlackPawnEP:
		return true
	case WhitePawn:
		return true
	case WhitePawnEP:
		return true
	default:
		return false
	}
}

func IsEPPawn(p rune) bool {
	if p == BlackPawnEP || p == WhitePawnEP {
		return true
	}
	return false
}
func IsKnight(p rune) bool {
	if p == BlackKnight || p == WhiteKnight {
		return true
	}
	return false
}
func IsBishop(p rune) bool {
	if p == BlackBishop || p == WhiteBishop {
		return true
	}
	return false
}
func IsRook(p rune) bool {
	if p == BlackRookC || p == WhiteRookC ||
		p == BlackRookNC || p == WhiteRookNC {
		return true
	}
	return false
}

func IsRookCastleable(p rune) bool {
	return p == BlackRookC || p == WhiteRookC
}
func IsQueen(p rune) bool {
	if p == BlackQueen || p == WhiteQueen {
		return true
	}
	return false
}

func IsKing(p rune) bool {
	if p == BlackKing || p == WhiteKing {
		return true
	}
	return false
}

const EnPassant rune = 'e'
const CastleBridge rune = 'c'

var BlackPawnPromotion = []rune{BlackQueen, BlackKnight, BlackBishop, BlackRookNC}
var WhitePawnPromotion = []rune{WhiteQueen, WhiteKnight, WhiteBishop, WhiteRookNC}
