package teafx

import (
	"golatro/pkg/teafx/view"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestFitHeightTextShorter(t *testing.T) {
	height := 5
	text := `a
bc
def
ghij`

	textFitted := view.FitHeight(text, height)
	assert.Equal(t, height-1, strings.Count(textFitted, "\n"))
	assert.Equal(t, `a
bc
def
ghij
`, textFitted)
}

func TestFitHeightTextTaller(t *testing.T) {
	height := 3
	text := `a
bc
def
ghij`

	textFitted := view.FitHeight(text, height)
	assert.Equal(t, height-1, strings.Count(textFitted, "\n"))
	assert.Equal(t, `a
bc
def`, textFitted)
}

func TestFitHeightAlreadyFits(t *testing.T) {
	height := 3
	text := `a
bc
def`

	textFitted := view.FitHeight(text, height)
	assert.Equal(t, height-1, strings.Count(textFitted, "\n"))
	assert.Equal(t, `a
bc
def`, textFitted)
}

func TestFitWidthTextShorter(t *testing.T) {
	width := 5
	text := `abcdefg
hij
klmno
pq`

	textFitted := view.FitWidth(text, width)
	assert.Equal(t, `abcde
hij  
klmno
pq   `, textFitted)
}

func TestFitDimensionsTextShorter(t *testing.T) {
	dimensions := view.Dimensions{Height: 3, Width: 5}
	text := `abcdefg
hij
klmno
pq`

	textFitted := view.FitDimensions(text, dimensions)
	assert.Equal(t, `abcde
hij  
klmno`, textFitted)
}

func TestSingleCardVisual(t *testing.T) {
	dimensions := view.Dimensions{Height: 9, Width: 13}
	card := `┌───────────┐
│2          │
│           │
│     ♥     │
│           │
│          2│
└───────────┘`

	textFitted := view.FitDimensions(card, dimensions)
	assert.Equal(t, `┌───────────┐
│2          │
│           │
│     ♥     │
│           │
│          2│
└───────────┘
             
             `, textFitted)
}

func TestFitWidthColoredText(t *testing.T) {
	dimensions := view.Dimensions{Height: 1, Width: 3}
	txt := lipgloss.NewStyle().Foreground(lipgloss.Color("#CC0000")).Render("123")

	textFitted := view.FitDimensions(txt, dimensions)
	assert.Equal(t, txt, textFitted)
}
