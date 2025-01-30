package ai

import "github.com/tawatts1/go_chess/board"

var ScoringDefaultPieceValue = "simple"
var ScoringHighKnightPieceValue = "highKnight"

func getScoringFunction(scoringName string) func(board.Board, bool) float64 {
	if scoringName == ScoringDefaultPieceValue {
		return getScoringFuncFromPieceValueMap(defaultPieceValue)
	} else if scoringName == ScoringHighKnightPieceValue {
		return getScoringFuncFromPieceValueMap(highKnightPieceValue)
	} else {
		panic("Unexpected scoring function name")
	}
}

func getScoringFuncFromPieceValueMap(pieceValueMap map[rune]float64) func(board.Board, bool) float64 {
	return func(b board.Board, isWhite bool) float64 {
		var out float64
		var multiplier float64 = 1
		if isWhite {
			multiplier = -1
		}
		for _, c := range board.AllCoordinates {
			out += multiplier * pieceValueMap[b.GetPiece(c)]
		}
		return out
	}
}

var defaultPieceValue = map[rune]float64{
	board.BlackPawn:   1,
	board.BlackKnight: 2.6,
	board.BlackBishop: 3,
	board.BlackRookNC: 4.99,
	board.BlackQueen:  9,
	board.BlackKing:   0,
	board.BlackRookC:  5, // a rook with the ability to castle
	board.BlackPawnEP: 1, // a pawn which can be taken via en passant
	//--------
	board.WhitePawn:   -1,
	board.WhiteKnight: -2.6,
	board.WhiteBishop: -3,
	board.WhiteRookNC: -4.99,
	board.WhiteQueen:  -9,
	board.WhiteKing:   0,
	board.WhiteRookC:  -5, // a rook with the ability to castle
	board.WhitePawnEP: -1, // a pawn which can be taken via en passant
}

var highKnightPieceValue = map[rune]float64{
	board.BlackPawn:   1,
	board.BlackKnight: 3.3,
	board.BlackBishop: 3,
	board.BlackRookNC: 4.99,
	board.BlackQueen:  9,
	board.BlackKing:   0,
	board.BlackRookC:  5, // a rook with the ability to castle
	board.BlackPawnEP: 1, // a pawn which can be taken via en passant
	//--------
	board.WhitePawn:   -1,
	board.WhiteKnight: -3.3,
	board.WhiteBishop: -3,
	board.WhiteRookNC: -4.99,
	board.WhiteQueen:  -9,
	board.WhiteKing:   0,
	board.WhiteRookC:  -5, // a rook with the ability to castle
	board.WhitePawnEP: -1, // a pawn which can be taken via en passant
}
