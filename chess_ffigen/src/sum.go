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
	return C.CString(board.GetBoardAfterMoveEncoded(C.GoString(boardStr), int(y1), int(x1), int(y2), int(x2)))
}

//export GetNextMoves
func GetNextMoves(boardStr *C.char, y C.int, x C.int) *C.char {
	return C.CString(board.GetMoveDestinationsEncoded(C.GoString(boardStr), int(y), int(x)))
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
