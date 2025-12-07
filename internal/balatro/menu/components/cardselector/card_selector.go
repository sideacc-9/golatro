package cardselector

import (
	"cmp"
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"
	"golatro/pkg/teafx/view"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type CardSelector struct {
	Cards         []balatro.Card
	SelectedCards []balatro.Card
	sortingType   balatro.SortType
	sortingFunc   func(c1 balatro.Card, c2 balatro.Card) int
	maxSelectable int

	Cursor int
	inputs control.Inputs

	cardModels []*CardModel
	hboxCards  view.HBox
}

var CARD_DIMENSIONS = view.Dimensions{Height: 7, Width: 11}
var HBOX_DIMENSIONS = view.Dimensions{Height: 9, Width: CARD_DIMENSIONS.Width*9 - 1}

func (m CardSelector) GetSelected() control.Option {
	return CardModel{card: m.Cards[m.Cursor], dimensions: CARD_DIMENSIONS, selected: true}
}

func (m *CardSelector) Up() tea.Cmd {
	m.Cursor--
	if m.Cursor < 0 {
		m.Cursor = len(m.Cards) - 1
	}
	return nil
}

func (m *CardSelector) Down() tea.Cmd {
	m.Cursor++
	if m.Cursor >= len(m.Cards) {
		m.Cursor = 0
	}
	return nil
}

var TOGGLEORDER = control.Input{
	Keys:        []string{"b"},
	Description: "Toggle order (rank/suit)",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, ToggleOrder
	},
}

func New(selectableCards []balatro.Card, maxSelectable int) CardSelector {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(control.LEFT)
	inputs.Add(control.RIGHT)
	inputs.Add(control.QUIT)
	inputs.Add(TOGGLEORDER)
	inputs.Add(control.SELECT)

	cs := CardSelector{
		Cards:         selectableCards,
		SelectedCards: make([]balatro.Card, 0),
		sortingType:   balatro.RANK,
		sortingFunc:   balatro.CompareRankSuit,
		maxSelectable: maxSelectable,
		inputs:        inputs,
		Cursor:        0,
	}
	cs.orderAndUpdateHbox()
	return cs
}

func (m *CardSelector) Reset(selectableCards []balatro.Card, maxSelectable, initCursor int) {
	m.Cards = selectableCards
	m.maxSelectable = maxSelectable
	m.Cursor = initCursor
	m.SelectedCards = make([]balatro.Card, 0)
	m.orderAndUpdateHbox()
}

func (m CardSelector) Init() tea.Cmd {
	return nil
}

func (m CardSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.hboxCards.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		input, ok := m.inputs.Inputs[msg.String()]
		if ok {
			_, cmdTemp := input.Action(&m)

			if cmdTemp != nil {
				cmd = tea.Batch(cmd, cmdTemp)
			}
		}
	case ToggleOrderMsg:
		if m.sortingType == balatro.RANK {
			m.sortingFunc = balatro.CompareSuitRank
			m.sortingType = balatro.SUIT
		} else {
			m.sortingFunc = balatro.CompareRankSuit
			m.sortingType = balatro.RANK
		}
		m.orderAndUpdateHbox()
	case CardSelectMsg:
		wasSelected := slices.ContainsFunc(m.SelectedCards, func(c balatro.Card) bool { return cmp.Compare(c.Uuid.ID(), msg.Card.Uuid.ID()) == 0 })
		isSelected := wasSelected
		if wasSelected {
			m.SelectedCards = slices.DeleteFunc(m.SelectedCards, func(c balatro.Card) bool { return cmp.Compare(c.Uuid.ID(), msg.Card.Uuid.ID()) == 0 })
			isSelected = false
		} else if len(m.SelectedCards) < m.maxSelectable {
			m.SelectedCards = append(m.SelectedCards, msg.Card)
			isSelected = true
		}
		for i := range len(m.cardModels) {
			cModel := m.cardModels[i]
			if cmp.Compare(cModel.card.Uuid.ID(), msg.Card.Uuid.ID()) == 0 {
				m.cardModels[i].selected = isSelected
				break
			}
		}
	}
	return m, cmd
}

func (m CardSelector) View() string {
	str := m.hboxCards.View()
	str += "\n"
	widthCursorCell := CARD_DIMENSIONS.Width + 1
	for i := range len(m.Cards) {
		if m.Cursor == i {
			str += ui.CenterString("▲", widthCursorCell, " ")
		} else {
			str += strings.Repeat(" ", widthCursorCell)
		}
	}
	return str
}

func (m *CardSelector) orderAndUpdateHbox() {
	balatro.Hand(m.Cards).SortFunc(m.sortingFunc)
	m.cardModels = cardsToModels(m.Cards, m.SelectedCards)

	m.hboxCards = view.NewSimpleHBox(HBOX_DIMENSIONS)
	for _, cm := range m.cardModels {
		m.hboxCards.Add(cm)
	}
}

func cardsToModels(cards []balatro.Card, selected []balatro.Card) []*CardModel {
	opts := make([]*CardModel, len(cards))
	for i, card := range cards {
		isSelected := slices.ContainsFunc(selected, func(c balatro.Card) bool { return cmp.Compare(c.Uuid.ID(), card.Uuid.ID()) == 0 })
		opts[i] = &CardModel{card: card, dimensions: CARD_DIMENSIONS, selected: isSelected}
	}
	return opts
}

type ToggleOrderMsg struct{}

var ToggleOrder = func() tea.Msg {
	return ToggleOrderMsg{}
}
