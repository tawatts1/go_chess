package ai

import (
	"testing"

	"github.com/tawatts1/go_chess/board"
)

func TestLruCache(t *testing.T) {
	// making the testing cache size such that it is filled completely in the test
	var testingLruCachePtr = NewLRUCache(69)
	cacheList := NewLruCacheList(10, 10)
	whiteBool := true
	depthInt := 2
	openingBoard := board.GetBoardFromString("onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO")
	pawnBoard := board.GetBoardFromString("00000000ppp0000000000k000000000000000PPP00000000000000000K000000")
	forkBoard := board.GetBoardFromString("onb0kb0opp000ppp00000n000N0qp0000000000000000Q00PPPP0PPPO0B0KB0O")
	openMoves := ScoreSortMoveList(newMoveList(openingBoard.GetLegalMoves(whiteBool)), openingBoard, whiteBool, depthInt, 0, ScoringDefaultPieceValue, cacheList, true)
	pawnMoves := ScoreSortMoveList(newMoveList(pawnBoard.GetLegalMoves(whiteBool)), pawnBoard, whiteBool, depthInt, 0, ScoringDefaultPieceValue, cacheList, true)
	forkMoves := ScoreSortMoveList(newMoveList(forkBoard.GetLegalMoves(whiteBool)), forkBoard, whiteBool, depthInt, 0, ScoringDefaultPieceValue, cacheList, true)
	//fmt.Printf("Opening, pawn, fork moves lengths: %v, %v, %v\n", openMoves.size, pawnMoves.size, forkMoves.size)
	addMoves := func(mList moveList, b board.Board, isWhite bool, depth int) {
		for i, m := range mList.moves {
			newBoard := board.GetBoardAfterMove(b, m)
			testingLruCachePtr.SetNewKey(newScoreArgs(newBoard, depth, isWhite), mList.scores[i])
		}
	}
	addMoves(openMoves, openingBoard, whiteBool, depthInt)
	addMoves(pawnMoves, pawnBoard, whiteBool, depthInt)
	addMoves(forkMoves, forkBoard, whiteBool, depthInt)
	//fmt.Printf("Current Cache size: %v\n", len(testingLruCachePtr.keyStack))
	verifyInCache := func(mList moveList, b board.Board, isWhite bool, depth int) {
		for i, m := range mList.moves {
			newBoard := board.GetBoardAfterMove(b, m)
			if !testingLruCachePtr.Has(newScoreArgs(newBoard, depth, isWhite)) {
				t.Error("Arg not in cache")
			} else {
				if val, _ := testingLruCachePtr.Get(newScoreArgs(newBoard, depth, isWhite)); val != mList.scores[i] {
					t.Error("value not equal")
				}
			}
		}
	}
	verifyInCache(openMoves, openingBoard, whiteBool, depthInt)
	verifyInCache(pawnMoves, pawnBoard, whiteBool, depthInt)
	verifyInCache(forkMoves, forkBoard, whiteBool, depthInt)
}

func TestLruCache2(t *testing.T) {
	// making the testing cache size such that it is filled completely in the test
	var testingLruCachePtr = NewLRUCache(9)
	boards := make([]board.Board, 0)
	boards = append(boards, board.GetBoardFromString("Nnbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO"))
	boards = append(boards, board.GetBoardFromString("B0000000ppp0000000000k000000000000000PPP00000000000000000K000000"))
	boards = append(boards, board.GetBoardFromString("Rnb0kb0opp000ppp00000n000N0qp0000000000000000Q00PPPP0PPPO0B0KB0O"))
	colors := [2]bool{false, true}
	depths := [5]int{1, 2, 3, 4, 5}
	testingArgs := make([]getScoreArgs, 0)
	testingValues := make([]float64, 0)
	var j float64 = 0
	for _, b := range boards {
		for _, d := range depths {
			for _, c := range colors {
				testingArgs = append(testingArgs, newScoreArgs(b, d, c))
				testingValues = append(testingValues, j)
				j += 1
			}
		}
	}
	//testing that adding a key does indeed add the key
	for i, arg := range testingArgs {
		testingLruCachePtr.SetNewKey(arg, testingValues[i])
		if !testingLruCachePtr.Has(arg) {
			t.Errorf("Setting new key failed, %v: %v\n", i, arg)
		}
		if len(testingLruCachePtr.keyStack) == testingLruCachePtr.maxSize {
			if !testingLruCachePtr.keyStack[testingLruCachePtr.maxSize-1].Equals(arg) {
				t.Errorf("arg should be in the last spot in stack, %v: %v\n", i, arg)
			}
		}
	}
	//testing that accessing a value will reorder the stack
	for i, arg := range testingArgs {
		lastIndex := len(testingLruCachePtr.keyStack) - 1
		//skip some args using i%2==0
		if testingLruCachePtr.Has(arg) && i%2 == 0 {
			_, ok := testingLruCachePtr.Get(arg)
			if !ok {
				t.Error("ok should be true since has is true. ")
			} else {
				if !testingLruCachePtr.Has(arg) {
					t.Error("The get somehow dropped the arg")
				}
				if !arg.Equals(testingLruCachePtr.keyStack[lastIndex]) {
					t.Error("Failed to move arg to last index. ")
				}
			}

		}
	}

	// for i, arg := range testingArgs {
	// 	has := testingLruCachePtr.Has(arg)
	// 	//if len(testingArgs)-testingLruCachePtr.maxSize < i {
	// 	fmt.Printf("%v, %v\n", i, has)
	// 	//}
	// }

}

func testWithoutCache(b board.Board, isWhite bool, depth int, scoringFunctionName string) string {
	mList1 := newMoveList(b.GetLegalMoves(isWhite))

	cacheList := NewLruCacheList(depth+1, 20*20*20)
	mList1 = ScoreSortMoveList(mList1, b, isWhite, depth, 0, scoringFunctionName, cacheList, false)
	mList2 := newMoveList(b.GetLegalMoves(isWhite))
	cacheList = NewLruCacheList(depth+1, 20*20*20)
	mList2 = ScoreSortMoveList(mList2, b, isWhite, depth, 0, scoringFunctionName, cacheList, true)
	//fmt.Println(cacheList)
	if !mList1.Equals(mList2) {
		return "Using cache changed score"
	} else {
		return ""
	}
}

func TestPawns3(t *testing.T) {
	// series of boards with very few pieces
	boards := []string{"00000000ppp0000000000k000000000000000PPP00000000000000000K000000",
		"0k000000p00000000n00p000000pP000000P0P000000N00000000000000000K0",
		"0k000000p000n00000000000000000000000000P00000000000N0000000000K0",
		"0000000000p0000p00000k00p00000000000000P00K00000P0000P0000000000",
	}
	depths := []int{0, 3, 4}
	for _, boardStr := range boards {
		for _, d := range depths {
			startingBoard := board.GetBoardFromString(boardStr)
			testWithoutCache(startingBoard, true, d, ScoringPiecePositionValue)
		}
	}
	for _, boardStr := range boards {
		startingBoard := board.GetBoardFromString(boardStr)
		testWithoutCache(startingBoard, true, 4, ScoringPiecePositionValue)
		testWithoutCache(startingBoard, false, 4, ScoringPiecePositionValue)
	}
}
