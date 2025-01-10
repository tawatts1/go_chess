package board

import "fmt"

// A way to encode a Move via a Move from one coordinate to another
type Move struct {
	a, b    Coord
	special rune
}

func NewMove(c1, c2 Coord, special rune) Move {
	return Move{a: c1, b: c2, special: special}
}

func (m1 Move) Equals(m2 Move) bool {
	return m1.a.Equals(m2.a) && m1.b.Equals(m2.b) && m1.special == m2.special
}

func Contains(moves []Move, m Move) bool {
	for _, m2 := range moves {
		if m2.Equals(m) {
			return true
		}
	}
	return false
}

func (m Move) String() string {
	return fmt.Sprintf("{%v, %v, %c }", m.a, m.b, m.special)
}

func (m Move) Encode() string {
	return fmt.Sprintf("%v,%v,%v,%v", m.a.y, m.a.x, m.b.y, m.b.x)
}

func (m Move) EncodeB() string {
	return m.b.Encode()
}

func GetFirstEqualMove(moves []Move, c1, c2 Coord) Move {
	for _, m := range moves {
		if m.a.Equals(c1) && m.b.Equals(c2) {
			return m
		}
	}
	panic("Move match not found")
}

func GetMovesFromBoardCoord(b Board, c Coord) []Move {
	var friends, enemies map[Coord]bool
	if IsWhite(b.GetPiece(c)) {
		friends = b.GetWhiteCoordMap()
		enemies = b.GetBlackCoordMap()
	} else {
		friends = b.GetBlackCoordMap()
		enemies = b.GetWhiteCoordMap()
	}
	return b.GetMoves(friends, enemies, c, true)
}

func (b Board) GetLegalMoves(isWhite bool) []Move {

	var friends, enemies map[Coord]bool

	if isWhite {
		friends = b.GetWhiteCoordMap()
		enemies = b.GetBlackCoordMap()
	} else {
		friends = b.GetBlackCoordMap()
		enemies = b.GetWhiteCoordMap()
	}
	out := make([]Move, 0, 2*len(friends))
	for _, c := range AllCoordinates {
		_, hasCoord := friends[c]
		if hasCoord {
			out = append(out, b.GetMoves(friends, enemies, c, true)...)
		}
	}
	return out
}

func (b Board) GetMoves(friends, enemies map[Coord]bool, c Coord, filterIllegalMoves bool) []Move {
	var out []Move
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
		castleMoves := make([]Move, 0, 1)
		if IsWhite(piece) && c.Equals(Coord{y: 7, x: 4}) {
			for _, rook_coord := range b.GetCastleableRooks(friends) {
				if rook_coord.Equals(BottomLeft) &&
					b.IsLocEmpty(7, 1) &&
					b.IsLocEmpty(7, 2) &&
					b.IsLocEmpty(7, 3) &&
					(!filterIllegalMoves || Contains(out, Move{a: c, b: Coord{y: 7, x: 3}, special: WhiteKing})) {
					castleMoves = append(castleMoves, Move{a: c, b: Coord{y: 7, x: 2}, special: CastleBridge})
				} else if rook_coord.Equals(BottomRight) &&
					b.IsLocEmpty(7, 5) &&
					b.IsLocEmpty(7, 6) &&
					(!filterIllegalMoves || Contains(out, Move{a: c, b: Coord{y: 7, x: 5}, special: WhiteKing})) {
					castleMoves = append(castleMoves, Move{a: c, b: Coord{y: 7, x: 6}, special: CastleBridge})
				}
			}
		} else if IsBlack(piece) && c.Equals(Coord{y: 0, x: 4}) {
			for _, rook_coord := range b.GetCastleableRooks(friends) {
				if rook_coord.Equals(TopLeft) &&
					b.IsLocEmpty(0, 1) &&
					b.IsLocEmpty(0, 2) &&
					b.IsLocEmpty(0, 3) &&
					(!filterIllegalMoves || Contains(out, Move{a: c, b: Coord{y: 0, x: 3}, special: BlackKing})) {
					castleMoves = append(castleMoves, Move{a: c, b: Coord{y: 0, x: 2}, special: CastleBridge})
				} else if rook_coord.Equals(TopRight) &&
					b.IsLocEmpty(0, 5) &&
					b.IsLocEmpty(0, 6) &&
					(!filterIllegalMoves || Contains(out, Move{a: c, b: Coord{y: 0, x: 5}, special: BlackKing})) {
					castleMoves = append(castleMoves, Move{a: c, b: Coord{y: 0, x: 6}, special: CastleBridge})
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

func FilterIllegalMoves(b Board, moveSlice []Move, piece rune) []Move {
	filteredOut := make([]Move, 0, len(moveSlice))
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

func (b Board) GetPawnMoves(friends, enemies map[Coord]bool, c Coord, heading int) []Move {
	out := make([]Move, 0, 2)
	piece := b.GetPiece(c)
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
			out = append(out, Move{a: c, b: forward})
		}

		// two steps ahead of pawn, for opening
		if (!isBlack && c.y == 6) ||
			(isBlack && c.y == 1) {
			forward2 := c.Copy().Add(heading*2, 0)
			_, hasFriend2 := friends[forward2]
			_, hasEnemy2 := enemies[forward2]
			if !hasFriend2 && !hasEnemy2 {
				out = append(out, Move{a: c, b: forward2, special: EnPassantMap[piece]})
			}
		}
	}
	//each En Passant
	for _, lr := range [2]int{-1, 1} {
		lrSquare := c.Copy().Add(0, lr)
		_, lrEnemy := enemies[lrSquare]
		if lrEnemy && IsEPPawn(b.GetPiece(lrSquare)) {
			lrSquare = lrSquare.Add(heading, 0)
			out = append(out, Move{a: c, b: lrSquare, special: EnPassant})
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
				out = append(out, Move{a: c, b: diagonalSquare})
			}
		}
	}
	return out
}

func getPawnPromotionMoves(from, to Coord, promotionCodes []rune) []Move {
	//when promoted, pawns can become knight, bishop, rook or queen.
	numPromotions := len(promotionCodes)
	var out = make([]Move, numPromotions, numPromotions)
	for i := range numPromotions {
		out[i] = Move{a: from, b: to, special: (promotionCodes)[i]}
	}
	return out
}

func GetBishopMoves(friends, enemies map[Coord]bool, c Coord) []Move {
	out := make([]Move, 0, 3)
	var newSquare Coord
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
				out = append(out, Move{a: c, b: newSquare})
				_, hasEnemy := enemies[newSquare]
				if hasEnemy {
					break // with the other vectors
				}
			}
		}
	}
	return out
}

