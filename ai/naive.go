package ai

import (
	"errors"

	"github.com/L-P/teafortwo/game"
)

// Naive uses Up and Left as much as it can, right and down when it has to.
type Naive struct {
	board    *game.Board
	nextIsUp bool
}

// NewNaive creates a new Naive AI
func NewNaive(board *game.Board) Naive {
	return Naive{
		board:    board,
		nextIsUp: true,
	}
}

// Solve attempts to get the highest score and win the game.
func (ai *Naive) Solve() error {
	for ai.board.HasMovesLeft() {
		best, err := ai.FindBest(*ai.board)
		if err != nil {
			return err
		}

		if ai.nextIsUp && best == game.DirUp {
			ai.nextIsUp = false
		} else if !ai.nextIsUp && best == game.DirLeft {
			ai.nextIsUp = true
		}

		if err := ai.board.Shift(best); err != nil {
			return err
		}
	}

	return nil
}

// FindBest TODO
func (ai Naive) FindBest(board game.Board) (game.Direction, error) {
	if ai.nextIsUp {
		if ok, _ := ai.board.CanShift(game.DirUp); ok {
			return game.DirUp, nil
		}
		if ok, _ := ai.board.CanShift(game.DirLeft); ok {
			return game.DirLeft, nil
		}
	} else {
		if ok, _ := ai.board.CanShift(game.DirLeft); ok {
			return game.DirLeft, nil
		}
		if ok, _ := ai.board.CanShift(game.DirUp); ok {
			return game.DirUp, nil
		}
	}

	if ok, _ := ai.board.CanShift(game.DirDown); ok {
		return game.DirDown, nil
	}
	if ok, _ := ai.board.CanShift(game.DirRight); ok {
		return game.DirRight, nil
	}

	return game.DirNone, errors.New("no possible direction")
}
