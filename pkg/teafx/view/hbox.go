package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type HBox struct {
	Models     []tea.Model
	spacing    int // empty lines of space between each model in the hbox
	dimensions Dimensions
}

func NewSimpleHBox(dimensions Dimensions, models ...tea.Model) HBox {
	return HBox{
		Models:     models,
		spacing:    1,
		dimensions: dimensions,
	}
}

func NewHBox(dimensions Dimensions, space int, models ...tea.Model) HBox {
	return HBox{
		Models:     models,
		spacing:    space,
		dimensions: dimensions,
	}
}

func (h *HBox) Add(model tea.Model) {
	h.Models = append(h.Models, model)
}

func (h HBox) Init() tea.Cmd {
	return nil
}

func (h HBox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds, iCmd tea.Cmd
	for i := range len(h.Models) {
		model := h.Models[i]
		model, iCmd = model.Update(msg)
		h.Models[i] = model
		cmds = tea.Batch(cmds, iCmd)
	}
	return h, cmds
}

func (h HBox) View() string {
	modelStrings := make([][]string, 0, len(h.Models))
	for _, model := range h.Models {
		modelStr := model.View()
		modelStrings = append(modelStrings, strings.Split(modelStr, "\n"))
	}
	result := ""
	for i := range h.dimensions.Height {
		for _, modelLines := range modelStrings {
			line := ""
			if i >= len(modelLines) {
				line = strings.Repeat(" ", len(modelLines[len(modelLines)-1]))
			} else {
				line = modelLines[i]
			}
			result += line + strings.Repeat(" ", h.spacing)
		}
		if i != h.dimensions.Height-1 {
			result += "\n"
		}
	}
	return result
}
