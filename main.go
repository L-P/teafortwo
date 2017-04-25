package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"

	"home.leo-peltier.fr/teafortwo/game"
)

func main() {
	board := &game.Board{}
	board.PlaceRandom()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("unable to init term: %s", err)
	}
	defer g.Close()

	g.SetManagerFunc(layout(board))

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("board", gocui.KeyArrowRight, gocui.ModNone, makeShiftCallback(board, game.DirRight)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("board", gocui.KeyArrowDown, gocui.ModNone, makeShiftCallback(board, game.DirDown)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("board", gocui.KeyArrowLeft, gocui.ModNone, makeShiftCallback(board, game.DirLeft)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("board", gocui.KeyArrowUp, gocui.ModNone, makeShiftCallback(board, game.DirUp)); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func layout(b *game.Board) func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
		if v, err := g.SetView("board", 0, 0, 30, 18); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			g.SetCurrentView("board")
			v.Frame = false
			fmt.Fprintln(v, b.String())
		}
		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func makeShiftCallback(b *game.Board, dir game.Direction) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if !b.Shift(dir) {
			return nil
		}

		if err := b.PlaceRandom(); err != nil {
			return err
		}

		v.Clear()
		fmt.Fprintln(v, b.String())
		return nil
	}
}
