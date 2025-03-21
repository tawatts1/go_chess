package ai

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
)

var verbosity int = 0 // 0 - nothing, 1 - just input lines, 2 - everything
var testFolder = "testingMoves/"

var defaultUseMultiprocessing bool = true
var defaultUseShuffle bool = true

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hasError(e error) bool {
	return e != nil
}

func testAiMoveFile(fname string, multiprocess bool, shuffle bool) string {
	data, err := os.ReadFile(fname)
	check(err)
	lines := strings.Split(string(data), "\n")
	//var a ai
	var N int
	var b board.Board
	for lineIndex, line := range lines {

		if len(line) < 5 || []rune(line)[0] == '#' {
			continue
		}
		if verbosity > 0 {
			fmt.Println(line)
		}
		args := strings.Split(line, ",")
		if args[0] == "new" {
			//a = GetAiFromString(args[1])

			b = board.GetBoardFromString(args[2])
			if verbosity > 1 {
				fmt.Println(b)
			}
		} else if args[0] == "move" || args[0] == "notmove" {
			Nparsed, okN := utility.StrToInt(args[1])
			color := args[2]
			y1, ok1 := utility.StrToInt(args[3])
			x1, ok2 := utility.StrToInt(args[4])
			y2, ok3 := utility.StrToInt(args[5])
			x2, ok4 := utility.StrToInt(args[6])
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
			if hasError(ok1) || hasError(ok2) || hasError(ok3) || hasError(ok4) || hasError(okN) {
				panic(fmt.Sprintf("failed to parse line %v", lineIndex+1))
			} else {
				c1 = board.NewCoord(y1, x1)
				c2 = board.NewCoord(y2, x2)
				N = Nparsed
			}

			r := []rune(strings.TrimSpace(args[7]))
			if len(r) > 1 {
				return fmt.Sprintf("%v: special column must be one or zero characters. ", lineIndex+1)
			} else if len(r) == 1 {
				special = r[0]
			}
			mExpected := board.NewMove(c1, c2, special)
			mResult := ChooseMove(b, isWhite, N, ScoringDefaultPieceValue, multiprocess, shuffle)
			if args[0] == "move" && !mResult.Equals(mExpected) {
				return fmt.Sprintf("%v\nline %v: Expected %v but got %v", b, lineIndex+1, mExpected, mResult)
			} else if args[0] == "notmove" && mResult.Equals(mExpected) {
				return fmt.Sprintf("line %v: Expected anything but %v", lineIndex+1, mResult)
			}
		} else if args[0] == "verbosity" {
			v, vok := utility.StrToInt(args[1])
			if !hasError(vok) {
				verbosity = v
			}
		}
	}
	return ""
}

func TestChooseMove(t *testing.T) {
	err := testAiMoveFile(testFolder+"aiTests.txt", defaultUseMultiprocessing, defaultUseShuffle)
	if err != "" {
		t.Error(err)
	}
}

func TestChooseMoveEndgame(t *testing.T) {
	err := testAiMoveFile(testFolder+"endgameTests.txt", defaultUseMultiprocessing, defaultUseShuffle)
	if err != "" {
		t.Error(err)
	}
}

func TestChooseMoveOrderedCoords(t *testing.T) {
	err := testAiMoveFile(testFolder+"orderedCoordTests.txt", false, false)
	if err != "" {
		t.Error(err)
	}
}

func ContainsMove(moveSlice []board.ScoredMove, m1 board.Move) bool {
	for _, m2 := range moveSlice {
		if m2.GetMove().Equals(m1) {
			return true
		}
	}
	return false
}

