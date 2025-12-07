package upgradeselector

import (
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"
	"golatro/pkg/teafx/view"

	tea "github.com/charmbracelet/bubbletea"
)

type UpgradeSelector struct {
	Upgrades []balatro.Consumable

	Cursor int
	inputs control.Inputs

	upgradeModels []*UpgradeModel
	hboxUpgrades  view.HBox

	priceFunc func(balatro.Consumable) int
}

var UPGRADE_DIMENSIONS = view.Dimensions{Height: 8, Width: 16}
var HBOX_DIMENSIONS = view.Dimensions{Height: 10, Width: UPGRADE_DIMENSIONS.Width*2 + 1}

var NONE_CONSUMABLE = UpgradeModel{Consumable: nil, dimensions: UPGRADE_DIMENSIONS, selected: false}

func (m UpgradeSelector) GetSelected() control.Option {
	if m.Cursor < 0 || m.Cursor >= len(m.Upgrades) {
		return NONE_CONSUMABLE
	}
	return UpgradeModel{Consumable: m.Upgrades[m.Cursor], dimensions: UPGRADE_DIMENSIONS, selected: true}
}

func (m *UpgradeSelector) Up() tea.Cmd {
	m.Cursor--
	if m.Cursor < 0 {
		m.Cursor = len(m.Upgrades) - 1
	}
	m.updateHbox()
	return nil
}

func (m *UpgradeSelector) Down() tea.Cmd {
	m.Cursor++
	if m.Cursor >= len(m.Upgrades) {
		m.Cursor = 0
	}
	m.updateHbox()
	return nil
}

func New(upgrades []balatro.Consumable) UpgradeSelector {
	return NewWithPriceFunc(upgrades, func(c balatro.Consumable) int { return c.Price() })
}

func NewWithPriceFunc(upgrades []balatro.Consumable, priceFunc func(balatro.Consumable) int) UpgradeSelector {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(control.LEFT)
	inputs.Add(control.RIGHT)
	inputs.Add(control.QUIT)
	inputs.Add(control.SELECT)

	cs := UpgradeSelector{
		Upgrades:  upgrades,
		inputs:    inputs,
		Cursor:    -1,
		priceFunc: priceFunc,
	}
	cs.updateHbox()
	return cs
}

func (m UpgradeSelector) Init() tea.Cmd {
	return nil
}

func (m UpgradeSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.hboxUpgrades.Update(msg)
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

func (m UpgradeSelector) View() string {
	str := m.hboxUpgrades.View()
	return str
}

func (m *UpgradeSelector) updateHbox() {
	m.upgradeModels = upgradesToModels(m.Upgrades, m.Cursor, m.priceFunc)

	m.hboxUpgrades = view.NewSimpleHBox(HBOX_DIMENSIONS)
	for _, cm := range m.upgradeModels {
		m.hboxUpgrades.Add(cm)
	}
}

func upgradesToModels(upgrades []balatro.Consumable, cursor int, priceFunc func(balatro.Consumable) int) []*UpgradeModel {
	opts := make([]*UpgradeModel, len(upgrades))
	for i, upgrade := range upgrades {
		isSelected := i == cursor
		opts[i] = &UpgradeModel{Consumable: upgrade, dimensions: UPGRADE_DIMENSIONS, selected: isSelected, priceFunc: priceFunc}
	}
	return opts
}

type ToggleOrderMsg struct{}

var ToggleOrder = func() tea.Msg {
	return ToggleOrderMsg{}
}
