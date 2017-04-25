package game

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/logrusorgru/aurora"
)

const BoardSide = 4

type Direction int

const (
	DirRight Direction = iota
	DirDown  Direction = iota
	DirLeft  Direction = iota
	DirUp    Direction = iota
)

type tileMap [BoardSide * BoardSide]int
type freezeMap [BoardSide * BoardSide]bool

type Board struct {
	tiles     tileMap
	freezeMap freezeMap
	score     int
	moves     int
}

func (b Board) get(x, y int) int {
	return b.tiles[positionToI(x, y)]
}

func (b *Board) set(x, y, v int) {
	b.tiles[positionToI(x, y)] = v
}

/* Shift pushes all tiles in the given direction until they either merge, reach
the border, or reach a tile with a different value.

If no movement occured, Shift returns false.

The algorithm is quite naive and was taken from C++ code I wrote at a 4 hours
hackathon years ago, it iterates over all cells to merge/displace them and
does this BoardSide times to ensure no gap is left between tiles.
To avoid collapsing a whole row (eg. 2 2 4 8 -> 16 instead of 0 4 4 8) each
tile that resulted from a merge is marked as "frozen" and will be skipped for
the next iterations.
*/
func (b *Board) Shift(dir Direction) bool {
	dX, dY := getShiftVector(dir)

	somethingHappened := false

	for j := 0; j < BoardSide; j++ {
		for i := 0; i < (BoardSide * BoardSide); i++ {
			if b.tiles[i] == 0 {
				continue
			}

			x, y := iToPosition(i)
			cur := b.get(x, y)
			neighX, neighY := dX+x, dY+y

			if neighX < 0 || neighX >= BoardSide ||
				neighY < 0 || neighY >= BoardSide {
				continue
			}

			neigh := b.get(neighX, neighY)

			if neigh == 0 {
				b.set(neighX, neighY, cur)
				b.set(x, y, 0)
				somethingHappened = true
			} else if cur == neigh && !b.isFrozen(x, y) {
				b.set(neighX, neighY, 2*cur)
				b.set(x, y, 0)
				b.freeze(x, y)
				b.freeze(neighX, neighY)
				b.score += 2 * cur
				somethingHappened = true
			}
		}
	}

	b.clearFreeze()
	if somethingHappened {
		b.moves += 1
	}

	return somethingHappened
}

func (b *Board) freeze(x, y int) {
	b.freezeMap[positionToI(x, y)] = true
}

func (b *Board) isFrozen(x, y int) bool {
	return b.freezeMap[positionToI(x, y)]
}

func (b *Board) clearFreeze() {
	b.freezeMap = freezeMap{}
}

// String returns a human-readable version of the Board.
func (b Board) String() string {
	str := fmt.Sprintln("┌──────┬──────┬──────┬──────┐")

	for y := 0; y < BoardSide; y++ {
		str += fmt.Sprintln("│      │      │      │      │")
		for x := 0; x < BoardSide; x++ {
			v := b.get(x, y)

			if v == 0 {
				str += fmt.Sprintf("│      ")
			} else {
				str += fmt.Sprintf(
					"│ %4d ",
					aurora.Colorize(v, getColor(v)),
				)
			}
		}
		str += fmt.Sprintln("│")
		str += fmt.Sprintln("│      │      │      │      │")

		if y < BoardSide-1 {
			str += fmt.Sprintln("├──────┼──────┼──────┼──────┤")
		}
	}

	str += fmt.Sprintln("└──────┴──────┴──────┴──────┘")

	return str
}

func getShiftVector(dir Direction) (dX, dY int) {
	switch dir {
	case DirRight:
		dX = 1
	case DirDown:
		dY = 1
	case DirLeft:
		dX = -1
	case DirUp:
		dY = -1
	}

	return dX, dY
}

func iToPosition(i int) (int, int) {
	return i % BoardSide, i / BoardSide
}

func positionToI(x int, y int) int {
	return y*BoardSide + x
}

// PlaceRandom places a 2 or a 4 (10% chance) in a random empty tile.
func (b *Board) PlaceRandom() error {
	if b.IsFull() {
		return errors.New("board is full")
	}

	available := make([]int, 0, BoardSide*BoardSide)
	for i := 0; i < (BoardSide * BoardSide); i++ {
		if b.tiles[i] == 0 {
			available = append(available, i)
		}
	}

	i := available[rand.Int()%len(available)]
	num := 2
	if rand.Intn(100) > 90 {
		num = 4
	}
	b.tiles[i] = num

	return nil
}

// IsFull returns true if the board has a value > 0 in every tile.
func (b Board) IsFull() bool {
	for i := 0; i < (BoardSide * BoardSide); i++ {
		if b.tiles[i] == 0 {
			return false
		}
	}

	return true
}

// HasMovesLeft returns true if the board can be played (ie. not in an endgame situation).
func (b Board) HasMovesLeft() bool {
	if !b.IsFull() {
		return true
	}

	if b.Shift(DirRight) {
		return true
	}
	if b.Shift(DirDown) {
		return true
	}
	if b.Shift(DirLeft) {
		return true
	}
	if b.Shift(DirUp) {
		return true
	}

	return false
}

func getColor(v int) aurora.Color {
	colors := map[int]aurora.Color{
		2:    aurora.GrayFg,
		4:    aurora.GrayFg,
		8:    aurora.BrownFg,
		16:   aurora.RedFg,
		32:   aurora.MagentaFg,
		64:   aurora.CyanFg,
		128:  aurora.GrayBg | aurora.CyanFg,
		256:  aurora.GrayBg | aurora.BlueFg,
		512:  aurora.GrayBg | aurora.MagentaFg,
		1024: aurora.GrayBg | aurora.RedFg,
		2048: aurora.GrayBg | aurora.BrownFg,
		4096: aurora.GrayBg | aurora.BlackFg,
		8192: aurora.GrayBg | aurora.GreenFg,
	}

	c, ok := colors[v]
	if !ok {
		return aurora.GrayFg
	}

	return c
}

// ColorTest fills the board with all "legal" values for testing purposes.
func (b *Board) ColorTest() {
	b.tiles = tileMap{
		8192, 4096, 2048, 1024,
		512, 256, 128, 64,
		32, 16, 8, 4,
		2, 0, 0, 0,
	}
}

func (b Board) Score() int {
	return b.score
}

func (b Board) Moves() int {
	return b.moves
}

func (b *Board) Reset() {
	b.score = 0
	b.moves = 0
	b.tiles = tileMap{}
	b.PlaceRandom()
}
