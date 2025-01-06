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

func (b board) GetMoves(friends, enemies map[coord]bool, c coord) []move {
	piece := b.GetPiece(c)
	if IsPawn(piece) {
		heading := -1
		if IsBlack(piece) {
			heading = 1
		}
		return b.GetPawnMoves(friends, enemies, c, heading)
	} else if IsBishop(piece) {
		return GetBishopMoves(friends, enemies, c)
	} else if IsRook(piece) {
		return GetRookMoves(friends, enemies, c)
	} else if IsQueen(piece) {
		return GetQueenMoves(friends, enemies, c)
	} else if IsKnight(piece) {
		return GetKnightMoves(friends, enemies, c)
	} else if IsKing(piece) {
		return GetKingMoves(friends, enemies, c)
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

func GetBishopMoves(friends, enemies map[coord]bool, c coord) []move {
	out := make([]move, 0, 3)
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

func GetRookMoves(friends, enemies map[coord]bool, c coord) []move {
	out := make([]move, 0, 3)
	var newSquare coord
	// vectors used to add diagonal moves to the starting square
	vectors := []coord{{y: 1, x: 0}, {y: -1, x: 0}, {y: 0, x: 1}, {y: 0, x: -1}}
	for _, v := range vectors {
		for newSquare = c.Copy().Add(v.y, v.x); newSquare.IsInBoard(); newSquare = newSquare.Copy().Add(v.y, v.x) {
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
	return out
}

func GetQueenMoves(friends, enemies map[coord]bool, c coord) []move {
	return append(GetBishopMoves(friends, enemies, c), GetRookMoves(friends, enemies, c)...)
}

func GetKnightMoves(friends, enemies map[coord]bool, c coord) []move {
	out := make([]move, 0, 4)
	var newSquare coord
	for _, sign := range [2]int{-1, 1} {
		for _, vector := range [4]coord{{y: 1, x: 2}, {y: 2, x: 1}, {y: -1, x: 2}, {y: -2, x: 1}} {
			newSquare = c.Add(sign*vector.y, sign*vector.x)
			if newSquare.IsInBoard() {
				_, hasFriend := friends[newSquare]
				if !hasFriend {
					out = append(out, move{a: c, b: newSquare})
				}
			}
		}
	}
	return out
}

func GetKingMoves(friends, enemies map[coord]bool, c coord) []move {
	out := make([]move, 0, 2)
	var newSquare coord
	for _, sign := range [2]int{-1, 1} {
		for _, vector := range [4]coord{{y: 1, x: 1}, {y: 1, x: 0}, {y: 0, x: 1}, {y: -1, x: 1}} {
			newSquare = c.Add(sign*vector.y, sign*vector.x)
			if newSquare.IsInBoard() {
				_, hasFriend := friends[newSquare]
				if !hasFriend {
					out = append(out, move{a: c, b: newSquare})
				}
			}
		}
	}
	return out
}

func (b board) IsInCheck(friends, enemies map[coord]bool) bool {
	king_coord := b.GetKingCoord(friends)
	for enemy_coord := range enemies {
		// note that we are checking what moves the enemy has, so the friends and enemies maps are switched.
		for _, m := range b.GetMoves(enemies, friends, enemy_coord) {
			if king_coord.Equals(m.b) {
				return true
			}
		}
	}
	return false
}
