package menu

import (
	"golatro/internal/balatro/ui"
	"golatro/pkg/balatro"
	"golatro/pkg/teafx/control"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type LogsMenu struct {
	returnTo tea.Model

	viewport viewport.Model

	show   bool
	inputs control.Inputs
}

func (m LogsMenu) GetInputs() control.Inputs {
	return m.inputs
}

func (m LogsMenu) ToggleShow() control.Inputer {
	m.show = !m.show
	return m
}

func (m LogsMenu) Shows() bool {
	return m.show
}

func NewLogsMenu(model tea.Model, gameState balatro.GameState) LogsMenu {
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

	vp.SetContent(gameState.GameLogger.All())
	vp.GotoBottom()

	return LogsMenu{
		inputs: inputs,
		show:   false,

		viewport: vp,

		returnTo: model,
	}
}

func (m LogsMenu) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen)
}

func (m LogsMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m LogsMenu) View() string {
	s := ui.WhiteForeground.Render("Game logs") + "\n\n\n"

	s += m.viewport.View() + "\n"

	s += "\n\n"

	s += control.DisplayHelp(m)

	s += "\n\n"
	return s
}
