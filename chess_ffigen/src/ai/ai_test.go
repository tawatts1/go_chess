package ai

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/tawatts1/go_chess/board"
)

var verbose bool = true
var testFolder = "testingMoves/"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hasError(e error) bool {
	return e != nil
}

func testAiMoveFile(fname string) string {
	data, err := os.ReadFile(fname)
	check(err)
	lines := strings.Split(string(data), "\n")
	//var a ai
	var N int
	var b board.Board
	for lineIndex, line := range lines {
		if verbose {
			fmt.Println(line)
		}
		if len(line) < 5 || []rune(line)[0] == '#' {
			continue
		}
		args := strings.Split(line, ",")
		if args[0] == "new" {
			//a = GetAiFromString(args[1])
			Nparsed, ok := strconv.Atoi(args[2])
			b = board.GetBoardFromString(args[3])
			if hasError(ok) {
				panic("failed to parse N")
			} else {
				N = Nparsed
			}
		} else if args[0] == "move" || args[0] == "notmove" {
			color := args[1]
			y1, ok1 := strconv.Atoi(strings.TrimSpace(args[2]))
			x1, ok2 := strconv.Atoi(strings.TrimSpace(args[3]))
			y2, ok3 := strconv.Atoi(strings.TrimSpace(args[4]))
			x2, ok4 := strconv.Atoi(strings.TrimSpace(args[5]))
			var isWhite bool
			if color == string(board.White) {
				isWhite = true
			} else if color == string(board.Black) {
				isWhite = false
			} else {
				panic("Color not recognized")
			}
			var special rune
			var c1, c2 board.Coord
			if hasError(ok1) || hasError(ok2) || hasError(ok3) || hasError(ok4) {
				panic(fmt.Sprintf("failed to parse line %v", lineIndex+1))
			} else {
				c1 = board.NewCoord(y1, x1)
				c2 = board.NewCoord(y2, x2)
			}

			r := []rune(strings.TrimSpace(args[6]))
			if len(r) > 1 {
				return fmt.Sprintf("%v: special column must be one or zero characters. ", lineIndex+1)
			} else if len(r) == 1 {
				special = r[0]
			}
			mExpected := board.NewMove(c1, c2, special)
			mResult := ChooseMove(b, isWhite, N)
			if args[0] == "move" && !mResult.Equals(mExpected) {
				return fmt.Sprintf("line %v: Expected %v but got %v", lineIndex+1, mExpected, mResult)
			} else if args[0] == "notmove" && mResult.Equals(mExpected) {
				return fmt.Sprintf("line %v: Expected anything but %v", lineIndex+1, mResult)
			}
		}
	}
	return ""
}

func TestChooseMove(t *testing.T) {
	err := testAiMoveFile(testFolder + "aiTests.txt")
	if err != "" {
		t.Error(err)
	}
}
