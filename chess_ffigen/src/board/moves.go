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

var BlackPawnPromotion = []rune{BlackKnight, BlackBishop, BlackRookNC, BlackQueen}
var WhitePawnPromotion = []rune{WhiteKnight, WhiteBishop, WhiteRookNC, WhiteQueen}

func (b board) GetMoves(friends, enemies map[coord]bool, c coord, filterIllegalMoves bool) []move {
	var out []move
	piece := b.GetPiece(c)
	if IsPawn(piece) {
		heading := -1
		if IsBlack(piece) {
			heading = 1
		}
		out = b.GetPawnMoves(friends, enemies, c, heading)
	} else if IsBishop(piece) {
		out = GetBishopMoves(friends, enemies, c)
	} else if IsRook(piece) {
		if IsRookCastleable(piece) {
			out = b.GetRookMoves(friends, enemies, c, true)
		} else {
			out = b.GetRookMoves(friends, enemies, c, false)
		}

	} else if IsQueen(piece) {
		out = b.GetQueenMoves(friends, enemies, c)
	} else if IsKnight(piece) {
		out = GetKnightMoves(friends, enemies, c)
	} else if IsKing(piece) {
		out = b.GetKingMoves(friends, enemies, c)
	} else {
		panic("Not implemented")
	}
	if filterIllegalMoves {
		out = FilterIllegalMoves(b, out, piece)
	}
	//Now add in castle/bridge moves manually
	if IsKing(piece) && (!filterIllegalMoves || !b.IsInCheck(friends, enemies, b.GetKingCoord(friends))) {
		castleMoves := make([]move, 0, 1)
		if IsWhite(piece) && c.Equals(coord{y: 7, x: 4}) {
			for _, rook_coord := range b.GetCastleableRooks(friends) {
				if rook_coord.Equals(BottomLeft) &&
					b.IsLocEmpty(7, 1) &&
					b.IsLocEmpty(7, 2) &&
					b.IsLocEmpty(7, 3) &&
					(!filterIllegalMoves || AnyEqual(out, move{a: c, b: coord{y: 7, x: 3}, special: WhiteKing})) {
					castleMoves = append(castleMoves, move{a: c, b: coord{y: 7, x: 2}, special: CastleBridge})
				} else if rook_coord.Equals(BottomRight) &&
					b.IsLocEmpty(7, 5) &&
					b.IsLocEmpty(7, 6) &&
					(!filterIllegalMoves || AnyEqual(out, move{a: c, b: coord{y: 7, x: 5}, special: WhiteKing})) {
					castleMoves = append(castleMoves, move{a: c, b: coord{y: 7, x: 6}, special: CastleBridge})
				}
			}
		} else if IsBlack(piece) && c.Equals(coord{y: 0, x: 4}) {
			for _, rook_coord := range b.GetCastleableRooks(friends) {
				if rook_coord.Equals(TopLeft) &&
					b.IsLocEmpty(0, 1) &&
					b.IsLocEmpty(0, 2) &&
					b.IsLocEmpty(0, 3) &&
					(!filterIllegalMoves || AnyEqual(out, move{a: c, b: coord{y: 0, x: 3}, special: BlackKing})) {
					castleMoves = append(castleMoves, move{a: c, b: coord{y: 0, x: 2}, special: CastleBridge})
				} else if rook_coord.Equals(TopRight) &&
					b.IsLocEmpty(0, 5) &&
					b.IsLocEmpty(0, 6) &&
					(!filterIllegalMoves || AnyEqual(out, move{a: c, b: coord{y: 0, x: 5}, special: BlackKing})) {
					castleMoves = append(castleMoves, move{a: c, b: coord{y: 0, x: 6}, special: CastleBridge})
				}
			}
		}
		if filterIllegalMoves {
			castleMoves = FilterIllegalMoves(b, castleMoves, piece)
		}
		out = append(out, castleMoves...)
	}
	return out
}

