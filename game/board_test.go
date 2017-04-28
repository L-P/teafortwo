package game

import (
	"reflect"
	"testing"
)

func TestSparse(t *testing.T) {
	b := Board{}
	b.set(1, 0, 2)
	b.set(3, 0, 2)
	b.Shift(DirLeft)

	if b.Get(0, 0) != 4 {
		t.Error("tiles were not merged")
	}

	for i := 1; i < BoardSide; i++ {
		if b.Get(i, 0) != 0 {
			t.Error("tiles were not reset to 0")
		}
	}
}

func TestNoRemerge(t *testing.T) {
	b := Board{}
	b.set(0, 0, 2)
	b.set(0, 1, 2)
	b.set(0, 2, 2)
	b.set(0, 3, 2)
	b.Shift(DirUp)

	if b.Get(0, 0) != 4 || b.Get(0, 1) != 4 {
		t.Error("tiles were not merged")
	}

	for i := 2; i < BoardSide; i++ {
		if b.Get(0, i) != 0 {
			t.Error("tiles were not reset to 0")
		}
	}

	b.Shift(DirUp)
	if b.Get(0, 0) != 8 {
		t.Error("tiles were not merged")
	}
}

func TestShift(t *testing.T) {
	for i, v := range getCases() {
		b := Board{tiles: v.before}
		b.Shift(v.direction)
		if !reflect.DeepEqual(b.tiles, v.expected) {
			t.Errorf("test case #%d failed", i)
		}
	}
}

type testCase struct {
	before    TileMap
	direction Direction
	expected  TileMap
}

func getCases() []testCase {
	return []testCase{
		testCase{
			before: TileMap{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
			direction: DirUp,
			expected: TileMap{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
		testCase{
			before: TileMap{
				2, 2, 2, 2,
				4, 0, 4, 0,
				0, 8, 0, 8,
				2, 4, 8, 16,
			},
			direction: DirLeft,
			expected: TileMap{
				4, 4, 0, 0,
				8, 0, 0, 0,
				16, 0, 0, 0,
				2, 4, 8, 16,
			},
		},
		testCase{
			before: TileMap{
				2, 2, 2, 2,
				4, 0, 4, 0,
				0, 8, 0, 8,
				2, 4, 8, 16,
			},
			direction: DirDown,
			expected: TileMap{
				0, 0, 0, 0,
				2, 2, 2, 2,
				4, 8, 4, 8,
				2, 4, 8, 16,
			},
		},
		testCase{
			before: TileMap{
				2, 2, 2, 2,
				2, 2, 2, 2,
				2, 2, 2, 2,
				2, 2, 2, 2,
			},
			direction: DirRight,
			expected: TileMap{
				0, 0, 4, 4,
				0, 0, 4, 4,
				0, 0, 4, 4,
				0, 0, 4, 4,
			},
		},
		testCase{
			before: TileMap{
				8, 0, 0, 0,
				4, 0, 0, 0,
				4, 0, 0, 0,
				0, 0, 0, 0,
			},
			direction: DirUp,
			expected: TileMap{
				8, 0, 0, 0,
				8, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
		testCase{
			before: TileMap{
				16, 0, 0, 0,
				8, 0, 0, 0,
				4, 0, 0, 0,
				4, 0, 0, 0,
			},
			direction: DirUp,
			expected: TileMap{
				16, 0, 0, 0,
				8, 0, 0, 0,
				8, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
	}
}