func TestSortMoveList(t *testing.T) {
	b := board.GetBoardFromString("000k00000np0p0000000000000n000000P000000000N0000000P0P000000K000")
	mList := newMoveList(b.GetLegalMoves(true))
	mList = ScoreSortMoveList(mList, b, true, 1, ScoringDefaultPieceValue)
	pawnAttack := board.NewMove(board.NewCoord(4, 1), board.NewCoord(3, 2), 0)
	knightAttack := board.NewMove(board.NewCoord(5, 3), board.NewCoord(3, 2), 0)
	if utility.IsClose(mList[0].GetScore(), mList[1].GetScore()) &&
		ContainsMove(mList[0:2], pawnAttack) &&
		ContainsMove(mList[0:2], knightAttack) {
	} else {
		t.Error("failed n=1")
	}

	mList = ScoreSortMoveList(mList, b, true, 2, ScoringDefaultPieceValue)
	if !utility.IsClose(mList[0].GetScore(), mList[1].GetScore()) &&
		mList[0].GetMove().Equals(pawnAttack) &&
		mList[1].GetMove().Equals(knightAttack) {
	} else {
		t.Error("failed n=2")
	}
	// after sorting with depth=3, either attack is fine, but the pawn attack
	// should still be ordered first because it is better for depth=2.
	mList = ScoreSortMoveList(mList, b, true, 3, ScoringDefaultPieceValue)
	if utility.IsClose(mList[0].GetScore(), mList[1].GetScore()) &&
		mList[0].GetMove().Equals(pawnAttack) &&
		mList[1].GetMove().Equals(knightAttack) {
	} else {
		t.Error("failed n=3")
	}
}

func TestSortScoredMoveListPositionScoring(t *testing.T) {
	//opening board after white puts a random pawn forward:
	b := board.GetBoardFromString("onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO")
	if !utility.IsClose(getPositionScoreValue(b, false), 0) {
		t.Error("error")
	}
	b_1bp := board.GetBoardFromString("0000k00000p0000000000000000000000000000000000000000000000000K000")
	//fmt.Println(b_1bp)
	if !utility.IsClose(getPositionScoreValue(b_1bp, false), 1+positionScalingFactor*.5) {
		t.Error("error")
	}

	if !utility.IsClose(getPositionScoreValue(b_1bp, true), -(1 + positionScalingFactor*.5)) {
		t.Error("error")
	}

	b_1wp := board.GetBoardFromString("0000k0000000000000000000000000000000000000000000000000P00000K000")
	if !utility.IsClose(getPositionScoreValue(b_1wp, false), -(1 + positionScalingFactor*.5)) {
		t.Error("error")
	}
	//moving knight should give the same score as moving a pawn two spaces.
	bAfterKnight := board.GetBoardAfterMove(b, board.NewMove(board.NewCoord(0, 1), board.NewCoord(2, 2), 0))
	bAfterPawn := board.GetBoardAfterMove(b, board.NewMove(board.NewCoord(1, 0), board.NewCoord(3, 0), 0))
	if !utility.IsClose(getPositionScoreValue(bAfterKnight, false),
		getPositionScoreValue(bAfterPawn, false)) {
		t.Error("error")
	}
}

