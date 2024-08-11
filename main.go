package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charlie-pecora/bubbles/levels"
	tea "github.com/charmbracelet/bubbletea"
)

const updateTimeMs = 10
const nextLevelTimeMs = 1000

func ticker() tea.Msg {
	time.Sleep(time.Millisecond * updateTimeMs)
	return levels.TickMessage
}

func resetTicker() tea.Msg {
	time.Sleep(time.Millisecond * nextLevelTimeMs)
	return levels.ResetMessage
}

type model struct {
	userLocation   levels.Coordinates
	TargetLocation levels.Coordinates
	walls          []levels.Coordinates
	userMove       levels.MoveValue
	level          int
	gameWon        bool
}

func initialModel() model {
	return model{
		userLocation: levels.Coordinates{X: 5, Y: 5},
		userMove:     levels.NilKey,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return ticker
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch s := levels.MoveValue(msg.String()); s {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case levels.UpKey, levels.DownKey, levels.LeftKey, levels.RightKey:
			m.userMove = s
		}

	case levels.ReadyMessageType:
		switch msg {
		case levels.TickMessage:
			// make a move and restart the ticker
			m.HandleMove()
			if m.LevelCompleted() {
				return m, resetTicker
			} else {
				return m, ticker
			}
		case levels.ResetMessage:
			m.NextLevel()
			return m, ticker
		}
	}
	return m, nil
}

func (m *model) HandleMove() {
	switch m.userMove {
	case levels.UpKey:
		if m.userLocation.Y > 0 {
			m.userLocation.Y -= 1
		}
	case levels.DownKey:
		if m.userLocation.Y < levels.GridSizeY-1 {
			m.userLocation.Y += 1
		}
	case levels.LeftKey:
		if m.userLocation.X > 0 {
			m.userLocation.X -= 1
		}
	case levels.RightKey:
		if m.userLocation.X < levels.GridSizeX-1 {
			m.userLocation.X += 1
		}
	}
	m.userMove = levels.NilKey
}

func (m model) LevelCompleted() bool {
	return m.userLocation == m.TargetLocation
}

func (m *model) NextLevel() {
	m.level += 1
	if m.level >= len(levels.LevelsArray) {
		m.gameWon = true
	} else {
		nextLevel := levels.LevelsArray[m.level]
		m.TargetLocation = nextLevel.TargetLocation
		m.walls = nextLevel.Walls
	}
}

func (m model) View() string {
	// The header
	s := "Grug Game\n\n"
	if m.gameWon {
		s += "You Win!!!!"
	} else if m.LevelCompleted() {
		s += fmt.Sprintf("You beat level %v!", m.level)
	} else {

		for yi := 0; yi < levels.GridSizeY; yi++ {
			rowString := make([]byte, 0, levels.GridSizeX)
			for xi := 0; xi < levels.GridSizeX; xi++ {
				switch (levels.Coordinates{X: xi, Y: yi}) {
				case m.userLocation:
					rowString = append(rowString, levels.UserChar)
				case m.TargetLocation:
					rowString = append(rowString, levels.TargetChar)
				default:
					rowString = append(rowString, levels.EmptyChar)
				}
			}
			s += string(rowString) + "\n"
		}
	}
		// The footer
	s += "\\nnPress q or ctrl+c to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
