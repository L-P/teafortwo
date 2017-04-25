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

	if err := g.SetKeybinding("", 'n', gocui.ModNone, makeResetCallback(board)); err != nil {
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