func TestGetPositionScore(t *testing.T) {
	testSingleFunc := func(y, x int, isWhite bool, p rune, expectedScore float64) {
		result := GetPositionScore(board.NewCoord(y, x), p)
		if !(utility.IsClose(expectedScore, result)) {
			t.Errorf("Expected %v but got %v ((%v, %v), %v, %c)", expectedScore, result, x, y, isWhite, p)
		}
	}
	testBothFunc := func(y, x int, isWhite bool, p rune, expectedScore float64, testBothColors, testSymmetries bool) {
		testSingleFunc(y, x, isWhite, p, expectedScore)
		if testBothColors {
			testSingleFunc(y, x, !isWhite, p, expectedScore)
		}
		if testSymmetries {
			testSingleFunc(y, 7-x, isWhite, p, expectedScore)
			testSingleFunc(7-y, x, isWhite, p, expectedScore)
			testSingleFunc(7-y, 7-x, isWhite, p, expectedScore)
			testSingleFunc(x, y, isWhite, p, expectedScore)
			testSingleFunc(x, 7-y, isWhite, p, expectedScore)
			testSingleFunc(7-x, y, isWhite, p, expectedScore)
			testSingleFunc(7-x, 7-y, isWhite, p, expectedScore)
		}
	}
	// PAWNS
	testBothFunc(1, 4, false, 'p', 0.5, false, false)
	testBothFunc(2, 4, false, 'p', 0, false, false)
	testBothFunc(3, 4, false, 'p', 1, false, false)
	testBothFunc(4, 4, false, 'p', 2, false, false)
	testBothFunc(5, 4, false, 'p', 3, false, false)
	testBothFunc(6, 4, false, 'p', 4, false, false)
	testBothFunc(7-1, 5, true, 'P', 0.5, false, false)
	testBothFunc(7-2, 5, true, 'P', 0, false, false)
	testBothFunc(7-3, 5, true, 'P', 1, false, false)
	testBothFunc(7-4, 5, true, 'P', 2, false, false)
	testBothFunc(7-5, 5, true, 'P', 3, false, false)
	testBothFunc(7-6, 5, true, 'P', 4, false, false)
	// BISHOPS
	testBothFunc(1, 4, false, 'b', 0.5, true, true)
	testBothFunc(2, 3, false, 'b', 1, true, true)
	testBothFunc(3, 4, false, 'b', 1.5, true, true)
	testBothFunc(4, 3, false, 'b', 1.5, true, true)
	testBothFunc(4, 4, false, 'b', 1.5, true, true)
	testBothFunc(5, 4, false, 'b', 1, true, true)
	testBothFunc(0, 7, false, 'b', 0.5, true, true)
	// QUEENS
	testBothFunc(1, 4, false, 'q', 0.5, true, true)
	testBothFunc(2, 3, false, 'q', 1, true, true)
	testBothFunc(3, 4, false, 'q', 1.5, true, true)
	testBothFunc(4, 3, false, 'q', 1.5, true, true)
	testBothFunc(4, 4, false, 'q', 1.5, true, true)
	testBothFunc(5, 4, false, 'q', 1, true, true)
	testBothFunc(0, 7, false, 'q', 0.5, true, true)
	// KNIGHTS
	testBothFunc(0, 0, false, 'n', 0.25, true, true)
	testBothFunc(1, 0, false, 'n', 3.0/8.0, true, false)
	testBothFunc(6, 0, false, 'n', 3.0/8.0, true, false)
	testBothFunc(1, 7, false, 'n', 3.0/8.0, true, false)
	testBothFunc(6, 7, false, 'n', 3.0/8.0, true, false)
	testBothFunc(0, 1, false, 'n', 0.5, true, false)
	testBothFunc(0, 6, false, 'n', 0.5, true, false)
	testBothFunc(7, 1, false, 'n', 0.5, true, false)
	testBothFunc(7, 6, false, 'n', 0.5, true, false)
	testBothFunc(1, 4, false, 'n', 0.75, true, true)
	testBothFunc(1, 1, false, 'n', 0.5, true, true)
	testBothFunc(2, 2, false, 'n', 1, true, true)
	testBothFunc(2, 3, false, 'n', 1, true, true)
	testBothFunc(4, 4, false, 'n', 1, true, true)
	testBothFunc(6, 4, false, 'n', 0.75, true, true)
}

func testWithoutMultiprocessing(b board.Board, isWhite bool, depth int, scoringFunctionName string) string {
	mList1 := CalculateScores(b, isWhite, depth, scoringFunctionName, false)
	mList2 := CalculateScores(b, isWhite, depth, scoringFunctionName, true)

	if !moveListsEffectivelyEqual(mList1, mList2) {
		return "Using Multiprocessing changed score"
	} else {
		return ""
	}
}

