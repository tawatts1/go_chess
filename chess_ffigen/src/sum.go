// sum.go file
package main

import "C"
import (
	"math"
	"time"

	"github.com/tawatts1/go_chess/ai"
	"github.com/tawatts1/go_chess/board"
)

//export GetAiChosenMove
func GetAiChosenMove(boardStr *C.char, isWhite C.int, aiName *C.char, N C.int) *C.char {
	whiteBool := int(isWhite) == 1
	return C.CString(getAiChosenMove(C.GoString(boardStr), whiteBool, C.GoString(aiName), int(N)))
}

func getAiChosenMove(boardStr string, isWhite bool, aiName string, N int) string {
	a := ai.GetAiFromString(aiName)
	b := board.GetBoardFromString(boardStr)
	m := a.ChooseMove(b, isWhite, N)
	time.Sleep(4 * time.Second)
	return m.Encode()
}

//export GetBoardAfterMove
func GetBoardAfterMove(boardStr *C.char, y1, x1, y2, x2 C.int) *C.char {
	return C.CString(getBoardAfterMoveEncoded(C.GoString(boardStr), int(y1), int(x1), int(y2), int(x2)))
}

func getBoardAfterMoveEncoded(boardStr string, y1, x1, y2, x2 int) string {
	b := board.GetBoardFromString(boardStr)
	c1 := board.NewCoord(y1, x1)
	c2 := board.NewCoord(y2, x2)
	moves := board.GetMovesFromBoardCoord(b, c1)
	m := board.GetFirstEqualMove(moves, c1, c2)
	b2 := board.GetBoardAfterMove(b, m)
	return b2.Encode()
}

//export GetNextMoves
func GetNextMoves(boardStr *C.char, y C.int, x C.int) *C.char {
	return C.CString(getMoveDestinationsEncoded(C.GoString(boardStr), int(y), int(x)))
}

func getMoveDestinationsEncoded(boardStr string, y, x int) string {
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
		out += m.EncodeB() + "|"
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
