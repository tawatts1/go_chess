package board

import (
	"fmt"
	"os"

	"github.com/tawatts1/go_chess/utility"
)

type board struct {
	grid [8][8]rune
}

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
var human_map = map[rune]string{
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
			out += human_map[spot] + "|"
		}
		out += divider
	}
	return out
}
