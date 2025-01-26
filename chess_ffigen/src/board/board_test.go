package board

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

var boardsDir string = "testingBoards/"
var verbose = false

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
	var b Board
	var bcm, wcm map[Coord]bool
	for lineIndex, line := range lines {
		if verbose {
			fmt.Println(line)
		}
		if len(line) < 5 || []rune(line)[0] == '#' {
			continue
		}
		args := strings.Split(line, ",")
		if args[0] == "new" {
			b = GetBoardFromString(args[1])
			bcm = b.GetBlackCoordMap()
			wcm = b.GetWhiteCoordMap()
			if verbose {
				fmt.Println(b)
			}

		} else if args[0] == "num" {
			y_, ok1 := strconv.Atoi(strings.TrimSpace(args[1]))
			x_, ok2 := strconv.Atoi(strings.TrimSpace(args[2]))
			expected_len, ok3 := strconv.Atoi(strings.TrimSpace(args[3]))
			var c Coord
			if hasError(ok1) || hasError(ok2) || hasError(ok3) {
				panic(fmt.Sprintf("failed to parse line %v", lineIndex+1))
			} else {
				c = Coord{y: y_, x: x_}
			}
			var actual_len int
			if IsBlack(b.GetPiece(c)) {
				actual_len = len(b.GetMoves(bcm, wcm, c, true))
			} else {
				actual_len = len(b.GetMoves(wcm, bcm, c, true))
			}
			if actual_len != expected_len {
				return fmt.Sprintf("%v: expected %v, got %v", lineIndex+1, expected_len, actual_len)
			}
		} else if args[0] == "has" {
			y1, ok1 := strconv.Atoi(strings.TrimSpace(args[1]))
			x1, ok2 := strconv.Atoi(strings.TrimSpace(args[2]))
			y2, ok3 := strconv.Atoi(strings.TrimSpace(args[3]))
			x2, ok4 := strconv.Atoi(strings.TrimSpace(args[4]))
			var special rune
			var c1, c2 Coord
			var m Move
			if hasError(ok1) || hasError(ok2) || hasError(ok3) || hasError(ok4) {
				panic(fmt.Sprintf("failed to parse line %v", lineIndex+1))
			} else {
				c1 = Coord{y: y1, x: x1}
				c2 = Coord{y: y2, x: x2}
			}
			if len(args) >= 6 {
				r := []rune(args[5])
				if len(r) > 1 {
					return fmt.Sprintf("%v: special column must be one or zero characters. ", lineIndex+1)
				} else if len(r) == 1 {
					special = r[0]
				}
			}
			m = Move{a: c1, b: c2, special: special}
			var moves []Move
			if IsBlack(b.GetPiece(c1)) {
				moves = b.GetMoves(bcm, wcm, c1, true)
			} else {
				moves = b.GetMoves(wcm, bcm, c1, true)
			}

			if !Contains(moves, m) {
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

func TestGetKingMoves(t *testing.T) {
	err := testMoveFile(boardsDir + "kingTests.txt")
	if err != "" {
		t.Error(err)
	}
}

func TestMovingInOutOfCheck(t *testing.T) {
	err := testMoveFile(boardsDir + "checkTests.txt")
	if err != "" {
		t.Error(err)
	}
}

func TestCastleMoves(t *testing.T) {
	err := testMoveFile(boardsDir + "castleTests.txt")
	if err != "" {
		t.Error(err)
	}
}
func TestIsInCheck(t *testing.T) {
	test := func(bstring string, whiteresult, blackresult bool) {
		b := GetBoardFromString(bstring)
		//fmt.Println(b)
		bcm := b.GetBlackCoordMap()
		wcm := b.GetWhiteCoordMap()
		if b.IsInCheck(wcm, bcm, b.GetKingCoord(wcm)) != whiteresult {
			t.Error("failed to detect check")
		}
		if b.IsInCheck(bcm, wcm, b.GetKingCoord(bcm)) != blackresult {
			t.Error("Detected check when there was none")
		}
	}
	test("onbq00no0ppppkbp00000000p0000000000K000000000P00PPPPP0PPONBQ0BNO", true, false)
	test("0000000000000k000000P00000000000000000b0000000000000NPP00r000K00", true, true)
	test("0000000Q00000k00000Q0P000pN00Q000Pn00qb0000000000000NPPr0R000K00", false, false)
}

func TestGetBoardAfterMove(t *testing.T) {
	test := func(b1str, b2str string, m Move) {
		b1 := GetBoardFromString(b1str)
		b2 := GetBoardAfterMove(b1, m)
		if b2.String() != GetBoardFromString(b2str).String() {
			fmt.Println(GetBoardFromString(b2str).String())
			fmt.Println(b2)
			t.Error("Move not effective")
		}
	}
	test("o000k00op0ppnpbpbpn00qp00000p0000000P0000QNP0P0NP00BB0PPO000K00O",
		"o000k00op0ppnQbpbpn00qp00000p0000000P00000NP0P0NP00BB0PPO000K00O",
		Move{a: Coord{y: 5, x: 1}, b: Coord{y: 1, x: 5}})
	test("onbqkbnopppp0ppp000000000000p0000000P00000000000PPPP0PPPONBQK00O",
		"onbqkbnopppp0ppp000000000000p0000000P00000000000PPPP0PPPRNBQ0RK0",
		Move{a: Coord{y: 7, x: 4}, b: Coord{y: 7, x: 6}, special: CastleBridge})
	test("onbqkbnopppp0ppp000000000000p0000000P00000000000PPPP0PPPRNBQK00O",
		"onbqkbnopppp0ppp000000000000p0000000P00000000000PPPPKPPPRNBQ000R",
		Move{a: Coord{y: 7, x: 4}, b: Coord{y: 6, x: 4}, special: WhiteKing})
	test("onbqk00opppp0ppp00000n0000b0p000000000000000P000PPPP0PPPONBQK0R0",
		"rnbq0rk0pppp0ppp00000n0000b0p000000000000000P000PPPP0PPPONBQK0R0",
		Move{a: Coord{y: 0, x: 4}, b: Coord{y: 0, x: 6}, special: CastleBridge})
	test("onbqkbnopp000ppp0000000000aPp0000000000000000000PPPP0PPPONBQKBNO",
		"onbqkbnopp000ppp000P000000p0p0000000000000000000PPPP0PPPONBQKBNO",
		Move{a: Coord{y: 3, x: 3}, b: Coord{y: 2, x: 3}})
	test("0000k00000q000000000000000000000000pA00000000000000P0P000000K000",
		"0000k000000000000000000000000000000pP00000000000000P0P0000q0K000",
		Move{a: Coord{y: 1, x: 2}, b: Coord{y: 7, x: 2}})
	test("o00qkb0op0pp0pppnp000p000000000Q00BP000000000000PPP00P0PONB0K0Nb",
		"o00qkb0op0pp0pppnp000p000000000Q00BP000000000000PPPK0P0PRNB000Nb",
		Move{a: Coord{y: 7, x: 4}, b: Coord{y: 6, x: 3}, special: WhiteKing})

}

func TestGetGameStatus(t *testing.T) {
	test := func(b1str, status string, isWhite bool) {
		b1 := GetBoardFromString(b1str)
		calcStatus := GetGameStatus(b1, isWhite)
		if calcStatus != status {
			t.Errorf("expected (%v) but got (%v)", status, calcStatus)
		}
	}
	test("onbqkbnopppp0Qpp000000000000p00000B0P00000000000PPPP0PPPONB0K0NO", StatusCheckMate, false)
	test("onbqkbnopppp00pp000000000000p0Q000B0P00000000000PPPP0PPPONB0K0NO", StatusBlackMove, false)
	test("0000k00000000000000000000000000000000000000000000r000000r000K000", StatusCheckMate, true)
	test("0000k000000000000000000000000000000r0r00000000000r0000000000K000", StatusStaleMate, true)
	test("00000k000000000000000NPP000000000B0r0r00000000000r0000000000K000", StatusBlackMove, false)
	test("00000k000000000000000NPP000P00000BPr0r000P0000000r0000000000K000", StatusCheckMate, false)
	test("00000k00000000000000000P00000000000r0r00000000000r0000000000K000", StatusWhiteMove, true)
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
	a1 := Coord{y: -5, x: 0}
	m1 := Move{a: a1, b: Coord{y: 1, x: 1}}
	a2 := a1.Copy()
	m2 := Move{a: a2, b: Coord{y: 1, x: 1}}
	if !m1.Equals(m2) {
		t.Error("moves not equal")
	}
}

func TestCoordAdd(t *testing.T) {
	ys := []int{0, 1, 2, 3, -5, -5, -5, 10000, 978, 67674664}
	xs := []int{-1, 1, -5, 0, 2, 3, -5, -5, -5, -100}
	testCoords := []Coord{{y: 0, x: 0}, {y: -1, x: 0}, {y: 1000, x: 3}}
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
	c1 := Coord{y: -4, x: 1}
	c2 := Coord{y: -4, x: 1}
	c3 := Coord{y: 0, x: 0}
	if !c1.Equals(c2) {
		t.Error("should be equal")
	}
	if c1.Equals(c3) {
		t.Error("shouldn't be equal")
	}
}
