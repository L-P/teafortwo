package main

import (
	"fmt"

	"home.leo-peltier.fr/teafortwo/game"
)

func main() {
	board := game.Board{}
	board.Set(0, 0, 512)
	board.Set(1, 0, 16)
	board.Set(2, 0, 16)
	board.Set(3, 0, 32)
	board.Set(2, 2, 8)
	board.Set(3, 1, 512)

	fmt.Println(board.String())

	board.Shift(game.DirLeft)
	fmt.Println(board.String())

	board.Shift(game.DirUp)
	fmt.Println(board.String())

	board.Shift(game.DirLeft)
	fmt.Println(board.String())
}
