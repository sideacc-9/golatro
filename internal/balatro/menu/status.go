package menu

import (
	"fmt"
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type StatusMenu struct {
	returnTo tea.Model

	viewport viewport.Model

	show   bool
	inputs control.Inputs
}

func (m StatusMenu) GetInputs() control.Inputs {
	return m.inputs
}

func (m StatusMenu) ToggleShow() control.Inputer {
	m.show = !m.show
	return m
}

func (m StatusMenu) Shows() bool {
	return m.show
}

type ReturnMsg struct{}

var ReturnPrev = func() tea.Msg {
	return ReturnMsg{}
}

var RETURNPREV = control.Input{
	Keys:        []string{tea.KeyBackspace.String()},
	Description: "Return",
	Action: func(m tea.Model) (tea.Model, tea.Cmd) {
		return m, ReturnPrev
	},
}

func NewStatusMenu(model tea.Model, gameState balatro.GameState, roundState *balatro.RoundState) StatusMenu {
	inputs := control.Inputs{
		Inputs: make(map[string]control.Input),
		Order:  make([]string, 0),
	}

	inputs.Add(control.QUIT)
	inputs.Add(RETURNPREV)
	inputs.Add(control.HELP)

	vp := viewport.New(120, 25)
	vp.KeyMap.Down.SetKeys("down")
	vp.KeyMap.HalfPageDown.SetKeys("ctrl+down")
	vp.KeyMap.PageDown.SetKeys("pgdown")
	vp.KeyMap.Up.SetKeys("up")
	vp.KeyMap.HalfPageUp.SetKeys("ctrl+up")
	vp.KeyMap.PageUp.SetKeys("pgup")

	content := ""
	if roundState != nil {
		content = RoundStateString(*roundState) + "\n\n"
	}
	content += GameStateString(gameState)

	vp.SetContent(content)

	return StatusMenu{
		inputs: inputs,
		show:   false,

		viewport: vp,

		returnTo: model,
	}
}

func (m StatusMenu) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen)
}

func (m StatusMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model

	m.viewport, cmd = m.viewport.Update(msg)

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
	case ReturnMsg:
		model = m.returnTo
	}

	if model == nil {
		model = m
	}

	return model, cmd
}

func (m StatusMenu) View() string {
	s := ui.WhiteForeground.Render("Game/round stats") + "\n\n\n"

	s += m.viewport.View() + "\n"

	s += "\n\n"

	s += control.DisplayHelp(m)

	s += "\n\n"
	return s
}

func GameStateString(g balatro.GameState) string {
	str := ui.WhiteForeground.Render("> Game stats") + "\n\n\n"

	str += ui.WhiteForeground.Render("Jokers") + "\n"

	for _, j := range g.Jokers {
		str += fmt.Sprintf("%v: %v\n", j.Type.String(), j.Type.Help())
	}

	str += "\n\n" + ui.WhiteForeground.Render("Deck") + "\n"
	balatro.Hand(g.Deck).SortFunc(balatro.CompareSuitRank)
	var prevSuit *balatro.Suit = nil
	for _, c := range g.Deck {
		if prevSuit == nil || c.Suit != *prevSuit {
			prevSuit = &c.Suit
			str += "\n" + c.Suit.Symbol + ": "
		}
		str += c.Rank.String()
		edStr := c.Edition.String()
		ehStr := c.Enhancement.String()
		if edStr != "" || ehStr != "" {
			str += " ("
			if edStr != "" {
				str += edStr
			}
			if edStr != "" && ehStr != "" {
				str += " "
			}
			if ehStr != "" {
				str += ehStr
			}
			str += ")"
		}
		str += ": "
	}

	str += "\n\n" + ui.WhiteForeground.Render("Hand levels") + "\n"
	for typ, lvl := range g.HandLevels {
		sum, mult := lvl.GetSumMult()
		str += fmt.Sprintf("%v (%v): %vx%v\n", typ.String(), lvl.Level, sum, mult)
	}

	return str
}

func RoundStateString(r balatro.RoundState) string {
	str := ui.WhiteForeground.Render("> Round stats") + "\n\n"

	str += ui.WhiteForeground.Render("Available cards") + "\n\n"
	balatro.Hand(r.AvailableCards).SortFunc(balatro.CompareSuitRank)
	var prevSuit *balatro.Suit = nil
	for _, c := range r.AvailableCards {
		if prevSuit == nil || c.Suit != *prevSuit {
			prevSuit = &c.Suit
			str += "\n" + c.Suit.Symbol + ": "
		}
		str += c.Rank.String()
		edStr := c.Edition.String()
		ehStr := c.Enhancement.String()
		if edStr != "" || ehStr != "" {
			str += " ("
			if edStr != "" {
				str += edStr
			}
			if edStr != "" && ehStr != "" {
				str += " "
			}
			if ehStr != "" {
				str += ehStr
			}
			str += ")"
		}
		str += ": "
	}

	return str
}
