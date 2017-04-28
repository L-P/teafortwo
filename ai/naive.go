package ai

import (
	"github.com/L-P/teafortwo/game"
)

// Naive uses Up and Left as much as it can, right and down when it has to.
type Naive struct {
	board *game.Board
}

func NewNaive(board *game.Board) Naive {
	return Naive{board: board}
}

func (ai *Naive) Solve() error {
	nextIsUp := true

	for ai.board.HasMovesLeft() {
		moved := false
		if nextIsUp {
			moved, _ = ai.board.Shift(game.DirUp)
			if !moved {
				moved, _ = ai.board.Shift(game.DirLeft)
			}
		} else {
			moved, _ = ai.board.Shift(game.DirLeft)
			if !moved {
				moved, _ = ai.board.Shift(game.DirUp)
			}
		}

		if moved {
			continue
		}

		if moved, _ := ai.board.Shift(game.DirRight); !moved {
			ai.board.Shift(game.DirDown)
		}
	}

	return nil
}
