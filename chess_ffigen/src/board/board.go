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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetBoardFromFile(fname string) *board {
	data, err := os.ReadFile(fname)
	check(err)
	return GetBoardFromString(string(data))
}

func GetBoardFromString(str string) *board {
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
	return &b
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

func IsBlack(piece rune) bool {
	blk, ok := blackMap[piece]
	return ok && blk
}

func IsWhite(piece rune) bool {
	wht, ok := whiteMap[piece]
	return ok && wht
}

func (b board) GetWhiteCoords() []coord {
	out := make([]coord, 0, 8)
	for y := range len(b.grid) {
		for x := range len(b.grid[0]) {
			if IsWhite(b.grid[y][x]) {
				out = append(out, coord{y: y, x: x})
			}
		}
	}
	return out
}

func (b board) GetBlackCoords() []coord {
	out := make([]coord, 0, 8)
	for y := range len(b.grid) {
		for x := range len(b.grid[0]) {
			if IsBlack(b.grid[y][x]) {
				out = append(out, coord{y: y, x: x})
			}
		}
	}
	return out
}

func (b board) GetBlackCoordMap() map[coord]bool {
	out := make(map[coord]bool)
	for _, c := range b.GetBlackCoords() {
		out[c] = true
	}
	return out
}

func (b board) GetWhiteCoordMap() map[coord]bool {
	out := make(map[coord]bool)
	for _, c := range b.GetWhiteCoords() {
		out[c] = true
	}
	return out
}

func (b board) GetPiece(c coord) rune {
	return b.grid[c.y][c.x]
}

func (b board) GetKingCoord(friends map[coord]bool) coord {
	for c := range friends {
		if IsKing(b.GetPiece(c)) {
			return c
		}
	}
	panic("King not found!")
}
