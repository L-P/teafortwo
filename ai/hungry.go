package ai

import "github.com/L-P/teafortwo/game"

// Hungry tries to maximize score with each move
type Hungry struct {
	board *game.Board
}

// NewHungry creates a new Hungry AI
func NewHungry(board *game.Board) Hungry {
	return Hungry{board: board}
}

// Solve attempts to get the highest score and win the game.
func (ai *Hungry) Solve() error {
	for ai.board.HasMovesLeft() {
		var bestDirection game.Direction
		var available game.Direction
		bestScore := 0
		hasBest := false

		for _, dir := range []game.Direction{game.DirRight, game.DirDown, game.DirLeft, game.DirUp} {
			board := *ai.board
			ok, err := board.Shift(dir)
			if err != nil {
				return err
			}

			if ok {
				available = dir
				if board.Score() > bestScore {
					bestDirection = dir
					hasBest = false
				}
			}
		}

		// If no best score found, just pick a valid direction.
		if !hasBest {
			bestDirection = available
		}

		if _, err := ai.board.Shift(bestDirection); err != nil {
			return err
		}
	}

	return nil
}
