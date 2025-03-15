package ai

import (
	"fmt"
	"sync"

	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
)

// Calculate moves and their scores and return one of the moves with the max score
func CalculateScores(b board.Board, isWhite bool, depth int, scoringFunctionName string, useMultiprocessing bool) []board.ScoredMove {
	mList := newMoveList(b.GetLegalMoves(isWhite))
	numCores := GetMaxNumCores(len(mList))
	if len(mList) == 0 {
		panic("Should not call calculateScores on a board with no possible moves.")
	}
	if useMultiprocessing && numCores > 1 && len(mList) > 4 && depth > 2 {
		if depth >= 2 {
			mList = ScoreSortMoveList(mList, b, isWhite, depth-2, scoringFunctionName)
		}
		L := len(mList)
		jobs := make(chan board.ScoredMove, L)
		finished_jobs := make(chan board.ScoredMove, L)
		scoreMutex := NewMutexScoreManager()
		var wg sync.WaitGroup
		for worker_i := 0; worker_i < numCores; worker_i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := range jobs {
					finished_jobs <- j.SetScore(-GetScoreMutex(
						board.GetBoardAfterMove(b, j.GetMove()),
						!isWhite,
						depth-1,
						scoreMutex,
						scoringFunctionName))
				}
			}()
		}
		// load the moves in the jobs channel
		for _, m := range mList {
			jobs <- m
		}
		close(jobs)       // close jobs so for loop in go func can start.
		wg.Wait()         //Wait until all workers are done.
		mList = mList[:0] // empty the moves list
		close(finished_jobs)
		for m := range finished_jobs {
			mList = append(mList, m)
		}
		mList = InsertionSort(mList)
		return mList
	} else {
		mList = ScoreSortMoveList(mList, b, isWhite, depth, scoringFunctionName)
		return mList
	}
}
func ChooseMove(b board.Board, isWhite bool, depth int, scoringFunctionName string, useMultiprocessing bool) board.Move {
	return GetMaxScoreMove(CalculateScores(b, isWhite, depth, scoringFunctionName, useMultiprocessing))
}

// Score and sort the move list.
func ScoreSortMoveList(mList []board.ScoredMove, b board.Board, isWhite bool, depth int, scoringFuncName string) []board.ScoredMove {
	// calculate score of moves for a certain depth, then return sorted []board.ScoredMove
	if depth == 0 {
		return mList
	} else if depth > 0 {
		if depth >= 2 {
			mList = ScoreSortMoveList(mList, b, isWhite, depth-2, scoringFuncName)
		}
		wcs := -utility.Infinity
		for i := range len(mList) {
			score_i := -GetScore(
				board.GetBoardAfterMove(b, mList[i].GetMove()),
				!isWhite,
				depth-1,
				-wcs,
				scoringFuncName)
			mList[i] = mList[i].SetScore(score_i)
			if score_i > wcs {
				wcs = score_i
			}
		}
		// Now sort moves list based on scores, in descending order.
		// Insertion sort
		mList = InsertionSort(mList)
		return mList
	} else {
		panic(fmt.Sprintf("Invalid depth: %v", depth))
	}
}

// Get the score by looking 'depth' number of moves ahead.
func GetScore(b board.Board, isWhite bool, depth int, parent_wcs float64, scoringFuncName string) float64 {
	if depth == 0 {
		return getScoringFunction(scoringFuncName)(b, isWhite)
	} else if depth > 0 {
		maxScore := -utility.Infinity
		moves, gameStatus := b.GetLegalMovesWithStatus(isWhite)
		mList := newMoveList(moves)
		if len(moves) > 0 {
			if depth >= 2 {
				mList = ScoreSortMoveList(mList, b, isWhite, depth-2, scoringFuncName)
			}
			var score_i float64
			for i := range len(mList) {
				score_i = -GetScore(
					board.GetBoardAfterMove(b, mList[i].GetMove()),
					!isWhite,
					depth-1,
					-maxScore,
					scoringFuncName)
				// Let us say if white does move A, the worst that will happen
				// after that is white gaining 5 points, so great for them!
				// Now white is considering move B. After move B, White pretends
				// to be black, and the value -5 is passed as the parent_wcs.
				// If after move B, black does move Alpha which can force a loss
				// of only 2 points, black would be glad since white didn't
				// do move A which would have been way worse for them.
				// So white should not consider move B anymore after seeing that
				// black could get a better deal.
				// the key metric there is that -2 > -5, or score > parent_wcs
				if score_i > maxScore {
					maxScore = score_i
					if score_i > parent_wcs && !utility.IsClose(score_i, parent_wcs) {
						break
					}
				}
			}
		} else {
			maxScore = getScoreFromGameStatus(gameStatus, depth)
		}

		return maxScore
	} else {
		panic("invalid depth")
	}
}

func getScoreFromGameStatus(status string, depth int) float64 {
	var out float64 = 0
	if status == board.StatusStaleMate {
		out = 0
	} else if status == board.StatusCheckMate {
		out = -(utility.Infinity * .9) - (0.1 * float64(depth))
	} else {
		panic(fmt.Sprintf("Invalid status to calculate score from game status: %v", status))
	}
	return out
}
