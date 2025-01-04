package board

import "fmt"

// A way to encode a move via a move from one coordinate to another
type move struct {
	a, b    *coord
	special rune
}

func (m1 move) Equals(m2 move) bool {
	return m1.a.Equals(*(m2.a)) && m1.b.Equals(*(m2.b)) && m1.special == m2.special
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
const PawnPromotion rune = 'p'
const CastleBridge rune = 'c'

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
	} else {
		panic("Not implemented")
	}
}

func (b board) GetPawnMoves(friends, enemies map[coord]bool, c coord, heading int) []move {
	out := make([]move, 0, 2)
	//directly ahead of pawn
	forward := c.Copy().Add(heading, 0)
	_, hasFriend := friends[forward]
	_, hasEnemy := enemies[forward]
	if !hasFriend && !hasEnemy {
		out = append(out, move{a: &c, b: &forward})
		// two steps ahead of pawn, for opening
		if (heading == -1 && c.y == 6) ||
			(heading == 1 && c.y == 1) {
			forward2 := c.Copy().Add(heading*2, 0)
			_, hasFriend2 := friends[forward2]
			_, hasEnemy2 := enemies[forward2]
			if !hasFriend2 && !hasEnemy2 {
				out = append(out, move{a: &c, b: &forward2})
			}
		}
	}
	//each En Passant
	for _, lr := range [2]int{-1, 1} {
		lrSquare := c.Copy().Add(0, lr)
		_, lrEnemy := enemies[lrSquare]
		if lrEnemy && IsEPPawn(b.GetPiece(lrSquare)) {
			lrSquare = lrSquare.Add(heading, 0)
			out = append(out, move{a: &c, b: &lrSquare, special: EnPassant})
		}
	}
	//each diagonal attack
	for _, lr := range [2]int{-1, 1} {
		diagonalSquare := c.Copy().Add(heading, lr)
		_, hasEnemy := enemies[diagonalSquare]
		if hasEnemy {
			out = append(out, move{a: &c, b: &diagonalSquare})
		}
	}
	return out
}
