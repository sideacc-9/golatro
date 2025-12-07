package cardselector

import (
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/view"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type CardSelectMsg struct {
	Card balatro.Card
}

type CardModel struct {
	card       balatro.Card
	dimensions view.Dimensions
	selected   bool
}

func (c CardModel) Init() tea.Cmd {
	return nil
}

func (c *CardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c CardModel) View() string {
	emptyLine := strings.Repeat(" ", c.dimensions.Width)
	str := CardStr(c.card, c.dimensions)
	if !c.selected {
		str = emptyLine + "\n" + emptyLine + "\n" + str
	} else {
		str = str + "\n" + emptyLine + "\n" + emptyLine
	}
	return str
}

func (c CardModel) Select() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return CardSelectMsg{Card: c.card}
	}
}

func (c CardModel) String(selected bool) string {
	return ""
}

func CardStr(card balatro.Card, dimensions view.Dimensions) string {
	if dimensions.Width < 7 || dimensions.Height < 5 {
		panic("Min card dimensions is 5 height x 7 Width")
	}

	lineTop := "┌" + strings.Repeat("─", dimensions.Width-2) + "┐\n"
	lineBottom := "└" + strings.Repeat("─", dimensions.Width-2) + "┘"

	rankEnhancementTop := "│"
	rankEnhancementTop += ui.PadRight(card.Rank.String(), 2, " ")
	rankEnhancementTop += strings.Repeat(" ", dimensions.Width-7) // total dimensions.Width -2 from borders, 2 from rank, 3 from enhancement abbreviation
	rankEnhancementTop += ui.PadLeft(card.Enhancement.Abbreviation(), 3, " ")
	rankEnhancementTop += "│\n"

	suit := card.Suit.Symbol
	if card.Suit == balatro.Hearts || card.Suit == balatro.Diamonds {
		suit = ui.Red.Render(suit)
	}
	suitLine := "│"
	suitLine += strings.Repeat(" ", (dimensions.Width-3)/2)
	suitLine += suit
	suitLine += strings.Repeat(" ", (dimensions.Width-3)/2)
	suitLine += "│\n"

	editionRankBottom := "│"
	editionRankBottom += ui.PadRight(card.Edition.Abbreviation(), 3, " ")
	editionRankBottom += strings.Repeat(" ", dimensions.Width-7)
	editionRankBottom += ui.PadLeft(card.Rank.String(), 2, " ")
	editionRankBottom += "│\n"

	emptyLine := "│" + strings.Repeat(" ", dimensions.Width-2) + "│\n"

	result := lineTop + rankEnhancementTop
	for i := 0; i < (dimensions.Height-5)/2; i++ {
		result += emptyLine
	}
	result += suitLine
	for i := 0; i < dimensions.Height-5-(dimensions.Height-5)/2; i++ {
		result += emptyLine
	}
	result += editionRankBottom + lineBottom

	return result
}
