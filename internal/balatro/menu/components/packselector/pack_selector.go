package packselector

import (
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"
	"golatro/pkg/teafx/view"

	tea "github.com/charmbracelet/bubbletea"
)

type PackSelector struct {
	Packs []balatro.Pack

	Cursor int
	inputs control.Inputs

	packModels []*PackModel
	hboxPacks  view.HBox
}

var PACK_DIMENSIONS = view.Dimensions{Height: 8, Width: 16}
var HBOX_DIMENSIONS = view.Dimensions{Height: 10, Width: PACK_DIMENSIONS.Width*2 + 1} // width: 2 packs + 1 spacing, height: pack height + 2 for selection space

var NONE_PACK = PackModel{Pack: nil, dimensions: PACK_DIMENSIONS, selected: false}

func (m PackSelector) GetSelected() control.Option {
	if m.Cursor < 0 || m.Cursor >= len(m.Packs) {
		return NONE_PACK
	}
	return PackModel{Pack: m.Packs[m.Cursor], dimensions: PACK_DIMENSIONS, selected: true}
}

func (m *PackSelector) Up() tea.Cmd {
	m.Cursor--
	if m.Cursor < 0 {
		m.Cursor = len(m.Packs) - 1
	}
	m.updateHbox()
	return nil
}

func (m *PackSelector) Down() tea.Cmd {
	m.Cursor++
	if m.Cursor >= len(m.Packs) {
		m.Cursor = 0
	}
	m.updateHbox()
	return nil
}

func New(packs []balatro.Pack) PackSelector {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(control.LEFT)
	inputs.Add(control.RIGHT)
	inputs.Add(control.QUIT)
	inputs.Add(control.SELECT)

	cs := PackSelector{
		Packs:  packs,
		inputs: inputs,
		Cursor: -1,
	}
	cs.updateHbox()
	return cs
}

func (m PackSelector) Init() tea.Cmd {
	return nil
}

func (m PackSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.hboxPacks.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		input, ok := m.inputs.Inputs[msg.String()]
		if ok {
			_, cmdTemp := input.Action(&m)

			if cmdTemp != nil {
				cmd = tea.Batch(cmd, cmdTemp)
			}
		}
	}
	return &m, cmd
}

func (m PackSelector) View() string {
	str := m.hboxPacks.View()
	return str
}

func (m *PackSelector) updateHbox() {
	m.packModels = packsToModels(m.Packs, m.Cursor)

	m.hboxPacks = view.NewSimpleHBox(HBOX_DIMENSIONS)
	for _, cm := range m.packModels {
		m.hboxPacks.Add(cm)
	}
}

func packsToModels(cards []balatro.Pack, cursor int) []*PackModel {
	opts := make([]*PackModel, len(cards))
	for i, pack := range cards {
		isSelected := i == cursor
		opts[i] = &PackModel{Pack: pack, dimensions: PACK_DIMENSIONS, selected: isSelected}
	}
	return opts
}

type ToggleOrderMsg struct{}

var ToggleOrder = func() tea.Msg {
	return ToggleOrderMsg{}
}
