package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/L-P/teafortwo/ai"
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
	runAi := flag.Bool("ai", false, "run AI")
	flag.Parse()

	if *runAi {
		runAI()
	} else {
		runPlay()
	}
}

func runAI() {
	rand.Seed(42)

	var highest = game.Board{}
	board := game.Board{}

	for i := 0; i < 5000; i++ {
		board.Reset()

		ai := ai.NewHungry(&board)
		if err := ai.Solve(); err != nil {
			panic(err)
		}

		if board.Highest() > highest.Highest() {
			highest = board
		}
	}

	fmt.Println(highest.String())
	fmt.Printf("score: %d\n", highest.Score())
	fmt.Printf("moves: %d\n", highest.Moves())
	fmt.Printf("largest tile: %d\n", highest.Highest())

}

func runPlay() {
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

const boardView = "board"
const scoreView = "score"
const messageView = "message"

func layout(s *GameState) func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
		if v, err := g.SetView(boardView, 0, 0, 30, 18); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			g.SetCurrentView(boardView)
			v.Frame = false
		}

		if _, err := g.SetView(scoreView, 32, 1, 46, 4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}

		if v, err := g.SetView(messageView, 0, 18, 30, 21); err != nil {
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

		{boardView, gocui.KeyArrowRight, makeShiftCallback(s, game.DirRight)},
		{boardView, gocui.KeyArrowDown, makeShiftCallback(s, game.DirDown)},
		{boardView, gocui.KeyArrowLeft, makeShiftCallback(s, game.DirLeft)},
		{boardView, gocui.KeyArrowUp, makeShiftCallback(s, game.DirUp)},

		// It would not be a real CLI game otherwise.
		{boardView, 'l', makeShiftCallback(s, game.DirRight)},
		{boardView, 'j', makeShiftCallback(s, game.DirDown)},
		{boardView, 'h', makeShiftCallback(s, game.DirLeft)},
		{boardView, 'k', makeShiftCallback(s, game.DirUp)},
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
		shifted, err := s.board.Shift(dir)
		if err != nil {
			return err
		}

		if !shifted {
			return nil
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
	board, err := g.View(boardView)
	if err != nil {
		return fmt.Errorf("unable to get board view: %s", err)
	}
	board.Clear()
	fmt.Fprintln(board, s.board.String())

	score, err := g.View(scoreView)
	if err != nil {
		return fmt.Errorf("unable to get score view: %s", err)
	}
	score.Clear()
	fmt.Fprintf(score, "score: %6d\n", s.board.Score())
	fmt.Fprintf(score, "moves: %6d\n", s.board.Moves())

	message, err := g.View(messageView)
	if err != nil {
		return fmt.Errorf("unable to get message view: %s", err)
	}
	message.Clear()
	fmt.Fprintf(message, s.message)

	return nil
}
