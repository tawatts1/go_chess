package ai

import (
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
	idealNumJobsPerWorker := 3
	var maxCores int = lenMoves / idealNumJobsPerWorker
	if numCores > maxCores {
		numCores = maxCores
	}
	if numCores < 1 {
		numCores = 1
	}
	return numCores
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
