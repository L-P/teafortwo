package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/L-P/teafortwo/game"
	"github.com/jroimartin/gocui"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	board := &game.Board{}
	board.Reset()

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("unable to init term: %s", err)
	}
	defer gui.Close()

	gui.SetManagerFunc(layout(board))
	if err := setBinds(board, gui); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

const BoardView = "board"
const ScoreView = "score"

func layout(b *game.Board) func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
		if v, err := g.SetView(BoardView, 0, 0, 30, 18); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			g.SetCurrentView(BoardView)
			v.Frame = false
		}

		if _, err := g.SetView(ScoreView, 32, 1, 46, 4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}

		redraw(b, g)
		return nil
	}
}

func setBinds(b *game.Board, g *gocui.Gui) error {
	binds := []struct {
		view string
		key  interface{}
		fn   func(g *gocui.Gui, v *gocui.View) error
	}{
		{"", gocui.KeyCtrlC, quit},
		{"", 'q', quit},
		{"", 'n', makeResetCallback(b)},
		{BoardView, gocui.KeyArrowRight, makeShiftCallback(b, game.DirRight)},
		{BoardView, gocui.KeyArrowDown, makeShiftCallback(b, game.DirDown)},
		{BoardView, gocui.KeyArrowLeft, makeShiftCallback(b, game.DirLeft)},
		{BoardView, gocui.KeyArrowUp, makeShiftCallback(b, game.DirUp)},
	}

	for _, v := range binds {
		if err := g.SetKeybinding(v.view, v.key, gocui.ModNone, v.fn); err != nil {
			return err
		}
	}

	return nil
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

		// TODO: actually handle the endgame, panic'ing is not exactly user-friendly.
		if !b.HasMovesLeft() {
			return fmt.Errorf("no moves left, score: %d", b.Score())
		}

		return redraw(b, g)
	}
}

func makeResetCallback(b *game.Board) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		b.Reset()
		return redraw(b, g)
	}
}

func redraw(b *game.Board, g *gocui.Gui) error {
	board, err := g.View(BoardView)
	if err != nil {
		return fmt.Errorf("unable to get board view: %s", err)
	}
	board.Clear()
	fmt.Fprintln(board, b.String())

	score, err := g.View(ScoreView)
	if err != nil {
		return fmt.Errorf("unable to get score view: %s", err)
	}
	score.Clear()
	fmt.Fprintf(score, "score: %6d\n", b.Score())
	fmt.Fprintf(score, "moves: %6d\n", b.Moves())

	return nil
}
