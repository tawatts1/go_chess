package ai

import (
	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
	"golang.org/x/exp/rand"
)

// type []board.ScoredMove struct {
// 	scores []float64
// 	moves  []board.Move
// 	size   int
// }
// move list is just []board.ScoredMove

func moveListsEqual(mList1 []board.ScoredMove, mList2 []board.ScoredMove) bool {
	if len(mList1) == len(mList2) {
		for i := range len(mList1) {
			matchFound := false
			for j := range len(mList1) {
				if mList1[i].Equals(mList2[j]) {
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

func newMoveList(moves []board.Move) []board.ScoredMove {
	l := len(moves)
	out := make([]board.ScoredMove, l, l)
	for i := range l {
		out[i] = board.NewScoredMove(moves[i])
	}
	return out
}

// after the scores are calculated, choose the best move. No need for this function to be efficient as it is called only once.
func GetMaxScoreMove(mList []board.ScoredMove) board.Move {
	if len(mList) > 0 {
		maxScore := mList[0].GetScore()
		for i := 1; i < len(mList); i++ {
			if mList[i].GetScore() > maxScore {
				maxScore = mList[i].GetScore()
			}
		}
		bestMoves := make([]board.Move, 0)
		for i := 0; i < len(mList); i++ {
			if utility.IsClose(maxScore, mList[i].GetScore()) {
				bestMoves = append(bestMoves, mList[i].GetMove())
			}
		}
		return bestMoves[rand.Intn(len(bestMoves))]
	} else {
		panic("No moves")
	}
}

func InsertionSort(mList []board.ScoredMove) []board.ScoredMove {
	for i := 1; i < len(mList); i++ {
		for j := i; j > 0; j-- {
			if mList[j].GreaterThan(mList[j-1]) { //.GetScore() > mList[j-1].GetScore() && !utility.IsClose(mList[j].GetScore(), mList[j-1].GetScore()) {
				mList[j], mList[j-1] = mList[j-1], mList[j]
			} else {
				break
			}
		}
	}
	return mList
}

func isSortedDesc(mList []board.ScoredMove) bool {
	for i := 1; i < len(mList); i++ {
		if !utility.IsApproxGreaterThanOrEq(mList[i-1].GetScore(), mList[i].GetScore()) {
			return false
		}
	}
	return true
}
