package main

import (
	"fmt"
	"strings"
)

const BoardSide = 4

type Direction int

const (
	DirRight Direction = iota
	DirDown  Direction = iota
	DirLeft  Direction = iota
	DirUp    Direction = iota
)

type Board struct {
	raw    [BoardSide * BoardSide]int
	freeze [BoardSide * BoardSide]bool
}

func (b Board) Get(x, y int) int {
	return b.raw[y*BoardSide+x]
}

func (b *Board) Set(x, y, v int) {
	b.raw[y*BoardSide+x] = v
}

func (b *Board) Freeze(x, y int) {
	b.freeze[y*BoardSide+x] = true
}

func (b *Board) IsFrozen(x, y int) bool {
	return b.freeze[y*BoardSide+x]
}

func (b *Board) ClearFreeze() {
	b.freeze = [BoardSide * BoardSide]bool{}
}

// Shift returns true if any movement occured
func (b *Board) Shift(dir Direction) bool {
	dX, dY := getShiftVector(dir)

	somethingHappened := false

	for j := 0; j < BoardSide; j++ {
		for i := 0; i < (BoardSide * BoardSide); i++ {
			if b.raw[i] == 0 {
				continue
			}

			x, y := iToPosition(i)
			cur := b.Get(x, y)
			neighX, neighY := dX+x, dY+y

			if neighX < 0 || neighX > BoardSide ||
				neighY < 0 || neighY > BoardSide {
				continue
			}

			neigh := b.Get(neighX, neighY)

			if neigh == 0 {
				b.Set(neighX, neighY, cur)
				b.Set(x, y, 0)
				somethingHappened = true
			} else if cur == neigh && !b.IsFrozen(x, y) {
				b.Set(neighX, neighY, 2*cur)
				b.Set(x, y, 0)
				b.Freeze(x, y)
				somethingHappened = true
			}
		}
	}

	b.ClearFreeze()

	return somethingHappened
}

func iToPosition(i int) (int, int) {
	return i % BoardSide, i / BoardSide
}

func (b *Board) Collate(dir Direction) {
	for y := 0; y < BoardSide; y++ {
		for x := 0; x < BoardSide; x++ {
		}
	}
}

func (b Board) String() string {
	pad := func(v int) string {
		num := fmt.Sprintf("%d", v)
		count := 4 - len(num)
		if count <= 0 {
			return num
		}

		return strings.Repeat(" ", count) + num
	}

	for y := 0; y < BoardSide; y++ {
		for x := 0; x < BoardSide; x++ {
			fmt.Printf(" %s ", pad(b.Get(x, y)))
		}

		fmt.Print("\n")
	}

	return ""
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
