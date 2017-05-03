package game

import (
	"reflect"
	"testing"
	"time"
)

func TestSparse(t *testing.T) {
	b := Board{}
	b.set(1, 0, 2)
	b.set(3, 0, 2)
	err := b.doShift(DirLeft)
	if err != nil {
		t.Fatal(err)
	}

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
	b.doShift(DirUp)

	if b.Get(0, 0) != 4 || b.Get(0, 1) != 4 {
		t.Error("tiles were not merged")
	}

	for i := 2; i < BoardSide; i++ {
		if b.Get(0, i) != 0 {
			t.Error("tiles were not reset to 0")
		}
	}

	b.doShift(DirUp)
	if b.Get(0, 0) != 8 {
		t.Error("tiles were not merged")
	}
}

func TestShift(t *testing.T) {
	for i, v := range getTestShiftCases() {
		b := Board{tiles: v.before}
		b.doShift(v.direction)
		if !reflect.DeepEqual(b.tiles, v.expected) {
			t.Errorf("test case #%d failed", i)
		}
	}
}

type testShiftCase struct {
	before    TileMap
	direction Direction
	expected  TileMap
}

func getTestShiftCases() []testShiftCase {
	return []testShiftCase{
		testShiftCase{
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
		testShiftCase{
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
		testShiftCase{
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
		testShiftCase{
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
		testShiftCase{
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
		testShiftCase{
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
		testShiftCase{
			before: TileMap{
				16, 0, 0, 0,
				0, 0, 2, 0,
				8, 0, 0, 2,
				8, 0, 2, 0,
			},
			direction: DirUp,
			expected: TileMap{
				16, 0, 4, 2,
				16, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
	}
}

type testPositionToICase struct {
	x        int
	y        int
	expected int
}

func getTestPositionToICases() []testPositionToICase {
	return []testPositionToICase{
		testPositionToICase{
			x:        0,
			y:        0,
			expected: 0,
		},
		testPositionToICase{
			x:        1,
			y:        0,
			expected: 1,
		},
		testPositionToICase{
			x:        0,
			y:        1,
			expected: 4,
		},
		testPositionToICase{
			x:        3,
			y:        3,
			expected: 15,
		},
	}
}

func TestPositionToICases(t *testing.T) {
	for _, v := range getTestPositionToICases() {
		actual := positionToI(v.x, v.y)
		if actual != v.expected {
			t.Errorf(
				"positionToI(%d, %d) = %d; expected %d",
				v.x,
				v.y,
				actual,
				v.expected,
			)
		}

		x, y := iToPosition(actual)
		if x != v.x || y != v.y {
			t.Errorf(
				"iToPosition(%d) = %d, %d; expected %d, %d",
				actual,
				x,
				y,
				v.x,
				v.y,
			)
		}
	}
}

func TestRandIsDeterministic(t *testing.T) {
	for i := 0; i < 100; i++ {
		now := time.Now().UnixNano()
		b1 := NewBoard(now)
		b2 := NewBoard(now)

		b1.placeRandom()
		b2.placeRandom()
		b1.placeRandom()
		b2.placeRandom()
		b1.placeRandom()
		b2.placeRandom()

		for j := 0; j < BoardSide*BoardSide; j++ {
			if b1.tiles[j] != b2.tiles[j] {
				t.Fatal("boards did not have the same randomness")
			}
		}
	}
}
