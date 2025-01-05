package board

import "fmt"

// A way to encode a move via a move from one coordinate to another
type move struct {
	a, b    coord
	special rune
}

func (m1 move) Equals(m2 move) bool {
	return m1.a.Equals(m2.a) && m1.b.Equals(m2.b) && m1.special == m2.special
}

func AnyEqual(moves []move, m move) bool {
	for _, m2 := range moves {
		if m2.Equals(m) {
			return true
		}
	}
	return false
}

func (m move) String() string {
	return fmt.Sprintf("{%v, %v, %c }", m.a, m.b, m.special)
}

const EnPassant rune = 'e'
const CastleBridge rune = 'c'

var BlackPawnPromotion = &[]rune{BlackKnight, BlackBishop, BlackRookNC, BlackQueen}
var WhitePawnPromotion = &[]rune{WhiteKnight, WhiteBishop, WhiteRookNC, WhiteQueen}

func (b board) GetMoves(c coord, bcm, wcm map[coord]bool) []move {
	piece := b.GetPiece(c)
	blk := IsBlack(piece)
	if IsPawn(piece) {
		if blk {
			heading := 1
			return b.GetPawnMoves(bcm, wcm, c, heading)
		} else {
			heading := -1
			return b.GetPawnMoves(wcm, bcm, c, heading)
		}
	} else if IsBishop(piece) {
		if blk {
			return b.GetBishopMoves(bcm, wcm, c)
		} else {
			return b.GetBishopMoves(wcm, bcm, c)
		}
	} else {
		panic("Not implemented")
	}
}

func (b board) GetPawnMoves(friends, enemies map[coord]bool, c coord, heading int) []move {
	out := make([]move, 0, 2)
	isBlack := heading == 1
	//directly ahead of pawn
	forward := c.Copy().Add(heading, 0)
	_, hasFriend := friends[forward]
	_, hasEnemy := enemies[forward]
	if !hasFriend && !hasEnemy {
		if (isBlack && c.y == 6) ||
			(!isBlack && c.y == 1) {
			var promotions *[]rune
			if isBlack {
				promotions = BlackPawnPromotion
			} else {
				promotions = WhitePawnPromotion
			}
			out = append(out, getPawnPromotionMoves(c, forward, promotions)...)
		} else {
			out = append(out, move{a: c, b: forward})
		}

		// two steps ahead of pawn, for opening
		if (!isBlack && c.y == 6) ||
			(isBlack && c.y == 1) {
			forward2 := c.Copy().Add(heading*2, 0)
			_, hasFriend2 := friends[forward2]
			_, hasEnemy2 := enemies[forward2]
			if !hasFriend2 && !hasEnemy2 {
				out = append(out, move{a: c, b: forward2})
			}
		}
	}
	//each En Passant
	for _, lr := range [2]int{-1, 1} {
		lrSquare := c.Copy().Add(0, lr)
		_, lrEnemy := enemies[lrSquare]
		if lrEnemy && IsEPPawn(b.GetPiece(lrSquare)) {
			lrSquare = lrSquare.Add(heading, 0)
			out = append(out, move{a: c, b: lrSquare, special: EnPassant})
		}
	}
	//each diagonal attack
	for _, lr := range [2]int{-1, 1} {
		diagonalSquare := c.Copy().Add(heading, lr)
		_, hasEnemy := enemies[diagonalSquare]
		if hasEnemy {
			if (isBlack && c.y == 6) ||
				(!isBlack && c.y == 1) {
				var promotions *[]rune
				if isBlack {
					promotions = BlackPawnPromotion
				} else {
					promotions = WhitePawnPromotion
				}
				out = append(out, getPawnPromotionMoves(c, diagonalSquare, promotions)...)
			} else {
				out = append(out, move{a: c, b: diagonalSquare})
			}
		}
	}
	return out
}

func getPawnPromotionMoves(from, to coord, promotionCodes *[]rune) []move {
	//when promoted, pawns can become knight, bishop, rook or queen.
	numPromotions := len(*promotionCodes)
	var out = make([]move, numPromotions, numPromotions)
	for i := range numPromotions {
		out[i] = move{a: from, b: to, special: (*promotionCodes)[i]}
	}
	return out
}

func (b board) GetBishopMoves(friends, enemies map[coord]bool, c coord) []move {
	out := make([]move, 0, 2)
	var newSquare coord
	// vectors used to add diagonal moves to the starting square
	v_xs, v_ys := [2]int{-1, 1}, [2]int{-1, 1}
	for _, v_x := range v_xs {
		for _, v_y := range v_ys {
			for newSquare = c.Copy().Add(v_y, v_x); newSquare.IsInBoard(); newSquare = newSquare.Copy().Add(v_y, v_x) {
				_, hasFriend := friends[newSquare]
				if hasFriend {
					break // with the other vectors
				}
				//add move whether there is an enemy there or it is empty
				out = append(out, move{a: c, b: newSquare})
				_, hasEnemy := enemies[newSquare]
				if hasEnemy {
					break // with the other vectors
				}
			}
		}
	}
	return out
}
