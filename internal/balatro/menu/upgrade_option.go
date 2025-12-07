package menu

import (
	"fmt"
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type UpgradeSelectMsg PackCardOption

type PackCardOption struct {
	Id   int
	Card balatro.Consumable
}

func (c PackCardOption) String(selected bool) string {
	str := ""
	if !selected {
		str += "                \n                \n"
	}
	emptySpaceLength := 14
	cardStrings := strings.Split(c.Card.String(), "\n")
	paddedName := ui.CenterString(cardStrings[0], emptySpaceLength, " ")
	paddedSubtitle := "              "
	if len(cardStrings) > 1 {
		paddedSubtitle = ui.CenterString(cardStrings[1], emptySpaceLength, " ")
	}
	str += fmt.Sprintf(`┌──────────────┐
│              │
│              │
│%v│
│%v│
│              │
└──────────────┘`, paddedName, paddedSubtitle)
	if selected {
		str += "\n                \n                "
	}
	return str
}
func (c PackCardOption) Select() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return UpgradeSelectMsg(c)
	}
}

func DisplayPackCards(cards, selected []PackCardOption, cursor int) string {
	packStrings := make([][]string, 0, len(cards))
	for i, c := range cards {
		isSelected := slices.ContainsFunc(selected, func(card PackCardOption) bool { return card.Id == c.Id })
		upgradeString := c.String(isSelected)
		if i == cursor {
			upgradeString += "\n        ▲       "
		} else {
			upgradeString += "\n                "
		}
		packStrings = append(packStrings, strings.Split(upgradeString, "\n"))
	}
	result := ""
	for i := range 10 {
		for _, cardLines := range packStrings {
			if len(cardLines) <= i {
				continue
			}
			result += cardLines[i] + " "
		}
		result += "\n"
	}
	return result
}
