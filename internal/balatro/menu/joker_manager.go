package menu

import (
	"golatro/internal/balatro/menu/components/upgradeselector"
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"

	tea "github.com/charmbracelet/bubbletea"
)

type JokerManager struct {
	gameState *balatro.GameState

	jokerSelector *upgradeselector.UpgradeSelector

	returnTo tea.Model

	show   bool
	inputs control.Inputs
}

type SellMsg struct {
	idx int
}

var Sell = func(idx int) tea.Cmd {
	return func() tea.Msg {
		return SellMsg{idx: idx}
	}
}

var SELL = control.Input{
	Keys:        []string{"x"},
	Description: "Sell item",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		jManager, ok := m.(*JokerManager)
		if !ok {
			return m, nil
		}
		return m, Sell(jManager.jokerSelector.Cursor)
	},
}

func halfPrice(c balatro.Consumable) int {
	return c.Price() / 2
}

func NewJokerManager(gameState *balatro.GameState, returnTo tea.Model) JokerManager {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(SELL)
	inputs.Add(control.LEFT)
	inputs.Add(control.RIGHT)
	inputs.Add(RETURNPREV)
	inputs.Add(control.QUIT)
	inputs.Add(control.HELP)

	jokers := make([]balatro.Consumable, len(gameState.Jokers))
	for i, j := range gameState.Jokers {
		jokers[i] = balatro.JokerCard(j)
	}
	jokerSelector := upgradeselector.NewWithPriceFunc(jokers, halfPrice)
	cs := JokerManager{
		gameState:     gameState,
		jokerSelector: &jokerSelector,
		returnTo:      returnTo,
		inputs:        inputs,
	}
	return cs
}

func (m JokerManager) GetInputs() control.Inputs {
	return m.inputs
}

func (m JokerManager) ToggleShow() control.Inputer {
	m.show = !m.show
	return m
}

func (m JokerManager) Shows() bool {
	return m.show
}

func (m JokerManager) Init() tea.Cmd {
	return nil
}

func (m JokerManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model = &m
	var cmd tea.Cmd
	jSel, cmd := m.jokerSelector.Update(msg)
	m.jokerSelector = jSel.(*upgradeselector.UpgradeSelector)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		input, ok := m.inputs.Inputs[msg.String()]
		if ok {
			modelTemp, cmdTemp := input.Action(&m)
			if modelTemp != nil {
				model = modelTemp
			}
			if cmdTemp != nil {
				cmd = tea.Batch(cmd, cmdTemp)
			}
		}
	case SellMsg:
		if m.gameState.SellJoker(msg.idx, halfPrice) {
			jokers := make([]balatro.Consumable, len(m.gameState.Jokers))
			for i, j := range m.gameState.Jokers {
				jokers[i] = balatro.JokerCard(j)
			}
			jokerSelector := upgradeselector.NewWithPriceFunc(jokers, halfPrice)
			m.jokerSelector = &jokerSelector
		}
	case ReturnMsg:
		model = m.returnTo
	}
	return model, cmd
}

func (m JokerManager) View() string {
	s := ui.WhiteForeground.Render("Joker management")
	s += "\n\n"
	if len(m.gameState.Jokers) == 0 {
		s += "No jokers available"
	} else {
		s += m.jokerSelector.View()
	}

	help := ""
	jOption := m.jokerSelector.GetSelected()
	j := jOption.(upgradeselector.UpgradeModel)
	if j.Consumable != nil {
		help = j.Consumable.Help()
	}
	s += "\n\nDescription: " + help

	s += "\n\n"

	s += control.DisplayHelp(m)

	s += "\n\n"
	return s
}
