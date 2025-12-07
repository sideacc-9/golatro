package upgradeselector

import (
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/view"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type UpgradeSelectMsg struct {
	Upgrade balatro.Consumable
}

type UpgradeModel struct {
	Consumable balatro.Consumable
	dimensions view.Dimensions
	selected   bool
	priceFunc  func(balatro.Consumable) int
}

func (c UpgradeModel) Init() tea.Cmd {
	return nil
}

func (c *UpgradeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c UpgradeModel) View() string {
	emptyLine := strings.Repeat(" ", c.dimensions.Width)
	str := upgradeStr(c.Consumable, c.dimensions, c.priceFunc)
	if !c.selected {
		str = emptyLine + "\n" + emptyLine + "\n" + str
	} else {
		str = str + "\n" + emptyLine + "\n" + emptyLine
	}
	return str
}

func (p UpgradeModel) Select() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return UpgradeSelectMsg{Upgrade: p.Consumable}
	}
}

func (p UpgradeModel) String(selected bool) string {
	return ""
}

func upgradeStr(consumable balatro.Consumable, dimensions view.Dimensions, priceFunc func(balatro.Consumable) int) string {
	if dimensions.Width < 16 || dimensions.Height < 8 {
		panic("Min card dimensions is 8 height x 16 Width")
	}

	lineTop := "┌" + strings.Repeat("─", dimensions.Width-2) + "┐\n"
	lineBottom := "└" + strings.Repeat("─", dimensions.Width-2) + "┘"

	nameLines := strings.Split(consumable.String(), "\n")

	consumableNameLines := make([]string, len(nameLines))
	for i, nameLine := range nameLines {
		consumableNameLines[i] = "│"
		consumableNameLines[i] += ui.CenterString(nameLine, dimensions.Width-2, " ")
		consumableNameLines[i] += "│\n"
	}

	emptyLine := "│" + strings.Repeat(" ", dimensions.Width-2) + "│\n"
	priceTag := emptyLine
	if priceFunc != nil {
		price := priceFunc(consumable)
		priceTag = "│" + ui.PadLeft(strconv.Itoa(price), dimensions.Width-3, " ") + "c│\n"
	}

	nEmptyLines := dimensions.Height - 3 - len(consumableNameLines)
	result := lineTop
	for i := 0; i < nEmptyLines/2; i++ {
		result += emptyLine
	}
	for _, line := range consumableNameLines {
		result += line
	}
	for i := 0; i < nEmptyLines-nEmptyLines/2; i++ {
		result += emptyLine
	}
	result += priceTag
	result += lineBottom

	return result
}
