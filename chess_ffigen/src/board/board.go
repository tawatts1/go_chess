package board

import (
	"fmt"
	"os"

	"github.com/tawatts1/go_chess/utility"
)

const StartingBoard string = "onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO"
const BoardHeight, BoardWidth = 8, 8

type board struct {
	grid [BoardHeight][BoardWidth]rune
}

func (b board) Copy() board {
	return board{grid: b.grid}
}

func (b1 board) Equals(b2 board) bool {
	for i := range BoardHeight {
		for j := range BoardWidth {
			if b1.grid[i][j] != b2.grid[i][j] {
				return false
			}
		}
	}
	return true
}

func (b board) SimpleMove(c1, c2 Coord) board {
	b.grid[c2.y][c2.x] = b.grid[c1.y][c1.x]
	b.grid[c1.y][c1.x] = Space
	return b
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetBoardFromFile(fname string) board {
	data, err := os.ReadFile(fname)
	check(err)
	return GetBoardFromString(string(data))
}

func GetBoardFromString(str string) board {
	s := utility.RemoveWhitespace(str)
	b := board{}
	l, w := len(b.grid), len(b.grid[0])
	numSquares := l * w
	for i, r := range s {
		if i >= numSquares {
			panic(fmt.Sprintf("Too many runes in input string: %s\ncleaned string: %s", str, s))
		}
		b.grid[i/w][i%w] = r
	}
	return b
}

func (b board) String() string {
	out := ""
	divider := "\n+--+--+--+--+--+--+--+--+\n"
	out += divider
	for _, row := range b.grid {
		out += "|"
		for _, spot := range row {
			out += humanMap[spot] + "|"
		}
		out += divider
	}
	return out
}

func (b board) Encode() string {
	out := ""
	for i := range BoardHeight {
		for j := range BoardWidth {
			out += string(b.grid[i][j])
		}
	}
	return out
}

func IsBlack(piece rune) bool {
	blk, ok := blackMap[piece]
	return ok && blk
}

func IsWhite(piece rune) bool {
	wht, ok := whiteMap[piece]
	return ok && wht
}

func GetColor(piece rune) rune {
	if IsWhite(piece) {
		return White
	} else {
		return Black
	}
}

func (b board) GetWhiteCoords() []Coord {
	out := make([]Coord, 0, 8)
	for y := range len(b.grid) {
		for x := range len(b.grid[0]) {
			if IsWhite(b.grid[y][x]) {
				out = append(out, Coord{y: y, x: x})
			}
		}
	}
	return out
}

func (b board) GetBlackCoords() []Coord {
	out := make([]Coord, 0, 8)
	for y := range len(b.grid) {
		for x := range len(b.grid[0]) {
			if IsBlack(b.grid[y][x]) {
				out = append(out, Coord{y: y, x: x})
			}
		}
	}
	return out
}

func (b board) GetBlackCoordMap() map[Coord]bool {
	out := make(map[Coord]bool)
	for _, c := range b.GetBlackCoords() {
		out[c] = true
	}
	return out
}

func (b board) GetWhiteCoordMap() map[Coord]bool {
	out := make(map[Coord]bool)
	for _, c := range b.GetWhiteCoords() {
		out[c] = true
	}
	return out
}

func (b board) GetPiece(c Coord) rune {
	return b.grid[c.y][c.x]
}

func (b board) GetKingCoord(friends map[Coord]bool) Coord {
	for c := range friends {
		if IsKing(b.GetPiece(c)) {
			return c
		}
	}
	panic("King not found!")
}

func (b board) GetCastleableRooks(friends map[Coord]bool) []Coord {
	out := make([]Coord, 0, 2)
	for c := range friends {
		if IsRookCastleable(b.GetPiece(c)) {
			out = append(out, c)
		}
	}
	return out
}

func (b board) IsCoordEmpty(c Coord) bool {
	return b.IsLocEmpty(c.y, c.x)
}

func (b board) IsLocEmpty(y, x int) bool {
	return b.grid[y][x] == Space
}
