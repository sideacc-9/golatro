package menu

import (
	"fmt"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SelectUpgradeMenu struct {
	cards    []PackCardOption
	packSize balatro.PackSize
	pack     balatro.Pack

	selected []PackCardOption

	gameState *balatro.GameState
	shopState balatro.ShopState

	cursor int

	show   bool
	inputs control.Inputs

	errorMsg string
}

func (m SelectUpgradeMenu) GetSelected() control.Option {
	return m.cards[m.cursor]
}

func (m *SelectUpgradeMenu) Up() tea.Cmd {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.cards) - 1
	}
	return nil
}

func (m *SelectUpgradeMenu) Down() tea.Cmd {
	m.cursor++
	if m.cursor >= len(m.cards) {
		m.cursor = 0
	}
	return nil
}

func (m SelectUpgradeMenu) GetInputs() control.Inputs {
	return m.inputs
}

func (m SelectUpgradeMenu) ToggleShow() control.Inputer {
	m.show = !m.show
	return m
}

func (m SelectUpgradeMenu) Shows() bool {
	return m.show
}

type ConfirmSelectionMsg struct{}

var ConfirmSelection = func() tea.Msg {
	return ConfirmSelectionMsg{}
}

var CONFIRMSELECTION = control.Input{
	Keys:        []string{"enter"},
	Description: "Confirm selection",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, ConfirmSelection
	},
}

type SkipMsg struct{}

var SkipSelection = func() tea.Msg {
	return SkipMsg{}
}

var SKIPSELECTION = control.Input{
	Keys:        []string{"x"},
	Description: "Skip pack",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, SkipSelection
	},
}

func NewSelectUpgradeMenu(gameState balatro.GameState, shopState balatro.ShopState, pack balatro.Pack) SelectUpgradeMenu {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(control.LEFT)
	inputs.Add(control.RIGHT)
	inputs.Add(CONFIRMSELECTION)
	inputs.Add(SKIPSELECTION)
	inputs.Add(MANAGEJOKERS)
	inputs.Add(control.QUIT)
	inputs.Add(control.SELECT)
	inputs.Add(control.HELP)

	upgradesRaw := pack.Open()
	upgrades := make([]PackCardOption, 0, len(upgradesRaw))
	for i, u := range upgradesRaw {
		upgrades = append(upgrades, PackCardOption{Card: u, Id: i})
	}
	return SelectUpgradeMenu{
		cursor:   0,
		cards:    upgrades,
		packSize: pack.Size(),
		pack:     pack,

		inputs: inputs,
		show:   false,

		gameState: &gameState,
		shopState: shopState,
	}
}

func (m SelectUpgradeMenu) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen)
}

func (m SelectUpgradeMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	var model tea.Model = &m

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
	case UpgradeSelectMsg:
		isSelected := slices.ContainsFunc(m.selected, func(card PackCardOption) bool { return card.Id == msg.Id })
		if isSelected {
			m.selected = slices.DeleteFunc(m.selected, func(card PackCardOption) bool { return card.Id == msg.Id })
		} else if len(m.selected) < m.packSize.ChoosableCards {
			m.selected = append(m.selected, PackCardOption(msg))
		}
	case ConfirmSelectionMsg:
		if len(m.selected) == m.packSize.ChoosableCards {
			erred := false
			for _, u := range m.selected {
				err := u.Card.Apply(m.gameState)
				if err != nil {
					m.errorMsg = err.Error()
					cmd = tea.Batch(cmd, control.TimedCmd(2*time.Second, control.ClearErrorMsg{}))
					erred = true
					break
				} else {
					m.gameState.GameLogger.Add(balatro.UPGRADE_APPLIED, u.Card.String())
				}
			}
			if !erred {
				model = NewShopMenu(*m.gameState, &m.shopState)
			}
		}
	case SkipMsg:
		model = NewShopMenu(*m.gameState, &m.shopState)
	case ManageJokersMsg:
		model = NewJokerManager(m.gameState, m)
	case control.ClearErrorMsg:
		m.errorMsg = ""
	}

	return model, cmd
}

func (m SelectUpgradeMenu) View() string {
	s := fmt.Sprintf("%s %s - Choose %v", m.packSize, m.pack, m.packSize.ChoosableCards)

	s += "\n\n"

	s += DisplayPackCards(m.cards, m.selected, m.cursor)

	s += "\nDescription: " + m.cards[m.cursor].Card.Help()
	if m.errorMsg != "" {
		s += "\nError: " + m.errorMsg
	}

	s += "\n\n"

	s += control.DisplayHelp(m)

	s += "\n\n"
	return s
}
