package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const gridSizeX int = 100
const gridSizeY int = 30
const userChar byte = 'A'
const backgroundChar byte = '.'
const updateTimems = 100

const (
	nilKey   string = ""
	upKey    string = "k"
	downKey  string = "j"
	leftKey  string = "h"
	rightKey string = "l"
)

type model struct {
	grid         [gridSizeY][gridSizeX]byte
	userLocation [2]int
	latestMove   string
}

func initialModel() model {
	var grid [gridSizeY][gridSizeX]byte
	for i := 0; i < gridSizeY; i++ {
		for j := 0; j < gridSizeX; j++ {
			grid[i][j] = backgroundChar
		}
	}
	return model{
		grid:         grid,
		userLocation: [2]int{5, 5},
		latestMove:   nilKey,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return ticker
}

type readyMessageType string

const readyMessage readyMessageType = "Ready"

func ticker() tea.Msg {
	time.Sleep(time.Millisecond * updateTimems)
	return readyMessage
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		// capture movement key presses
		case upKey, downKey, leftKey, rightKey:
			m.latestMove = msg.String()
		}
	case readyMessageType:
		switch m.latestMove {
		case upKey:
			if m.userLocation[0] > 0 {
				m.userLocation[0] -= 1
			}
		case downKey:
			if m.userLocation[0] < gridSizeY-1 {
				m.userLocation[0] += 1
			}
		case leftKey:
			if m.userLocation[1] > 0 {
				m.userLocation[1] -= 1
			}
		case rightKey:
			if m.userLocation[1] < gridSizeX-1 {
				m.userLocation[1] += 1
			}
		}
		m.latestMove = nilKey
		return m, ticker
	}
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Grug Game\n\n"

	for i, row := range m.grid {
		rowString := row[:]
		if m.userLocation[0] == i {
			rowString[m.userLocation[1]] = userChar
		}
		s += string(rowString) + "\n"
	}

	// The footer
	s += "\nPress q to quit.\n"

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
