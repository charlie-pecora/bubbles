package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/charlie-pecora/bubbles/levels"
	tea "github.com/charmbracelet/bubbletea"
)

const updateTimeMs = 10
const ticketPerLevelUpdate = 50
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
	userLocation  levels.Coordinates
	level         levels.Level
	userMove      levels.MoveValue
	gameWon       bool
	gameLost      bool
	displayWidth  int
	displayHeight int
	counter       int
}

func initialModel() model {
	level, _ := levels.GetLevel(0)
	return model{
		userLocation: levels.Coordinates{X: 5, Y: 5},
		userMove:     levels.NilKey,
		level:        level,
	}
}

func (m model) Init() tea.Cmd {
	return ticker
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.displayHeight = msg.Height
		m.displayWidth = msg.Width

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch s := levels.MoveValue(msg.String()); s {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// Restart the game if you have won
		case "enter":
			if m.gameLost || m.gameWon {
				m = initialModel()
				return m, ticker
			}

		case levels.UpKey, levels.DownKey, levels.LeftKey, levels.RightKey:
			m.userMove = s
		}

	case levels.ReadyMessageType:
		switch msg {
		case levels.TickMessage:
			// make a move and restart the ticker
			m.HandleMove()
			m.counter += 1
			if m.counter%ticketPerLevelUpdate == 0 {
				m.counter = 0
				m.level.Update()
			}
			m.level.Update()
			if m.CheckLost() {
				m.gameLost = true
				return m, nil
			} else if m.LevelCompleted() {
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

func (m model) CheckLost() bool {
	return slices.Contains(m.level.Enemies, m.userLocation)
}

func (m model) LevelCompleted() bool {
	return m.userLocation == m.level.TargetLocation
}

func (m *model) NextLevel() {
	nextLevel, ok := levels.GetLevel(m.level.Number + 1)
	if !ok {
		m.gameWon = true
	} else {
		m.level = nextLevel
	}
}

func (m model) View() string {
	// The header
	s := "Grug Game\n\n"
	s += fmt.Sprintf("(%v, %v)\n", m.displayHeight, m.displayWidth)
	if m.gameLost {
		s += "You Lose!!!!\n\nPress enter to play again!"
	} else if m.gameWon {
		s += "You Win!!!!\n\nPress enter to play again!"
	} else if m.LevelCompleted() {
		s += fmt.Sprintf("You beat level %v!", m.level.Number)
	} else {

		for yi := 0; yi < levels.GridSizeY; yi++ {
			rowString := make([]byte, 0, levels.GridSizeX)
			for xi := 0; xi < levels.GridSizeX; xi++ {
				current := levels.Coordinates{X: xi, Y: yi}
				if slices.Contains(m.level.Enemies, current) {
					rowString = append(rowString, levels.EnemyChar)
				} else if current == m.userLocation {
					rowString = append(rowString, levels.UserChar)
				} else if current == m.level.TargetLocation {
					rowString = append(rowString, levels.TargetChar)
				} else {
					rowString = append(rowString, levels.EmptyChar)
				}
			}
			s += string(rowString) + "\n"
		}
	}
	// The footer
	s += "\n\nPress q or ctrl+c to quit.\n"

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
