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

// Check if two move lists have the same highest scoring moves.
// If moves are checked in differing orders, it may turn out that lower scoring
// moves are not equal. This is because as soon as a move is known to be lower
// than the highest-score-so-far, getScore moves on to a different move.
func moveListsEffectivelyEqual(mList1 []board.ScoredMove, mList2 []board.ScoredMove) bool {
	if len(mList1) == len(mList2) {
		if len(mList1) == 0 {
			return true
		} else {
			mList1 = InsertionSort(mList1)
			highScore := mList1[0].GetScore()
			topScoreList1 := make([]board.ScoredMove, 0)
			topScoreList2 := make([]board.ScoredMove, 0)
			for _, m1 := range mList1 {
				if utility.IsClose(m1.GetScore(), highScore) {
					topScoreList1 = append(topScoreList1, m1)
				}
			}
			for _, m2 := range mList2 {
				if utility.IsClose(m2.GetScore(), highScore) {
					topScoreList2 = append(topScoreList2, m2)
				}
			}
			if len(topScoreList1) == len(topScoreList2) {
				for _, m1 := range topScoreList1 {
					matchFound := false
					for _, m2 := range topScoreList2 {
						if m1.GetMove().Equals(m2.GetMove()) {
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
