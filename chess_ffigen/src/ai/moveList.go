package ai

import (
	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
	"golang.org/x/exp/rand"
)

type moveList struct {
	scores []float64
	moves  []board.Move
	size   int
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

func newMoveList(moves []board.Move) moveList {
	return moveList{moves: moves, scores: make([]float64, len(moves), len(moves)), size: len(moves)}
}

func (mList moveList) AddMoves(m ...board.Move) moveList {
	mList.moves = append(mList.moves, m...)
	mList.scores = append(mList.scores, 0)
	mList.size = len(mList.moves)
	return mList
}

func (mList1 moveList) Combine(mList2 moveList) moveList {
	return moveList{moves: append(mList1.moves, mList2.moves...),
		scores: append(mList1.scores, mList2.scores...),
		size:   mList1.size + mList2.size}
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

func (mList moveList) InsertionSort() moveList {
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
	return mList
}

func (mList moveList) isSortedDesc() bool {
	for i := 1; i < mList.size; i++ {
		if !utility.IsApproxGreaterThanOrEq(mList.scores[i-1], mList.scores[i]) {
			return false
		}
	}
	return true
}