func FilterIllegalMoves(b board, moveSlice []move, piece rune) []move {
	filteredOut := make([]move, 0, len(moveSlice))
	for i := 0; i < len(moveSlice); i++ {
		m := moveSlice[i]
		newBoard := GetBoardAfterMove(b, m)
		wcm := newBoard.GetWhiteCoordMap()
		bcm := newBoard.GetBlackCoordMap()
		if IsWhite(piece) {
			if !newBoard.IsInCheck(wcm, bcm, newBoard.GetKingCoord(wcm)) {
				filteredOut = append(filteredOut, m)
			}
		} else {
			// black piece moving
			if !newBoard.IsInCheck(bcm, wcm, newBoard.GetKingCoord(bcm)) {
				filteredOut = append(filteredOut, m)
			}
		}
	}
	return filteredOut
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
			var promotions []rune
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
				var promotions []rune
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

func getPawnPromotionMoves(from, to coord, promotionCodes []rune) []move {
	//when promoted, pawns can become knight, bishop, rook or queen.
	numPromotions := len(promotionCodes)
	var out = make([]move, numPromotions, numPromotions)
	for i := range numPromotions {
		out[i] = move{a: from, b: to, special: (promotionCodes)[i]}
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

func (b board) GetRookMoves(friends, enemies map[coord]bool, c coord, isCastleable bool) []move {
	out := make([]move, 0, 3)
	piece := b.GetPiece(c)
	var newSquare coord
	// vectors used to add diagonal moves to the starting square
	vectors := [4]coord{{y: 1, x: 0}, {y: -1, x: 0}, {y: 0, x: 1}, {y: 0, x: -1}}
	for _, v := range vectors {
		for newSquare = c.Copy().Add(v.y, v.x); newSquare.IsInBoard(); newSquare = newSquare.Copy().Add(v.y, v.x) {
			_, hasFriend := friends[newSquare]
			if hasFriend {
				break // with the other vectors
			}
			//add move whether there is an enemy there or it is empty
			if isCastleable {
				out = append(out, move{a: c, b: newSquare, special: piece})
			} else {
				out = append(out, move{a: c, b: newSquare})
			}
			_, hasEnemy := enemies[newSquare]
			if hasEnemy {
				break // with the other vectors
			}
		}
	}
	return out
}

func (b board) GetQueenMoves(friends, enemies map[coord]bool, c coord) []move {
	return append(GetBishopMoves(friends, enemies, c), b.GetRookMoves(friends, enemies, c, false)...)
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

func (b board) GetKingMoves(friends, enemies map[coord]bool, c coord) []move {
	out := make([]move, 0, 2)
	piece := b.GetPiece(c)
	var newSquare coord
	for _, sign := range [2]int{-1, 1} {
		for _, vector := range [4]coord{{y: 1, x: 1}, {y: 1, x: 0}, {y: 0, x: 1}, {y: -1, x: 1}} {
			newSquare = c.Add(sign*vector.y, sign*vector.x)
			if newSquare.IsInBoard() {
				_, hasFriend := friends[newSquare]
				if !hasFriend {
					if c.Equals(BlackKingHome) || c.Equals(WhiteKingHome) {
						out = append(out, move{a: c, b: newSquare, special: piece})
					} else {
						out = append(out, move{a: c, b: newSquare})
					}
				}
			}
		}
	}
	return out
}

func (b board) IsInCheck(friends, enemies map[coord]bool, king_coord coord) bool {
	for enemy_coord := range enemies {
		// note that we are checking what moves the enemy has, so the friends and enemies maps are switched.
		for _, m := range b.GetMoves(enemies, friends, enemy_coord, false) {
			if king_coord.Equals(m.b) {
				return true
			}
		}
	}
	return false
}

func GetBoardAfterMove(b board, m move) board {
	out := b.Copy()
	color := GetColor(b.GetPiece(m.a))
	epPawnPiece := EnPassantMap[color] // note that this piece is of the oposite color, according to the map.
	for i := range BoardHeight {
		for j := range BoardWidth {
			if out.grid[i][j] == epPawnPiece {
				out.grid[i][j] = EnPassantMap[epPawnPiece]
			}
		}
	}

	if m.special == 0 {
		return out.SimpleMove(m.a, m.b)
	} else if m.special == EnPassant {
		a, b := m.a, m.b
		out.grid[b.y][b.x] = out.grid[a.y][a.x]
		out.grid[a.y][b.x] = Space
		out.grid[a.y][a.x] = Space
		return out
	} else if m.special == CastleBridge {
		a, b := m.a, m.b
		out = out.SimpleMove(a, b)
		if b.x == 2 {
			rook_destination := b.Add(0, 1)
			out = out.SimpleMove(b.Add(0, -2), rook_destination)
		} else if b.x == 6 {
			rook_destination := b.Add(0, -1)
			out = out.SimpleMove(b.Add(0, 1), rook_destination)
			castlingRookCode := out.grid[rook_destination.y][rook_destination.x]
			for i := range BoardHeight {
				for j := range BoardWidth {
					if out.grid[i][j] == castlingRookCode {
						out.grid[i][j] = CastleMap[castlingRookCode]
					}
				}
			}
		} else {
			panic("Bad bridge move")
		}
		return out
	} else if m.special == BlackKing || m.special == WhiteKing {
		//black/white king was moving from home position, check for castleable rooks and make then not castleable.
		out = out.SimpleMove(m.a, m.b)
		castleableRookCode := CastleMap[m.special]
		for i := range BoardHeight {
			for j := range BoardWidth {
				if out.grid[i][j] == castleableRookCode {
					out.grid[i][i] = CastleMap[castleableRookCode]
				}
			}
		}
		return out
	} else if m.special == BlackRookC || m.special == WhiteRookC {
		a, b := m.a, m.b
		out.grid[b.y][b.x] = CastleMap[m.special]
		out.grid[a.y][a.x] = Space
		return out
	}
	for _, promotion := range append(BlackPawnPromotion, WhitePawnPromotion...) {
		if m.special == promotion {
			a, b := m.a, m.b
			out.grid[b.y][b.x] = promotion
			out.grid[a.y][a.x] = Space
			return out
		}
	}

	panic("Not implemented - special move")
}
