package menu

import (
	"fmt"
	"golatro/internal/balatro/menu/components/packselector"
	"golatro/internal/balatro/menu/components/upgradeselector"
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ShopMenu struct {
	gameState *balatro.GameState
	shopState balatro.ShopState

	subMenus []control.Optioner
	cursor   int

	show   bool
	inputs control.Inputs

	errorMsg string
}

func (m *ShopMenu) Up() tea.Cmd {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.subMenus) - 1
	}
	return nil
}

func (m *ShopMenu) Down() tea.Cmd {
	m.cursor++
	if m.cursor >= len(m.subMenus) {
		m.cursor = 0
	}
	return nil
}

func (m ShopMenu) GetInputs() control.Inputs {
	return m.inputs
}

func (m ShopMenu) ToggleShow() control.Inputer {
	m.show = !m.show
	return m
}

func (m ShopMenu) Shows() bool {
	return m.show
}

type ExitShopMsg struct{}

var ExitShop = func() tea.Msg {
	return ExitShopMsg{}
}

var EXITSHOP = control.Input{
	Keys:        []string{"x"},
	Description: "Exit shop",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, ExitShop
	},
}

type ManageJokersMsg struct{}

var ManageJokers = func() tea.Msg {
	return ManageJokersMsg{}
}

var MANAGEJOKERS = control.Input{
	Keys:        []string{"j"},
	Description: "Exit shop",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, ManageJokers
	},
}

var UPSHOP = control.Input{
	Keys:        []string{"up"},
	Description: "Move up",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		shop, ok := m.(*ShopMenu)
		if !ok {
			return m, nil
		}
		shop.Up()
		return shop, nil
	},
}

var DOWNSHOP = control.Input{
	Keys:        []string{"down"},
	Description: "Move down",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		shop, ok := m.(*ShopMenu)
		if !ok {
			return m, nil
		}
		shop.Down()
		return shop, nil
	},
}

func NewShopMenu(gameState balatro.GameState, shopState *balatro.ShopState) ShopMenu {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(control.LEFT)
	inputs.Add(control.RIGHT)
	inputs.Add(control.QUIT)
	inputs.Add(UPSHOP)
	inputs.Add(DOWNSHOP)
	inputs.Add(TOGGLESTATS)
	inputs.Add(TOGGLELOGS)
	inputs.Add(MANAGEJOKERS)
	inputs.Add(EXITSHOP)
	inputs.Add(control.HELP)

	if shopState == nil {
		temp := balatro.NewShopState(2, 2)
		shopState = &temp
		for range 2 {
			shopState.Packs = append(shopState.Packs, balatro.GetRandomPack())
			shopState.Upgrades = append(shopState.Upgrades, balatro.GetRandomConsumable())
		}
	}

	packSelector := packselector.New(shopState.Packs)
	upgradesSelector := upgradeselector.New(shopState.Upgrades)

	return ShopMenu{
		cursor:   0,
		subMenus: []control.Optioner{&packSelector, &upgradesSelector},

		inputs: inputs,
		show:   false,

		gameState: &gameState,
		shopState: *shopState,
	}
}

func (m ShopMenu) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen)
}

func (m ShopMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model = &m

	submenu, cmd := m.subMenus[m.cursor].Update(msg)
	m.subMenus[m.cursor] = submenu.(control.Optioner)

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
	case upgradeselector.UpgradeSelectMsg:
		if msg.Upgrade != nil && m.gameState.Money-msg.Upgrade.Price() >= m.gameState.MinMoney {
			err := msg.Upgrade.Apply(m.gameState)
			if err != nil {
				m.errorMsg = err.Error()
				cmd = tea.Batch(cmd, control.TimedCmd(2*time.Second, control.ClearErrorMsg{}))
			} else {
				m.gameState.Money -= msg.Upgrade.Price()
				m.gameState.GameLogger.Add(balatro.UPGRADE_APPLIED, msg.Upgrade.String())
				m.gameState.GameLogger.Add(balatro.MONEY_SPENT, fmt.Sprintf("Bought %v for %vc", msg.Upgrade, msg.Upgrade.Price()))
				m.shopState.RemoveUpgrade(msg.Upgrade)
				UpdateShopMenus(&m)
			}
		}
	case packselector.PackSelectMsg:
		if msg.Pack != nil && m.gameState.Money-msg.Pack.Size().Price >= m.gameState.MinMoney {
			m.gameState.Money -= msg.Pack.Size().Price
			m.gameState.GameLogger.Add(balatro.MONEY_SPENT, fmt.Sprintf("Bought %v for %vc", msg.Pack, msg.Pack.Size().Price))
			m.shopState.RemovePack(msg.Pack)
			model = NewSelectUpgradeMenu(*m.gameState, m.shopState, msg.Pack)
		}
	case ExitShopMsg:
		model = NewRoundMenu(*m.gameState)
	case ToggleStatsMsg:
		model = NewStatusMenu(m, *m.gameState, nil)
	case ToggleLogsMsg:
		model = NewLogsMenu(m, *m.gameState)
	case ManageJokersMsg:
		model = NewJokerManager(m.gameState, m)
	case control.ClearErrorMsg:
		m.errorMsg = ""
	}

	return model, cmd
}

func UpdateShopMenus(m *ShopMenu) {
	packSelector := packselector.New(m.shopState.Packs)
	upgradesSelector := upgradeselector.New(m.shopState.Upgrades)
	m.subMenus = []control.Optioner{&packSelector, &upgradesSelector}
}

func (m ShopMenu) View() string {
	s := ui.WhiteForeground.Render("The store")

	s += "\n"
	s += fmt.Sprintf("Money: %vc\n\n", m.gameState.Money)

	var option balatro.Helper = nil
	for i, menu := range m.subMenus {
		if i == m.cursor {
			s += menu.View()
			menuOpt := menu.GetSelected()
			switch purchase := menuOpt.(type) {
			case packselector.PackModel:
				option = purchase.Pack
			case upgradeselector.UpgradeModel:
				option = purchase.Consumable
			}
		} else {
			s += ui.GreyedOut.Render(menu.View())
		}
		s += "\n"
	}

	help := ""
	if option != nil {
		help = option.Help()
	}
	s += "\nDescription: " + help
	if m.errorMsg != "" {
		s += "\nError: " + m.errorMsg
	}

	s += "\n\n"

	s += control.DisplayHelp(m)

	s += "\n\n"
	return s
}
