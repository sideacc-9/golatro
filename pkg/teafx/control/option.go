package control

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Option interface {
	String(selected bool) string
	Select() (tea.Model, tea.Cmd)
}

type Exit struct{}

func (c Exit) Select() (tea.Model, tea.Cmd) {
	return nil, tea.Quit
}

var Bold = lipgloss.NewStyle().Bold(true)

func (c Exit) String(selected bool) string {
	if !selected {
		return "Exit"
	}
	return Bold.Render("Exit")
}

type ModelOption interface {
	tea.Model
	Option
}
