package ai

import (
	"github.com/tawatts1/go_chess/board"
	"golang.org/x/exp/rand"
)

var simpleAICode = "simple"

type ai struct {
	s0 boardPieceScore
	//eventually the scoring ai will be more complicated, including other scoring methods like position.
}

func (a ai) GetScoreAfterMove(b board.Board, m board.Move, isWhite bool) float64 {
	var out float64 = 0
	boardAfterMove := board.GetBoardAfterMove(b, m)
	out += a.s0.getScore(boardAfterMove, isWhite)
	return out
}

func GetAiFromString(aiCode string) ai {
	if aiCode == simpleAICode {
		return ai{boardPieceScore{pieceMap: defaultPieceValue}}
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

func (mList moveList) GetMaxScoreMove() board.Move {
	if mList.size > 0 {
		maxScore := mList.scores[0]
		maxScoreMove := mList.moves[0]
		for i := 1; i < mList.size; i++ {
			if mList.scores[i] > maxScore {
				maxScore = mList.scores[i]
				maxScoreMove = mList.moves[i]
			}
		}
		return maxScoreMove
	} else {
		panic("No moves")
	}
}

func (a ai) ChooseMove(b board.Board, isWhite bool, depth int) board.Move {
	mList := newMoveList(b.GetLegalMoves(isWhite))
	if depth == 0 {
		return mList.moves[rand.Intn(mList.size)]
	} else if depth == 1 {
		for i := range mList.size {
			mList.scores[i] = a.GetScoreAfterMove(b, mList.moves[i], isWhite)
		}
		return mList.GetMaxScoreMove()
	} else {
		panic("This depth is not implemented. ")
	}

}
