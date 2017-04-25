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

	if b.get(0, 0) != 4 {
		t.Error("tiles were not merged")
	}

	for i := 1; i < BoardSide; i++ {
		if b.get(i, 0) != 0 {
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

	if b.get(0, 0) != 4 || b.get(0, 1) != 4 {
		t.Error("tiles were not merged")
	}

	for i := 2; i < BoardSide; i++ {
		if b.get(0, i) != 0 {
			t.Error("tiles were not reset to 0")
		}
	}

	b.Shift(DirUp)
	if b.get(0, 0) != 8 {
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
	before    tileMap
	direction Direction
	expected  tileMap
}

func getCases() []testCase {
	return []testCase{
		testCase{
			before: tileMap{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
			direction: DirUp,
			expected: tileMap{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
		testCase{
			before: tileMap{
				2, 2, 2, 2,
				4, 0, 4, 0,
				0, 8, 0, 8,
				2, 4, 8, 16,
			},
			direction: DirLeft,
			expected: tileMap{
				4, 4, 0, 0,
				8, 0, 0, 0,
				16, 0, 0, 0,
				2, 4, 8, 16,
			},
		},
		testCase{
			before: tileMap{
				2, 2, 2, 2,
				4, 0, 4, 0,
				0, 8, 0, 8,
				2, 4, 8, 16,
			},
			direction: DirDown,
			expected: tileMap{
				0, 0, 0, 0,
				2, 2, 2, 2,
				4, 8, 4, 8,
				2, 4, 8, 16,
			},
		},
		testCase{
			before: tileMap{
				2, 2, 2, 2,
				2, 2, 2, 2,
				2, 2, 2, 2,
				2, 2, 2, 2,
			},
			direction: DirRight,
			expected: tileMap{
				0, 0, 4, 4,
				0, 0, 4, 4,
				0, 0, 4, 4,
				0, 0, 4, 4,
			},
		},
		testCase{
			before: tileMap{
				8, 0, 0, 0,
				4, 0, 0, 0,
				4, 0, 0, 0,
				0, 0, 0, 0,
			},
			direction: DirUp,
			expected: tileMap{
				8, 0, 0, 0,
				8, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
		testCase{
			before: tileMap{
				16, 0, 0, 0,
				8, 0, 0, 0,
				4, 0, 0, 0,
				4, 0, 0, 0,
			},
			direction: DirUp,
			expected: tileMap{
				16, 0, 0, 0,
				8, 0, 0, 0,
				8, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
	}
}
