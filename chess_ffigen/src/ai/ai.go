package ai

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
)

type MutexScoreManager struct {
	lock  sync.RWMutex
	score float64
}

func NewMutexScoreManager() *MutexScoreManager {
	return &MutexScoreManager{score: utility.Infinity}
}

func (mtx *MutexScoreManager) Read() float64 {
	mtx.lock.RLock()
	defer mtx.lock.RUnlock()
	return mtx.score
}

// Update the score if the newScore is smaller
func (mtx *MutexScoreManager) Update(newScore float64) {
	currentScore := mtx.Read()
	if newScore < currentScore {
		mtx.lock.Lock()
		defer mtx.lock.Unlock()
		if newScore < mtx.score {
			mtx.score = newScore
		}
	}
}

func GetMaxNumCores(lenMoves int) int {
	numCores := (runtime.NumCPU() * 3) / 4 // don't use all the cores
	idealNumProcessesPerProcess := 2
	var maxCores int = lenMoves / idealNumProcessesPerProcess
	if numCores > maxCores {
		numCores = maxCores
	}
	if numCores < 1 {
		numCores = 1
	}
	return numCores
}

// Calculate moves and their scores and return one of the moves with the max score
func CalculateScores(b board.Board, isWhite bool, depth int, scoringFunctionName string, useMultiprocessing bool) []board.ScoredMove {
	mList := newMoveList(b.GetLegalMoves(isWhite))
	numCores := GetMaxNumCores(len(mList))
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
					//fmt.Printf("Worker %v starting job with presumed score %v\n", worker_i, j.GetScore())
					finished_jobs <- j.SetScore(-GetScoreMutex(
						board.GetBoardAfterMove(b, j.GetMove()),
						!isWhite,
						depth-1,
						scoreMutex,
						scoringFunctionName))
					//fmt.Printf("Worker %v finished the job!\n", worker_i)
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
		mList := newMoveList(b.GetLegalMoves(isWhite))
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
		return maxScore
	} else {
		panic("invalid depth")
	}
}

// Get the score by looking 'depth' number of moves ahead.
func GetScoreMutex(b board.Board, isWhite bool, depth int, mtx *MutexScoreManager, scoringFuncName string) float64 {
	if depth == 0 {
		return getScoringFunction(scoringFuncName)(b, isWhite)
	} else if depth > 0 {
		maxScore := -utility.Infinity
		var parent_wcs float64
		mList := newMoveList(b.GetLegalMoves(isWhite))
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
			if score_i > maxScore {
				maxScore = score_i
				parent_wcs = mtx.Read()
				if score_i > parent_wcs && !utility.IsClose(score_i, parent_wcs) {
					break
				}
			}
		}
		mtx.Update(maxScore)
		return maxScore
	} else {
		panic("invalid depth")
	}
}
