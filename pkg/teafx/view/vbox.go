package view

import tea "github.com/charmbracelet/bubbletea"

type VBox struct {
	models     []tea.Model
	spacing    int // empty lines of space between each model in the vbox
	dimensions Dimensions
}

func NewSimpleVBox(dimensions Dimensions, models ...tea.Model) VBox {
	return VBox{
		models:     models,
		spacing:    1,
		dimensions: dimensions,
	}
}

func NewVBox(dimensions Dimensions, space int, models ...tea.Model) VBox {
	return VBox{
		models:     models,
		spacing:    space,
		dimensions: dimensions,
	}
}

func (v VBox) Init() tea.Cmd {
	return nil
}

func (v VBox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds, iCmd tea.Cmd
	for i := range len(v.models) {
		model := v.models[i]
		model, iCmd = model.Update(msg)
		v.models[i] = model
		cmds = tea.Batch(cmds, iCmd)
	}
	return v, cmds
}

func (v VBox) View() string {
	str := ""
	for _, model := range v.models {
		str += model.View() + "\n"
		for range v.spacing {
			str += "\n"
		}
	}
	return FitDimensions(str, v.dimensions)
}
