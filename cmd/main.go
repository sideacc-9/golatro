package main

import (
	"fmt"
	"golatro/internal/balatro/menu"
	"golatro/pkg/balatro"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	gameState := balatro.NewBasicGameState()
	p := tea.NewProgram(menu.NewRoundMenu(gameState))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
