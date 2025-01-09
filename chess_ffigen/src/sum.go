// sum.go file
package main

import "C"
import (
	"math"
	"time"

	"github.com/tawatts1/go_chess/board"
)

//export GetBoardAfterMove
func GetBoardAfterMove(boardStr *C.char, y1, x1, y2, x2 C.int) *C.char {
	return C.CString(GetBoardAfterMoveEncoded(C.GoString(boardStr), int(y1), int(x1), int(y2), int(x2)))
}

//export GetNextMoves
func GetNextMoves(boardStr *C.char, y C.int, x C.int) *C.char {
	return C.CString(GetMoveDestinationsEncoded(C.GoString(boardStr), int(y), int(x)))
}

func GetBoardAfterMoveEncoded(boardStr string, y1, x1, y2, x2 int) string {
	b := board.GetBoardFromString(boardStr)
	c1 := board.NewCoord(y1, x1)
	c2 := board.NewCoord(y2, x2)
	moves := board.GetMovesFromBoardCoord(b, c1)
	m := board.GetFirstEqualMove(moves, c1, c2)
	b2 := board.GetBoardAfterMove(b, m)
	return b2.Encode()
}

func GetMoveDestinationsEncoded(boardStr string, y, x int) string {
	b := board.GetBoardFromString(boardStr)
	c := board.NewCoord(y, x)
	var friends, enemies map[board.Coord]bool
	if board.IsWhite(b.GetPiece(c)) {
		friends = b.GetWhiteCoordMap()
		enemies = b.GetBlackCoordMap()
	} else {
		friends = b.GetBlackCoordMap()
		enemies = b.GetWhiteCoordMap()
	}
	moves := b.GetMoves(friends, enemies, c, true)
	out := ""
	for _, m := range moves {
		out += m.Encode()
	}
	return out
}

//export sum
func sum(a C.int, b C.int) C.int {
	return C.int(math.Floor(math.Log(float64(a + b))))
	//return a + b
}

//export longSum
func longSum(a C.int, b C.int) C.int {
	time.Sleep(4 * time.Second)
	return a + b + 100
}

//export enforce_binding
func enforce_binding() {}

func main() {}
