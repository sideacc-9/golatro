package menu

import (
	"fmt"
	"golatro/internal/balatro/menu/components/cardselector"
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type RoundMenu struct {
	gameState  balatro.GameState
	roundState balatro.RoundState

	cardSelector cardselector.CardSelector

	lastSum  int
	lastMult float64
	handType *balatro.HandType

	show   bool
	inputs control.Inputs
}

func (m RoundMenu) GetInputs() control.Inputs {
	return m.inputs
}

func (m RoundMenu) ToggleShow() control.Inputer {
	m.show = !m.show
	return m
}

func (m RoundMenu) Shows() bool {
	return m.show
}

var PLAYHAND = control.Input{
	Keys:        []string{"enter"},
	Description: "Play hand",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, PlayHand
	},
}

var DISCARD = control.Input{
	Keys:        []string{"x"},
	Description: "Discard hand",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, Discard
	},
}

var TOGGLESTATS = control.Input{
	Keys:        []string{"s"},
	Description: "Toggle show stats",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, ToggleStats
	},
}

var TOGGLELOGS = control.Input{
	Keys:        []string{"l"},
	Description: "Toggle show game logs",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, ToggleLogs
	},
}

func NewRoundMenu(gameState balatro.GameState) RoundMenu {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(control.QUIT)
	inputs.Add(PLAYHAND)
	inputs.Add(DISCARD)
	inputs.Add(TOGGLESTATS)
	inputs.Add(TOGGLELOGS)
	inputs.Add(control.HELP)

	roundState := gameState.NextRound()
	return RoundMenu{
		gameState:  gameState,
		roundState: roundState,

		cardSelector: cardselector.New(roundState.SelectableCards, gameState.MaxHandSize),

		inputs: inputs,
		show:   false,
	}
}

func (m RoundMenu) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen)
}

func (m RoundMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model = &m

	updatedCardSelector, cmd := m.cardSelector.Update(msg)
	m.cardSelector = updatedCardSelector.(cardselector.CardSelector)
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
	case PlayHandMsg:
		switch m.roundState.RoundStatus(m.gameState.MaxHands) {
		case balatro.WON:
			return NewShopMenu(m.gameState, nil), nil
		case balatro.LOST:
			return m, tea.Quit
		}
		if m.roundState.RoundStatus(m.gameState.MaxHands) == balatro.WON {
			return NewShopMenu(m.gameState, nil), nil
		}
		if len(m.cardSelector.SelectedCards) == 0 {
			break
		}
		cursor := m.cardSelector.Cursor
		m.handType, m.lastSum, m.lastMult = m.roundState.ScoreHand(&m.gameState, m.cardSelector.SelectedCards)
		m.cardSelector.Reset(m.roundState.SelectableCards, m.gameState.MaxHandSize, cursor)
	case DiscardMsg:
		if m.roundState.Discards >= m.gameState.MaxDiscards || len(m.cardSelector.SelectedCards) == 0 ||
			m.roundState.RoundStatus(m.gameState.MaxHands) != balatro.INPROGRESS {
			break
		}
		cursor := m.cardSelector.Cursor
		m.roundState.Discard(m.cardSelector.SelectedCards, m.gameState.MaxSelectableCards)
		m.cardSelector.Reset(m.roundState.SelectableCards, m.gameState.MaxHandSize, cursor)
	case ToggleStatsMsg:
		model = NewStatusMenu(m, m.gameState, &m.roundState)
	case ToggleLogsMsg:
		model = NewLogsMenu(m, m.gameState)
	}

	return model, cmd
}

func (m RoundMenu) View() string {
	s := ""

	s += "Round: " + strconv.Itoa(m.gameState.Round) + " / 3 - Ante: " + strconv.Itoa(m.gameState.Ante) + " / 8\n"
	s += "─────────────────────────────\n"
	s += "Hands: " + strconv.Itoa(m.gameState.MaxHands-m.roundState.Hand+1) + "\n\n"
	s += "Discards: " + strconv.Itoa(m.gameState.MaxDiscards-m.roundState.Discards) + "\n\n"
	s += fmt.Sprintf("Points: %v / %v", m.roundState.Points, m.roundState.Target)
	if m.handType != nil {
		s += fmt.Sprintf(" (%v x %v) %v", m.lastSum, m.lastMult, m.handType)
	}
	s += "\n\n"

	s += m.cardSelector.View() + "\n\n"

	var cardDesc = "-"
	if m.cardSelector.Cursor >= 0 && m.cardSelector.Cursor < len(m.roundState.SelectableCards) {
		cardDesc = m.roundState.SelectableCards[m.cardSelector.Cursor].String()
	}
	s += "Card details: " + cardDesc

	s += "\n\n"

	switch m.roundState.RoundStatus(m.gameState.MaxHands) {
	case balatro.WON:
		s += ui.WhiteForeground.Render("> Next") + "\n\n"
	case balatro.LOST:
		s += "YOU LOST!\n\n"
		s += "Check your stats with '" + TOGGLESTATS.Keys[0] + "' or the game logs with '" + TOGGLELOGS.Keys[0] + "'"
	}

	s += control.DisplayHelp(m)

	s += "\n\n"
	return s
}

type GameResultMsg struct {
	won bool
}

var GameResult = func(won bool) tea.Cmd {
	return func() tea.Msg {
		return GameResultMsg{won: won}
	}
}

type ToggleStatsMsg struct{}

var ToggleStats = func() tea.Msg {
	return ToggleStatsMsg{}
}

type ToggleLogsMsg struct{}

var ToggleLogs = func() tea.Msg {
	return ToggleLogsMsg{}
}

type PlayHandMsg struct{}

var PlayHand = func() tea.Msg {
	return PlayHandMsg{}
}

type DiscardMsg struct{}

var Discard = func() tea.Msg {
	return DiscardMsg{}
}
