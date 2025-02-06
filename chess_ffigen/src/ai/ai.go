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
	return &MutexScoreManager{score: -utility.Infinity}
}

func (mtx *MutexScoreManager) Read() float64 {
	mtx.lock.RLock()
	defer mtx.lock.RUnlock()
	return mtx.score
}

func (mtx *MutexScoreManager) Update(newScore float64) {
	currentScore := mtx.Read()
	if newScore > currentScore {
		mtx.lock.Lock()
		defer mtx.lock.Unlock()
		if newScore > mtx.score {
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
func CalculateScores(b board.Board, isWhite bool, depth int, scoringFunctionName string, useMultiprocessing bool) moveList {
	mList := newMoveList(b.GetLegalMoves(isWhite))
	numCores := GetMaxNumCores(mList.size)
	if useMultiprocessing && numCores > 1 && mList.size > 4 {
		if depth > 1 { //!!! should this be depth>2 ?
			mList = ScoreSortMoveList(mList, b, isWhite, depth-2, scoringFunctionName)
		}
		slcMList := make([]moveList, numCores)
		for i, m := range mList.moves {
			index := i % numCores
			slcMList[index] = slcMList[index].AddMoves(m)
		}
		scoredMLists := make(chan moveList, numCores)
		var wg sync.WaitGroup
		scoreMutex := NewMutexScoreManager()
		for _, coreMList := range slcMList {
			wg.Add(1)
			go func() {
				defer wg.Done()
				scoredMLists <- ScoreSortMoveListMutex(coreMList, b, isWhite, depth, scoringFunctionName, scoreMutex)
			}()
		}
		resultMList := newMoveList(make([]board.Move, 0))
		wg.Wait()           // wait for all cpus to finish the processes
		close(scoredMLists) // close the channel to new input.
		for ml := range scoredMLists {
			resultMList = resultMList.Combine(ml)
		}
		resultMList = resultMList.InsertionSort()
		return resultMList

	} else {
		mList = ScoreSortMoveList(mList, b, isWhite, depth, scoringFunctionName)
		return mList
	}
}
func ChooseMove(b board.Board, isWhite bool, depth int, scoringFunctionName string, useMultiprocessing bool) board.Move {
	return CalculateScores(b, isWhite, depth, scoringFunctionName, useMultiprocessing).GetMaxScoreMove()
}

// Score and sort the move list.
func ScoreSortMoveList(mList moveList, b board.Board, isWhite bool, depth int, scoringFuncName string) moveList {
	// calculate score of moves for a certain depth, then return sorted moveList
	if depth == 0 {
		return mList
	} else if depth > 0 {
		if depth > 1 { //!!! should this be depth>2 ?
			mList = ScoreSortMoveList(mList, b, isWhite, depth-2, scoringFuncName)
		}
		wcs := -utility.Infinity
		for i := range mList.size {
			score_i := -GetScore(
				board.GetBoardAfterMove(b, mList.moves[i]),
				!isWhite,
				depth-1,
				-wcs,
				scoringFuncName)
			mList.scores[i] = score_i
			if score_i > wcs {
				wcs = score_i
			}
		}
		// Now sort moves list based on scores, in descending order.
		// Insertion sort
		mList = mList.InsertionSort()
		return mList
	} else {
		panic(fmt.Sprintf("Invalid depth: %v", depth))
	}
}

// Score and sort the move list.
func ScoreSortMoveListMutex(mList moveList, b board.Board, isWhite bool, depth int, scoringFuncName string, scoreMutex *MutexScoreManager) moveList {
	// calculate score of moves for a certain depth, then return sorted moveList
	if depth == 0 {
		return mList
	} else if depth > 0 {
		// No presorting because that was already done outside this function.
		var wcs float64
		for i := range mList.size {
			wcs = scoreMutex.Read()
			score_i := -GetScore(
				board.GetBoardAfterMove(b, mList.moves[i]),
				!isWhite,
				depth-1,
				-wcs,
				scoringFuncName)
			mList.scores[i] = score_i
			if score_i > wcs {
				scoreMutex.Update(score_i)
			}
		}
		// Now sort moves list based on scores, in descending order.
		// Insertion sort
		mList = mList.InsertionSort()
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
		wcs := -utility.Infinity
		maxScore := wcs
		mList := newMoveList(b.GetLegalMoves(isWhite))
		if depth > 2 { // should this be depth>2 or 1?
			mList = ScoreSortMoveList(mList, b, isWhite, depth-2, scoringFuncName)
		}
		for i := range mList.size {
			score := -GetScore(
				board.GetBoardAfterMove(b, mList.moves[i]),
				!isWhite,
				depth-1,
				-wcs,
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
			if score > maxScore {
				maxScore = score
				if score > parent_wcs && !utility.IsClose(score, parent_wcs) {
					break
				}
				if score > wcs {
					wcs = score
				}
			}
		}

		return maxScore
	} else {
		panic("invalid depth")
	}
}
