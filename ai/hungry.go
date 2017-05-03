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
		best, err := findBest(*ai.board)
		if err != nil {
			return err
		}

		if _, err := ai.board.Shift(best); err != nil {
			return err
		}
	}

	return nil
}

func findBest(board game.Board) (game.Direction, error) {
	var bestDirection game.Direction
	var available game.Direction
	bestScore := 0
	hasBest := false

	for _, dir := range []game.Direction{game.DirRight, game.DirDown, game.DirLeft, game.DirUp} {
		board := board
		ok, err := board.Shift(dir)
		if err != nil {
			return game.DirNone, err
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

	return bestDirection, nil
}
