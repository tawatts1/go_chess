package ai

import (
	"fmt"

	"github.com/tawatts1/go_chess/board"
	"github.com/tawatts1/go_chess/utility"
	"golang.org/x/exp/rand"
)

var simpleAICode = "simple"

func GetScoreAfterMove(b board.Board, m board.Move, isWhite bool) float64 {
	var out float64 = 0
	boardAfterMove := board.GetBoardAfterMove(b, m)
	out += getScoreFromBoard(boardAfterMove, isWhite)
	return out
}

func GetAiFromString(aiCode string) func(board.Board, bool) float64 {
	if aiCode == simpleAICode {
		return getScoreFromBoard
	} else {
		panic("ai code does not match")
	}
}

type moveList struct {
	scores []float64
	moves  []board.Move
	size   int
}

func newMoveList(moves []board.Move) moveList {
	return moveList{moves: moves, scores: make([]float64, len(moves), len(moves)), size: len(moves)}
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

func ChooseMove(b board.Board, isWhite bool, depth int) board.Move {
	mList := newMoveList(b.GetLegalMoves(isWhite))
	mList = SortMoveList(mList, b, isWhite, depth)
	return mList.GetMaxScoreMove()
}

func SortMoveList(mList moveList, b board.Board, isWhite bool, depth int) moveList {
	// calculate score of moves for a certain depth, then return sorted moveList
	if depth == 0 {
		return mList
	} else if depth > 0 {
		if depth > 1 {
			mList = SortMoveList(mList, b, isWhite, depth-1)
		}
		wcs := -utility.Infinity
		for i := range mList.size {
			score_i := -GetScore(
				board.GetBoardAfterMove(b, mList.moves[i]),
				!isWhite,
				depth-1,
				-wcs)
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

func GetScore(b board.Board, isWhite bool, depth int, parent_wcs float64) float64 {
	wcs := -utility.Infinity
	if depth == 0 {
		return getScoreFromBoard(b, isWhite)
	} else if depth > 0 {
		moves := b.GetLegalMoves(isWhite)
		maxScore := wcs
		for i := range len(moves) {
			score := -GetScore(
				board.GetBoardAfterMove(b, moves[i]),
				!isWhite,
				depth-1,
				-wcs)
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
