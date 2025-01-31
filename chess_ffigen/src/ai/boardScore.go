package ai

import "github.com/tawatts1/go_chess/board"

var ScoringDefaultPieceValue = "simple"
var ScoringHighKnightPieceValue = "highKnight"
var ScoringPiecePositionValue = "position"
var positionScalingFactor float64 = 0.01

func getScoringFunction(scoringName string) func(board.Board, bool) float64 {
	if scoringName == ScoringDefaultPieceValue {
		return getScoringFuncFromPieceValueMap(defaultPieceValue)
	} else if scoringName == ScoringHighKnightPieceValue {
		return getScoringFuncFromPieceValueMap(highKnightPieceValue)
	} else if scoringName == ScoringPiecePositionValue {
		return getPositionScoreValue
	} else {
		panic("Unexpected scoring function name")
	}
}

func getScoringFuncFromPieceValueMap(pieceValueMap map[rune]float64) func(board.Board, bool) float64 {
	return func(b board.Board, isWhite bool) float64 {
		var out float64
		var multiplier float64 = 1
		if !isWhite {
			multiplier = -1
		}
		var p rune
		for _, c := range board.AllCoordinates {
			p = b.GetPiece(c)
			if p != board.Space {
				out += multiplier * pieceValueMap[p]
			}
		}
		return out
	}
}

func getPositionScoreValue(b board.Board, isWhitesMove bool) float64 {
	var out float64
	var moveColorMultiplier float64 = 1
	if !isWhitesMove {
		moveColorMultiplier = -1
	}
	var p rune
	var pieceColorMultiplier float64 = 1
	for _, c := range board.AllCoordinates {
		p = b.GetPiece(c)
		if board.IsWhite(p) {
			pieceColorMultiplier = 1
		} else {
			pieceColorMultiplier = -1
		}
		if p != board.Space {
			//add the inherent value of the piece
			//the piece color multiplier is built into the map.
			out += moveColorMultiplier * defaultPieceValue[p]
			//calculate and add the value based on the pieces position
			out += moveColorMultiplier * pieceColorMultiplier * positionScalingFactor * GetPositionScore(c, p)
		}
	}
	return out
}

// Returns a positive value no matter the piece color.
func GetPositionScore(c board.Coord, p rune) float64 {
	var out float64
	if board.IsPawn(p) {
		//Pawns are more valuable if they are close to promotion.
		y := c.GetY()
		if board.IsWhite(p) {
			y = 7 - y // normalize for color
		}
		if y == 1 {
			out = 0.5
		} else {
			out = float64(y - 2)
		}

	} else if board.IsBishop(p) || board.IsQueen(p) {
		y, x := c.GetY(), c.GetX()
		min, max := minimum(y, x), -minimum(-y, -x)
		distToEdge := minimum(7-max, min)
		out = float64(distToEdge) * 0.5
		if (y == x || 7-y == x) && distToEdge != 3 {
			//favor the diagonals for bishops/queens
			out += 0.5
		}
	} else if board.IsKnight(p) {
		y, x := c.GetY(), c.GetX()
		isInHomePosition := false
		if y == 0 || y == 7 {
			if x == 1 || x == 6 {
				isInHomePosition = true
			}
		}
		//flip x and y to upper left corner of board.
		if y > 3 {
			y = 7 - y
		}
		if x > 3 {
			x = 7 - x
		}
		if y > 1 && x > 1 {
			out = 1
		} else if y > 1 || x > 1 {
			if y == 0 || x == 0 {
				out = .5
			} else {
				out = .75
			}
		} else {
			if isInHomePosition {
				out = .5
			} else {
				out = float64(x+y+2) / 8.0
			}
		}
	}
	// if board.IsWhite(p) {
	// 	return out
	// } else {
	// 	return out * -1
	// }
	return out
}

func minimum(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

var defaultPieceValue = map[rune]float64{
	board.BlackPawn:   -1,
	board.BlackKnight: -2.6,
	board.BlackBishop: -3,
	board.BlackRookNC: -4.99,
	board.BlackQueen:  -9,
	board.BlackKing:   0,
	board.BlackRookC:  -5, // a rook with the ability to castle
	board.BlackPawnEP: -1, // a pawn which can be taken via en passant
	//--------
	board.WhitePawn:   1,
	board.WhiteKnight: 2.6,
	board.WhiteBishop: 3,
	board.WhiteRookNC: 4.99,
	board.WhiteQueen:  9,
	board.WhiteKing:   0,
	board.WhiteRookC:  5, // a rook with the ability to castle
	board.WhitePawnEP: 1, // a pawn which can be taken via en passant
}

var highKnightPieceValue = map[rune]float64{
	board.BlackPawn:   -1,
	board.BlackKnight: -3.3,
	board.BlackBishop: -3,
	board.BlackRookNC: -4.99,
	board.BlackQueen:  -9,
	board.BlackKing:   0,
	board.BlackRookC:  -5, // a rook with the ability to castle
	board.BlackPawnEP: -1, // a pawn which can be taken via en passant
	//--------
	board.WhitePawn:   1,
	board.WhiteKnight: 3.3,
	board.WhiteBishop: 3,
	board.WhiteRookNC: 4.99,
	board.WhiteQueen:  9,
	board.WhiteKing:   0,
	board.WhiteRookC:  5, // a rook with the ability to castle
	board.WhitePawnEP: 1, // a pawn which can be taken via en passant
}
