package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/L-P/teafortwo/game"
	"github.com/jroimartin/gocui"
	"github.com/logrusorgru/aurora"
)

// GameState holds the global state of a game.
type GameState struct {
	board   game.Board
	message string
}

const howtoMessage = "q: quit, n: new game\n→↓←↑: move tiles around"

func main() {
	rand.Seed(time.Now().UnixNano())
	state := &GameState{message: howtoMessage}
	state.board.Reset()

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("unable to init term: %s", err)
	}
	defer gui.Close()

	gui.SetManagerFunc(layout(state))
	if err := setBinds(state, gui); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

const BoardView = "board"
const ScoreView = "score"
const MessageView = "message"

func layout(s *GameState) func(*gocui.Gui) error {
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

		if v, err := g.SetView(MessageView, 0, 18, 30, 21); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Frame = false
		}

		redraw(s, g)
		return nil
	}
}

func setBinds(s *GameState, g *gocui.Gui) error {
	binds := []struct {
		view string
		key  interface{}
		fn   func(g *gocui.Gui, v *gocui.View) error
	}{
		{"", gocui.KeyCtrlC, quit},
		{"", 'q', quit},
		{"", 'n', makeResetCallback(s)},

		{BoardView, gocui.KeyArrowRight, makeShiftCallback(s, game.DirRight)},
		{BoardView, gocui.KeyArrowDown, makeShiftCallback(s, game.DirDown)},
		{BoardView, gocui.KeyArrowLeft, makeShiftCallback(s, game.DirLeft)},
		{BoardView, gocui.KeyArrowUp, makeShiftCallback(s, game.DirUp)},

		// It would not be a real CLI game otherwise.
		{BoardView, 'l', makeShiftCallback(s, game.DirRight)},
		{BoardView, 'j', makeShiftCallback(s, game.DirDown)},
		{BoardView, 'h', makeShiftCallback(s, game.DirLeft)},
		{BoardView, 'k', makeShiftCallback(s, game.DirUp)},
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

func makeShiftCallback(s *GameState, dir game.Direction) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if !s.board.Shift(dir) {
			return nil
		}

		if err := s.board.PlaceRandom(); err != nil {
			return err
		}

		if s.board.Won() {
			s.message = aurora.Green(
				"You won! You can keep playing\nor reset the game with 'n'.",
			).String()
		}

		if !s.board.HasMovesLeft() {
			s.message = aurora.Red(
				"No moves left.\nPress 'n' to start a new game.",
			).String()
		}

		return redraw(s, g)
	}
}

func makeResetCallback(s *GameState) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		s.message = howtoMessage
		s.board.Reset()
		return redraw(s, g)
	}
}

func redraw(s *GameState, g *gocui.Gui) error {
	board, err := g.View(BoardView)
	if err != nil {
		return fmt.Errorf("unable to get board view: %s", err)
	}
	board.Clear()
	fmt.Fprintln(board, s.board.String())

	score, err := g.View(ScoreView)
	if err != nil {
		return fmt.Errorf("unable to get score view: %s", err)
	}
	score.Clear()
	fmt.Fprintf(score, "score: %6d\n", s.board.Score())
	fmt.Fprintf(score, "moves: %6d\n", s.board.Moves())

	message, err := g.View(MessageView)
	if err != nil {
		return fmt.Errorf("unable to get message view: %s", err)
	}
	message.Clear()
	fmt.Fprintf(message, s.message)

	return nil
}
