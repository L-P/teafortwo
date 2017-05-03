package game

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/logrusorgru/aurora"
)

// BoardSide is the board edge size in tiles.
const BoardSide = 4

// Direction represents a cardinal direction.
type Direction int

const (
	// DirRight goes east.
	DirRight Direction = iota

	// DirDown goes south.
	DirDown Direction = iota

	// DirLeft goes west.
	DirLeft Direction = iota

	// DirUp goes north.
	DirUp Direction = iota
)

// Directions returns all possible shift directions
func Directions() []Direction {
	return []Direction{
		DirRight,
		DirDown,
		DirLeft,
		DirUp,
	}
}

// TileMap is an array of tile values.
type TileMap [BoardSide * BoardSide]int
type freezeMap [BoardSide * BoardSide]bool

// Board represents the game board, its tiles, and score.
type Board struct {
	tiles     TileMap
	freezeMap freezeMap
	score     int
	moves     int
	highest   int
	won       bool
}

// Get returns the value of the tile at the given position.
func (b Board) Get(x, y int) int {
	return b.tiles[positionToI(x, y)]
}

// GetTiles returns a copy of the internal tiles array.
func (b Board) GetTiles() TileMap {
	return b.tiles
}

func (b *Board) set(x, y, v int) {
	b.tiles[positionToI(x, y)] = v
}

// Shift pushes all tiles in the given direction until they either merge, reach
// the border, or reach a tile with a different value.
//
// If no movement occured, Shift returns an ImpossibleShift error.
// If something else and very wrong happened, a generic error is returned.
//
// A random tile is placed after shifting.
func (b *Board) Shift(dir Direction) error {
	err := b.doShift(dir)
	if err != nil {
		return err
	}

	b.moves++

	if err := b.placeRandom(); err != nil {
		return err
	}

	return nil
}

// doShift does the actual shifting so we can test it without the placeRandom part
//
// The algorithm is quite naive and was taken from C++ code I wrote at a 4 hours
// hackathon years ago, it iterates over all cells to merge/displace them and
// does this BoardSide times to ensure no gap is left between tiles.
// To avoid collapsing a whole row (eg. 2 2 4 8 -> 16 instead of 0 4 4 8) each
// tile that resulted from a merge is marked as "frozen" and will be skipped for
// the next iterations.
func (b *Board) doShift(dir Direction) error {
	dX, dY := getShiftVector(dir)
	defer b.clearFreeze()
	somethingHappened := false

	for j := 0; j < BoardSide; j++ {
		for i := 0; i < (BoardSide * BoardSide); i++ {
			if b.tiles[i] == 0 {
				continue
			}

			x, y := iToPosition(i)
			neighX, neighY := x+dX, y+dY

			if neighX < 0 || neighX >= BoardSide ||
				neighY < 0 || neighY >= BoardSide {
				continue
			}

			neigh := b.Get(neighX, neighY)

			if neigh == 0 {
				// No merging, move tiles around
				b.set(neighX, neighY, b.tiles[i])
				b.set(x, y, 0)
				somethingHappened = true
			} else if b.tiles[i] == neigh && !b.isFrozen(x, y) {
				// Merge adjacent identical values if they did not result from a merge this Shift call.
				new := 2 * b.tiles[i]

				// Place values
				b.set(neighX, neighY, new)
				b.set(x, y, 0)

				// Disable merging these tiles until next Shift call
				b.freeze(x, y)
				b.freeze(neighX, neighY)

				b.updateScore(new)
				somethingHappened = true
			}
		}
	}

	if !somethingHappened {
		return &ImpossibleShift{dir}
	}

	return nil
}

// ImpossibleShift is returned by Shift() when a shift can't be performed in a given direction
type ImpossibleShift struct {
	dir Direction
}

func (e *ImpossibleShift) Error() string {
	return fmt.Sprintf("can't shift %s", DirectionToName(e.dir))
}

// DirectionToName returns the direction name from its value.
func DirectionToName(dir Direction) string {
	return map[Direction]string{
		DirRight: "right",
		DirDown:  "down",
		DirLeft:  "left",
		DirUp:    "up",
	}[dir]
}

func (b *Board) updateScore(tileValue int) {
	if tileValue > b.highest {
		b.highest = tileValue
	}
	if tileValue >= 2048 {
		b.won = true
	}
	b.score += tileValue
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
			v := b.Get(x, y)

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

func iToPosition(i int) (x int, y int) {
	return i % BoardSide, i / BoardSide
}

func positionToI(x int, y int) int {
	return y*BoardSide + x
}

// placeRandom places a 2 or a 4 (10% chance) in a random empty tile.
func (b *Board) placeRandom() error {
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

	for _, v := range Directions() {
		if err := b.doShift(v); err == nil {
			return true
		}
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
	b.tiles = TileMap{
		8192, 4096, 2048, 1024,
		512, 256, 128, 64,
		32, 16, 8, 4,
		2, 0, 0, 0,
	}
}

// Score returns the current game score.
// The score is the sum of all merged values.
func (b Board) Score() int {
	return b.score
}

// Moves returns the number of successful Shift done on the board.
func (b Board) Moves() int {
	return b.moves
}

// Reset resets the board to its initial state (no score, only one random tile).
func (b *Board) Reset() {
	*b = Board{}
	b.placeRandom()
}

// Won returns true if board holds a winning game, meaning the player
// reached 2048 (it's the name of the game).
func (b Board) Won() bool {
	return b.won
}

// Highest returns the highest value on the board.
func (b Board) Highest() int {
	return b.highest
}
