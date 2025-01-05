package board

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

var boardsDir string = "testingBoards/"

func hasError(e error) bool {
	if e != nil {
		return true
	}
	return false
}

func testMoveFile(fname string) string {
	data, err := os.ReadFile(fname)
	check(err)
	lines := strings.Split(string(data), "\n")
	var b *board
	var bcm, wcm map[coord]bool
	for lineIndex, line := range lines {
		//fmt.Println(line)
		if len(line) < 5 || []rune(line)[0] == '#' {
			continue
		}
		args := strings.Split(line, ",")
		if args[0] == "new" {
			b = GetBoardFromString(args[1])
			bcm = b.GetBlackCoordMap()
			wcm = b.GetWhiteCoordMap()
			//fmt.Println(b)
		} else if args[0] == "num" {
			y_, ok1 := strconv.Atoi(args[1])
			x_, ok2 := strconv.Atoi(args[2])
			expected_len, ok3 := strconv.Atoi(args[3])
			var c coord
			if hasError(ok1) || hasError(ok2) || hasError(ok3) {
				panic(fmt.Sprintf("failed to parse line %v", lineIndex+1))
			} else {
				c = coord{y: y_, x: x_}
			}
			actual_len := len(b.GetMoves(c, bcm, wcm))
			if actual_len != expected_len {
				return fmt.Sprintf("%v: expected %v, got %v", lineIndex+1, expected_len, actual_len)
			}
		} else if args[0] == "has" {
			y1, ok1 := strconv.Atoi(args[1])
			x1, ok2 := strconv.Atoi(args[2])
			y2, ok3 := strconv.Atoi(args[3])
			x2, ok4 := strconv.Atoi(args[4])
			var special rune
			var c1, c2 coord
			var m move
			if hasError(ok1) || hasError(ok2) || hasError(ok3) || hasError(ok4) {
				panic(fmt.Sprintf("failed to parse line %v", lineIndex+1))
			} else {
				c1 = coord{y: y1, x: x1}
				c2 = coord{y: y2, x: x2}
			}
			if len(args) >= 6 {
				r := []rune(args[5])
				if len(r) > 1 {
					return fmt.Sprintf("%v: special column must be one or zero characters. ", lineIndex+1)
				} else if len(r) == 1 {
					special = r[0]
				}
			}
			m = move{a: c1, b: c2, special: special}
			moves := b.GetMoves(c1, bcm, wcm)
			if !AnyEqual(moves, m) {
				return fmt.Sprintf("%v: expected %v to contain %v", lineIndex+1, moves, m)
			}
		}
	}
	return ""
}

func TestGetPawnMoves(t *testing.T) {
	err := testMoveFile(boardsDir + "pawnTests.txt")
	if err != "" {
		t.Error(err)
	}
}

func TestGetBishopMoves(t *testing.T) {
	err := testMoveFile(boardsDir + "bishopTests.txt")
	if err != "" {
		t.Error(err)
	}
}

func TestGetRookMoves(t *testing.T) {
	err := testMoveFile(boardsDir + "rookTests.txt")
	if err != "" {
		t.Error(err)
	}
}

func TestGetQueenMoves(t *testing.T) {
	err := testMoveFile(boardsDir + "queenTests.txt")
	if err != "" {
		t.Error(err)
	}
}

func TestGetKnightMoves(t *testing.T) {
	err := testMoveFile(boardsDir + "knightTests.txt")
	if err != "" {
		t.Error(err)
	}
}

// func TestGetBoardFromString(t *testing.T) {
// 	b := GetBoardFromString(StartingBoard)
// 	fmt.Println(b)
// }

// func TestGetBoardFromFile(t *testing.T) {
// 	b := GetBoardFromFile(boardsDir + "kpo.txt")
// 	fmt.Println(b)

// 	b2 := GetBoardFromFile(boardsDir + "qgo.txt")
// 	fmt.Println(b2)
// }

// func TestGetWhiteCoords(t *testing.T) {
// 	b := GetBoardFromString(StartingBoard)
// 	fmt.Println("Checking white coordinates: ")
// 	fmt.Println(b)
// 	fmt.Println(b.GetWhiteCoords())
// }

// func TestGetBlackCoords(t *testing.T) {
// 	b := GetBoardFromString(StartingBoard)
// 	fmt.Println("Checking black coordinates: ")
// 	fmt.Println(b)
// 	fmt.Println(b.GetBlackCoords())
// }

func TestMoveEquals(t *testing.T) {
	a1 := coord{y: -5, x: 0}
	m1 := move{a: a1, b: coord{y: 1, x: 1}}
	a2 := a1.Copy()
	m2 := move{a: a2, b: coord{y: 1, x: 1}}
	if !m1.Equals(m2) {
		t.Error("moves not equal")
	}
}

func TestCoordAdd(t *testing.T) {
	ys := []int{0, 1, 2, 3, -5, -5, -5, 10000, 978, 67674664}
	xs := []int{-1, 1, -5, 0, 2, 3, -5, -5, -5, -100}
	testCoords := []coord{{y: 0, x: 0}, {y: -1, x: 0}, {y: 1000, x: 3}}
	for i := range len(ys) {
		y := ys[i]
		x := xs[i]
		for _, c := range testCoords {
			sum := c.Copy().Add(y, x)
			if y+c.y != sum.y || x+c.x != sum.x {
				t.Errorf("failure to add: %v", c)
			}
		}
	}
}

func TestCoordEquals(t *testing.T) {
	c1 := coord{y: -4, x: 1}
	c2 := coord{y: -4, x: 1}
	c3 := coord{y: 0, x: 0}
	if !c1.Equals(c2) {
		t.Error("should be equal")
	}
	if c1.Equals(c3) {
		t.Error("shouldn't be equal")
	}
}