func (b Board) GetRookMoves(friends, enemies map[Coord]bool, c Coord, isCastleable bool) []Move {
	out := make([]Move, 0, 3)
	piece := b.GetPiece(c)
	var newSquare Coord
	// vectors used to add diagonal moves to the starting square
	vectors := [4]Coord{{y: 1, x: 0}, {y: -1, x: 0}, {y: 0, x: 1}, {y: 0, x: -1}}
	for _, v := range vectors {
		for newSquare = c.Copy().Add(v.y, v.x); newSquare.IsInBoard(); newSquare = newSquare.Copy().Add(v.y, v.x) {
			_, hasFriend := friends[newSquare]
			if hasFriend {
				break // with the other vectors
			}
			//add move whether there is an enemy there or it is empty
			if isCastleable {
				out = append(out, Move{a: c, b: newSquare, special: piece})
			} else {
				out = append(out, Move{a: c, b: newSquare})
			}
			_, hasEnemy := enemies[newSquare]
			if hasEnemy {
				break // with the other vectors
			}
		}
	}
	return out
}

func (b Board) GetQueenMoves(friends, enemies map[Coord]bool, c Coord) []Move {
	return append(GetBishopMoves(friends, enemies, c), b.GetRookMoves(friends, enemies, c, false)...)
}

func GetKnightMoves(friends, enemies map[Coord]bool, c Coord) []Move {
	out := make([]Move, 0, 4)
	var newSquare Coord
	for _, sign := range [2]int{-1, 1} {
		for _, vector := range [4]Coord{{y: 1, x: 2}, {y: 2, x: 1}, {y: -1, x: 2}, {y: -2, x: 1}} {
			newSquare = c.Add(sign*vector.y, sign*vector.x)
			if newSquare.IsInBoard() {
				_, hasFriend := friends[newSquare]
				if !hasFriend {
					out = append(out, Move{a: c, b: newSquare})
				}
			}
		}
	}
	return out
}

func (b Board) GetKingMoves(friends, enemies map[Coord]bool, c Coord) []Move {
	out := make([]Move, 0, 2)
	piece := b.GetPiece(c)
	var newSquare Coord
	for _, sign := range [2]int{-1, 1} {
		for _, vector := range [4]Coord{{y: 1, x: 1}, {y: 1, x: 0}, {y: 0, x: 1}, {y: -1, x: 1}} {
			newSquare = c.Add(sign*vector.y, sign*vector.x)
			if newSquare.IsInBoard() {
				_, hasFriend := friends[newSquare]
				if !hasFriend {
					if c.Equals(BlackKingHome) || c.Equals(WhiteKingHome) {
						out = append(out, Move{a: c, b: newSquare, special: piece})
					} else {
						out = append(out, Move{a: c, b: newSquare})
					}
				}
			}
		}
	}
	return out
}

func (b Board) IsInCheck(friends, enemies map[Coord]bool, king_coord Coord) bool {
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

func GetBoardAfterMove(b Board, m Move) Board {
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
	if m.special == BlackPawnEP || m.special == WhitePawnEP {
		a, b := m.a, m.b
		out.grid[b.y][b.x] = m.special
		out.grid[a.y][a.x] = Space
		return out
	}

	panic("Not implemented - special move")
}
