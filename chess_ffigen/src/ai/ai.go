package ai

import (
	"fmt"

	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
	"golang.org/x/exp/rand"
)

// func GetScoreAfterMove(b board.Board, m board.Move, isWhite bool) float64 {
// 	var out float64 = 0
// 	boardAfterMove := board.GetBoardAfterMove(b, m)
// 	out += getScoreFromBoard(boardAfterMove, isWhite)
// 	return out
// }

// func GetAiFromString(aiCode string) func(board.Board, bool) float64 {
// 	if aiCode == simpleAICode {
// 		return getScoreFromBoard
// 	} else {
// 		panic("ai code does not match")
// 	}
// }

type moveList struct {
	scores []float64
	moves  []board.Move
	size   int
}

func newMoveList(moves []board.Move) moveList {
	return moveList{moves: moves, scores: make([]float64, len(moves), len(moves)), size: len(moves)}
}

func (mList1 moveList) Equals(mList2 moveList) bool {
	if mList1.size == mList2.size {
		for i := range mList1.size {
			matchFound := false
			for j := range mList1.size {
				if mList1.moves[i].Equals(mList2.moves[j]) && mList1.scores[i] == mList2.scores[j] {
					matchFound = true
					break
				}
			}
			if !matchFound {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// after the scores are calculated, choose the best move. No need for this function to be efficient as it is called only once.
func (mList moveList) GetMaxScoreMove() board.Move {
	if mList.size > 0 {
		maxScore := mList.scores[0]
		for i := 1; i < mList.size; i++ {
			if mList.scores[i] > maxScore {
				maxScore = mList.scores[i]
			}
		}
		bestMoves := make([]board.Move, 0)
		for i := 0; i < mList.size; i++ {
			if utility.IsClose(maxScore, mList.scores[i]) {
				bestMoves = append(bestMoves, mList.moves[i])
			}
		}
		return bestMoves[rand.Intn(len(bestMoves))]
	} else {
		panic("No moves")
	}
}

// Calculate moves and their scores and return one of the moves with the max score
func ChooseMove(b board.Board, isWhite bool, depth int, scoringFunctionName string, useCache bool, verbosity int) board.Move {
	mList := newMoveList(b.GetLegalMoves(isWhite))
	cache := NewLruCacheList(depth+1, 800000)
	mList = ScoreSortMoveList(mList, b, isWhite, depth, 0, scoringFunctionName, cache, useCache)
	if verbosity > 0 {
		fmt.Println(cache)
	}
	return mList.GetMaxScoreMove()
}

// Score and sort the move list.
// arguments of note:
// depth: how many moves to look ahead.
// movesAhead: how many moves ahead it is already looking. This is used in caching.
// scoringFuncName: The identifier for the desired scoring function
// cache: Pointer to cache which is used to store board calculations
// useCache: Whether to read and write to the cache. Should be true except for testing cases.
func ScoreSortMoveList(mList moveList, b board.Board, isWhite bool,
	depth int, movesAhead int, scoringFuncName string,
	cache *lruCacheList, useCache bool) moveList {

	if depth == 0 {
		return mList
	} else if depth > 0 {
		if depth > 2 {
			mList = ScoreSortMoveList(mList, b, isWhite, depth-2, movesAhead, scoringFuncName, cache, useCache)
		}
		wcs := -utility.Infinity
		newDepth := depth - 1
		newMovesAhead := movesAhead + 1
		newColor := !isWhite
		var score_i float64
		for i := range mList.size {
			newBoard := board.GetBoardAfterMove(b, mList.moves[i])
			score_i = -GetScore(
				newBoard,
				newColor,
				newDepth,
				newMovesAhead,
				-wcs,
				scoringFuncName,
				cache,
				useCache)
			mList.scores[i] = score_i
			if score_i > wcs {
				wcs = score_i
			}
		}
		// Now sort moves list based on scores, in descending order.
		// Insertion sort
		for i := 1; i < mList.size; i++ {
			for j := i; j > 0; j-- {
				if mList.scores[j] > mList.scores[j-1] && !utility.IsClose(mList.scores[j], mList.scores[j-1]) {
					mList.scores[j], mList.scores[j-1] = mList.scores[j-1], mList.scores[j]
					mList.moves[j], mList.moves[j-1] = mList.moves[j-1], mList.moves[j]
				} else {
					break
				}
			}
		}
		if !mList.isSortedDesc() {
			panic("list not sorted properly")
		}
		return mList
	} else {
		panic(fmt.Sprintf("Invalid depth: %v", depth))
	}
}

func (mList moveList) isSortedDesc() bool {
	for i := 1; i < mList.size; i++ {
		if !utility.IsApproxGreaterThanOrEq(mList.scores[i-1], mList.scores[i]) {
			return false
		}
	}
	return true
}

// Get the score by looking 'depth' number of moves ahead.
// arguments of note:
// depth: how many moves to look ahead.
// movesAhead: how many moves ahead it is already looking. This is used in caching.
// parent_wcs: The worst case scenario one recursion step above. Used to cut a calculation short.
// scoringFuncName: The identifier for the desired scoring function
// cache: Pointer to cache which is used to store board calculations
// useCache: Whether to read and write to the cache. Should be true except for testing cases.
func GetScore(b board.Board, isWhite bool, depth, movesAhead int, parent_wcs float64,
	scoringFuncName string, cache *lruCacheList, useCache bool) float64 {

	var maxScore float64
	current_args := newScoreArgs(b, depth, isWhite)
	// check cache to see if this score has already been calculated.
	if useCache && depth > MinCacheDepth {
		current_cached, current_hit := cache.Get(movesAhead, current_args)
		if current_hit {
			return current_cached
		}
	}

	if depth == 0 {
		maxScore = getScoringFunction(scoringFuncName)(b, isWhite)
		if useCache && depth > MinCacheDepth {
			cache.SetNewKey(movesAhead, current_args, maxScore)
		}
		return maxScore
	} else if depth > 0 {
		wcs := -utility.Infinity
		maxScore = wcs
		mList := newMoveList(b.GetLegalMoves(isWhite))
		if depth > 2 {
			mList = ScoreSortMoveList(mList, b, isWhite, depth-2, movesAhead, scoringFuncName, cache, useCache)
		}
		addScoreToCache := true
		var score_i float64
		newDepth := depth - 1
		newColor := !isWhite
		newMovesAhead := movesAhead + 1
		for i := range mList.size {
			newBoard := board.GetBoardAfterMove(b, mList.moves[i])
			score_i = -GetScore(
				newBoard,
				newColor,
				newDepth,
				newMovesAhead,
				-wcs,
				scoringFuncName,
				cache,
				useCache)

			// Let us say if white does move A, the worst that will happen
			// after that is white gaining 5 points, so great for them!
			// Now white is considering move B. After move B, White pretends
			// to be black, and the value -5 is passed as the parent_wcs.
			// If after move B, black does move Alpha which can force a loss
			// of only 2 points, black would be glad since white didn't
			// do move A which would have been way worse for them.
			// So white should not consider move B anymore after seeing that
			// black could get a better deal.
			// the key metric there is that -2 > -5, or score_i > parent_wcs
			if score_i > maxScore {
				maxScore = score_i
				if score_i > parent_wcs && !utility.IsClose(score_i, parent_wcs) {
					// stopping because the score is too high. So the actual
					// score for board b is greater than or equal to 'score'
					addScoreToCache = false
					break
				}
				if score_i > wcs {
					wcs = score_i
				}
			}

		} // end for loop
		if addScoreToCache && useCache && depth > MinCacheDepth {
			cache.SetNewKey(movesAhead, current_args, maxScore)
		}

		return maxScore
	} else {
		panic("invalid depth")
	}
}
