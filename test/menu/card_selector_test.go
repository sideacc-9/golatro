package menu

import (
	"fmt"
	"golatro/internal/balatro/menu/components/cardselector"
	"golatro/pkg/balatro"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDisplaySomeSelected(t *testing.T) {
	cards := []balatro.Card{
		balatro.NewBasicCard(balatro.RandomRank(), balatro.RandomSuit()),
		balatro.NewBasicCard(balatro.RandomRank(), balatro.RandomSuit()),
		balatro.NewBasicCard(balatro.RandomRank(), balatro.RandomSuit()),
		balatro.NewBasicCard(balatro.RandomRank(), balatro.RandomSuit()),
	}
	var selector tea.Model = cardselector.New(cards, 5)
	selector, _ = selector.Update(cardselector.CardSelectMsg{Card: cards[0]})
	selector, _ = selector.Update(cardselector.CardSelectMsg{Card: cards[1]})
	fmt.Println(selector.View())
}
