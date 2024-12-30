package board

import (
	"fmt"
	"os"

	"github.com/tawatts1/go_chess/utility"
)

type board struct {
	grid [8][8]rune
}

// coord, y is the column, x is the row of the board.
type coord struct {
	y, x int
}

func (c coord) String() string {
	return fmt.Sprintf("(y:%v,x:%v)", c.y, c.x)
}

const Black, White rune = 'b', 'w'
const Space rune = '0'
const BlackPawn, WhitePawn rune = 'p', 'P'
const BlackKnight, WhiteKnight rune = 'n', 'N'
const BlackBishop, WhiteBishop rune = 'b', 'B'
const BlackRookNC, WhiteRookNC rune = 'r', 'R'
const BlackQueen, WhiteQueen rune = 'q', 'Q'
const BlackKing, WhiteKing rune = 'k', 'K'
const BlackRookC, WhiteRookC rune = 'o', 'O'
const BlackPawnEP, WhitePawnEP rune = 'a', 'A'

// A map that maps the encoded board runes to human readable strings
var humanMap = map[rune]string{
	Space:       "  ",
	BlackPawn:   "bp",
	BlackKnight: "bn",
	BlackBishop: "bb",
	BlackRookNC: "br",
	BlackQueen:  "bq",
	BlackKing:   "bk",
	BlackRookC:  "BR", // a rook with the ability to castle
	BlackPawnEP: "BP", // a pawn which can be taken via en passant
	//--------
	WhitePawn:   "wp",
	WhiteKnight: "wn",
	WhiteBishop: "wb",
	WhiteRookNC: "wr",
	WhiteQueen:  "wq",
	WhiteKing:   "wk",
	WhiteRookC:  "WR", // a rook with the ability to castle
	WhitePawnEP: "WP", // a pawn which can be taken via en passant
}

var blackMap = map[rune]bool{
	BlackPawn:   true,
	BlackKnight: true,
	BlackBishop: true,
	BlackRookNC: true,
	BlackQueen:  true,
	BlackKing:   true,
	BlackRookC:  true,
	BlackPawnEP: true,
}

var whiteMap = map[rune]bool{
	WhitePawn:   true,
	WhiteKnight: true,
	WhiteBishop: true,
	WhiteRookNC: true,
	WhiteQueen:  true,
	WhiteKing:   true,
	WhiteRookC:  true,
	WhitePawnEP: true,
}

const StartingBoard string = "onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO"

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
