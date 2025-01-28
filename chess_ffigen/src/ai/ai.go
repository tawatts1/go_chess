package ai

import (
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
	if depth == 0 {
		return mList.GetMaxScoreMove()
	} else if depth > 0 {
		for i := range mList.size {
			mList.scores[i] = -GetScore(
				board.GetBoardAfterMove(b, mList.moves[i]),
				!isWhite,
				depth-1)
		}
		return mList.GetMaxScoreMove()
	} else {
		panic("This depth is not implemented. ")
	}
}

func GetScore(b board.Board, isWhite bool, depth int) float64 {
	if depth == 0 {
		return getScoreFromBoard(b, isWhite)
	} else if depth > 0 {
		moves := b.GetLegalMoves(isWhite)
		maxScore := -2 * utility.Infinity
		for i := range len(moves) {
			score := -GetScore(
				board.GetBoardAfterMove(b, moves[i]),
				!isWhite,
				depth-1)
			if score > maxScore {
				maxScore = score
			}
		}
		return maxScore
	} else {
		panic("invalid depth")
	}
}
