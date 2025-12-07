package packselector

import (
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/view"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type PackSelectMsg struct {
	Pack balatro.Pack
}

type PackModel struct {
	Pack       balatro.Pack
	dimensions view.Dimensions
	selected   bool
}

func (c PackModel) Init() tea.Cmd {
	return nil
}

func (c *PackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c PackModel) View() string {
	emptyLine := strings.Repeat(" ", c.dimensions.Width)
	str := PackStr(c.Pack, c.dimensions)
	if !c.selected {
		str = emptyLine + "\n" + emptyLine + "\n" + str
	} else {
		str = str + "\n" + emptyLine + "\n" + emptyLine
	}
	return str
}

func (p PackModel) Select() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return PackSelectMsg{Pack: p.Pack}
	}
}

func (p PackModel) String(selected bool) string {
	return ""
}

func PackStr(pack balatro.Pack, dimensions view.Dimensions) string {
	if dimensions.Width < 16 || dimensions.Height < 8 {
		panic("Min card dimensions is 5 height x 7 Width")
	}

	lineTop := "┌" + strings.Repeat("─", dimensions.Width-2) + "┐\n"
	lineBottom := "└" + strings.Repeat("─", dimensions.Width-2) + "┘"

	packSizeName := "│"
	packSizeName += ui.CenterString(pack.Size().String(), dimensions.Width-2, " ")
	packSizeName += "│\n"

	packTypeName := "│"
	packTypeName += ui.CenterString(pack.String(), dimensions.Width-2, " ")
	packTypeName += "│\n"

	priceTag := "│" + strings.Repeat(" ", dimensions.Width-4) + strconv.Itoa(pack.Size().Price) + "c│\n"

	emptyLine := "│" + strings.Repeat(" ", dimensions.Width-2) + "│\n"

	result := lineTop
	for i := 0; i < (dimensions.Height-5)/2; i++ {
		result += emptyLine
	}
	result += packSizeName
	result += packTypeName
	for i := 0; i < dimensions.Height-5-(dimensions.Height-5)/2; i++ {
		result += emptyLine
	}
	result += priceTag
	result += lineBottom

	return result
}
