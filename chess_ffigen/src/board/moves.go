package board

// A way to encode a move via a move from one coordinate to another
type move struct {
	a, b    *coord
	special rune
}

const EnPassant rune = 'e'
const PawnPromotion rune = 'p'
const CastleBridge rune = 'c'

func (b board) GetBlackMoves() []move {
	out := make([]move, 0, 8)
	bcm := b.GetBlackCoordMap()
	wcm := b.GetWhiteCoordMap()
	for c, _ := range bcm { // black piece coordinate
		bp := b.GetPiece(c)
		if IsPawn(bp) {
			out = append(out, b.GetPawnMoves(bcm, wcm, c, 1)...)
		}
	}
	return out

}

func (b board) GetPawnMoves(friends, enemies map[coord]bool, c coord, heading int) []move {
	out := make([]move, 0, 8)
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
				out = append(out, move{a: &c, b: &forward2, special: EnPassant})
			}
		}
	}
	// each diagonal attack
	for _, lr := range [2]int{-1, 1} {
		diagonalSquare := c.Copy().Add(heading, lr)
		_, hasEnemy := enemies[diagonalSquare]
		if hasEnemy {
			out = append(out, move{a: &c, b: &diagonalSquare})
		}
	}

	return out
}