func TestMultiprocessing(t *testing.T) {
	boards := []string{
		"onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO",
		"00000000ppp0000000000k000000000000000PPP00000000000000000K000000",
		"0k000000p00000000n00p000000pP000000P0P000000N00000000000000000K0",
		"0k000000p000n00000000000000000000000000P00000000000N0000000000K0",
		"0000000000p0000p00000k00p00000000000000P00K00000P0000P0000000000",
		"0000k00000000000000ppp00000000000000000000PPP00000000000000K0000",
		"00b0k0000n00000000ppppr00P000000000PPP00000BNB00000000000000K000",
		"0000kn000000000000000000000000000000000000000000p00PPP000000K00O",
	}
	depth := 3
	depthMap := make(map[int]int)
	depthMap[0] = 3
	depthMap[2] = 3
	depthMap[3] = 4
	depthMap[4] = 2
	colors := []bool{false, true}
	for _, b := range boards {
		for _, c := range colors {
			err := testWithoutMultiprocessing(board.GetBoardFromString(b),
				c, depth, ScoringPiecePositionValue)
			if err != "" {
				t.Error(err)
			}
			depth = depthMap[depth]
		}
	}
}

func BenchmarkOpening3(b *testing.B) {
	startingBoard := board.GetBoardFromString("onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO")
	for n := 0; n < b.N; n++ {
		ChooseMove(startingBoard, true, 3, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
	}
}

func BenchmarkOpening4(b *testing.B) {
	startingBoard := board.GetBoardFromString("onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO")
	for n := 0; n < b.N; n++ {
		ChooseMove(startingBoard, true, 4, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
	}
}

func BenchmarkPawns7(b *testing.B) {
	startingBoard := board.GetBoardFromString("00000000ppp0000000000k000000000000000PPP00000000000000000K000000")
	for n := 0; n < b.N; n++ {
		ChooseMove(startingBoard, true, 7, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
	}
}

func BenchmarkSimpleBoards6(b *testing.B) {
	simpleBoards := []string{"00000000ppp0000000000k000000000000000PPP00000000000000000K000000",
		"0k000000p00000000n00p000000pP000000P0P000000N00000000000000000K0",
		"0k000000p000n00000000000000000000000000P00000000000N0000000000K0",
		"0000000000p0000p00000k00p00000000000000P00K00000P0000P0000000000",
		"0000k00000000000000ppp00000000000000000000PPP00000000000000K0000",
	}
	boards := make([]board.Board, 0)
	for _, boardString := range simpleBoards {
		boards = append(boards, board.GetBoardFromString(boardString))
	}
	for n := 0; n < b.N; n++ {
		for _, startingBoard := range boards {
			ChooseMove(startingBoard, true, 6, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
		}

	}
}

func BenchmarkFork4and5(b *testing.B) {
	startingBoard := board.GetBoardFromString("onb0kb0opp000ppp00000n000N0qp0000000000000000Q00PPPP0PPPO0B0KB0O")
	for n := 0; n < b.N; n++ {
		ChooseMove(startingBoard, true, 4, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
		ChooseMove(startingBoard, true, 5, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
	}
}

func BenchmarkCaptureChains5(b *testing.B) {
	sb1 := board.GetBoardFromString("00b0k0000n00000000ppppr00P000000000PPP00000BNB00000000000000K000")
	sb2 := board.GetBoardFromString("000nk00000npp00000b00p00R00p00000000P00000NP0000000PPP000000K000")
	for n := 0; n < b.N; n++ {
		ChooseMove(sb1, true, 5, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
		ChooseMove(sb2, true, 5, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
	}
}

func BenchmarkBishopsVsRook5(b *testing.B) {
	sb := board.GetBoardFromString("0k00r0000000pppp00000000000000000000000000000000PPP00000000BB0K0")
	for n := 0; n < b.N; n++ {
		ChooseMove(sb, true, 5, ScoringPiecePositionValue, defaultUseMultiprocessing, true)
	}
}
