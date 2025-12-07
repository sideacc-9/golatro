package balatro

import (
	"golatro/pkg/balatro"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultFullHouseUpgrade(t *testing.T) {
	game := balatro.NewBasicGameState()
	fhLevel := game.HandLevels[balatro.FullHouse]
	fhLevel.Level = 2
	game.HandLevels[balatro.FullHouse] = fhLevel // 65 * 6 initial score

	round := game.NextRound()
	round.SelectableCards = []balatro.Card{
		balatro.NewCard(balatro.R3, balatro.Clubs, balatro.Normal{}, balatro.NewMult()),
		balatro.NewBasicCard(balatro.R3, balatro.Hearts),
		balatro.NewCard(balatro.R5, balatro.Clubs, balatro.Normal{}, balatro.NewMult()),
		balatro.NewBasicCard(balatro.R5, balatro.Diamonds),
		balatro.NewBasicCard(balatro.R5, balatro.Spades),
		balatro.NewCard(balatro.K, balatro.Clubs, balatro.Normal{}, balatro.NewSteel()),
		balatro.NewBasicCard(balatro.Q, balatro.Diamonds),
		balatro.NewBasicCard(balatro.J, balatro.Spades),
	}

	typeH, s, m := round.ScoreHand(&game, balatro.Hand(round.SelectableCards[:5]))

	assert.Equal(t, balatro.FullHouse, *typeH)
	assert.Equal(t, (6.0+4+4)*1.5, m)
	assert.Equal(t, 86, s)
	assert.Equal(t, 86*((6.0+4+4)*1.5), round.Points)
}
