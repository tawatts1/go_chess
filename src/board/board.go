package board

import (
	"fmt"
	"os"

	"github.com/tawatts1/go_chess/utility"
)

type board struct {
	grid [8][8]rune
}

// A map that maps the encoded board runes to human readable strings
var human_map = map[rune]string{
	'0': "  ",
	'p': "bp",
	'n': "bn",
	'b': "bb",
	'r': "br",
	'q': "bq",
	'k': "bk",
	'o': "BR", // a rook with the ability to castle
	'a': "BP", // a pawn which can be taken via en passant
	//--------
	'P': "wp",
	'N': "wn",
	'B': "wb",
	'R': "wr",
	'Q': "wq",
	'K': "wk",
	'O': "WR", // a rook with the ability to castle
	'A': "WP", // a pawn which can be taken via en passant
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
