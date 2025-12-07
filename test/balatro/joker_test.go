package balatro

import (
	"golatro/pkg/balatro"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNA(t *testing.T) {
	game := balatro.NewBasicGameState()
	game.Jokers = append(game.Jokers, balatro.NewJoker(balatro.DNA))

	round := game.NextRound()

	assert.Equal(t, 52, len(game.Deck))
	c := round.SelectableCards[0]

	round.ScoreHand(&game, balatro.Hand([]balatro.Card{c}))

	assert.Equal(t, 53, len(game.Deck))

	rankMatch, suitMatch := 0, 0
	for _, v := range game.Deck {
		if v.Suit == c.Suit {
			suitMatch++
		}
		if v.Rank == c.Rank {
			rankMatch++
		}
	}
	assert.Equal(t, 5, rankMatch)
	assert.Equal(t, 14, suitMatch)
}

func TestZanyJoker(t *testing.T) {
	game := balatro.NewBasicGameState()
	game.Jokers = append(game.Jokers, balatro.NewJoker(balatro.ZanyJoker))

	round := game.NextRound()

	handT, s, m := round.ScoreHand(&game, balatro.Hand([]balatro.Card{
		balatro.NewBasicCard(balatro.A, balatro.Diamonds),
		balatro.NewBasicCard(balatro.A, balatro.Spades),
		balatro.NewBasicCard(balatro.A, balatro.Hearts),
		balatro.NewBasicCard(balatro.R4, balatro.Hearts),
		balatro.NewBasicCard(balatro.R6, balatro.Hearts),
	}))

	assert.Equal(t, balatro.ThreeOfAKind, *handT)
	assert.Equal(t, 30+11*3, s)
	assert.Equal(t, 3.0+12, m)
}

func TestWilyJoker(t *testing.T) {
	game := balatro.NewBasicGameState()
	game.Jokers = append(game.Jokers, balatro.NewJoker(balatro.WilyJoker))

	round := game.NextRound()

	handT, s, m := round.ScoreHand(&game, balatro.Hand([]balatro.Card{
		balatro.NewBasicCard(balatro.A, balatro.Diamonds),
		balatro.NewBasicCard(balatro.A, balatro.Spades),
		balatro.NewBasicCard(balatro.A, balatro.Hearts),
		balatro.NewBasicCard(balatro.R4, balatro.Hearts),
		balatro.NewBasicCard(balatro.R6, balatro.Hearts),
	}))

	assert.Equal(t, balatro.ThreeOfAKind, *handT)
	assert.Equal(t, 30+11*3+100, s)
	assert.Equal(t, 3.0, m)
}

func TestFourFingers(t *testing.T) {
	game := balatro.NewBasicGameState()
	j := balatro.NewJoker(balatro.FourFingers)
	game.Jokers = append(game.Jokers, j)
	j.Type.Effects[0].Effect(&game, nil, nil, -1, nil)

	round := game.NextRound()

	handT, _, _ := round.ScoreHand(&game, balatro.Hand([]balatro.Card{
		balatro.NewBasicCard(balatro.A, balatro.Diamonds),
		balatro.NewBasicCard(balatro.R7, balatro.Hearts),
		balatro.NewBasicCard(balatro.K, balatro.Hearts),
		balatro.NewBasicCard(balatro.R4, balatro.Hearts),
		balatro.NewBasicCard(balatro.R6, balatro.Hearts),
	}))

	assert.Equal(t, balatro.Flush, *handT)

	handT, _, _ = round.ScoreHand(&game, balatro.Hand([]balatro.Card{
		balatro.NewBasicCard(balatro.A, balatro.Hearts),
		balatro.NewBasicCard(balatro.R7, balatro.Hearts),
		balatro.NewBasicCard(balatro.K, balatro.Hearts),
		balatro.NewBasicCard(balatro.R4, balatro.Hearts),
		balatro.NewBasicCard(balatro.R6, balatro.Hearts),
	}))

	assert.Equal(t, balatro.Flush, *handT)

	j.Type.Effects[0].RemoveEffect(&game, nil)

	handT, _, _ = round.ScoreHand(&game, balatro.Hand([]balatro.Card{
		balatro.NewBasicCard(balatro.A, balatro.Diamonds),
		balatro.NewBasicCard(balatro.R7, balatro.Hearts),
		balatro.NewBasicCard(balatro.K, balatro.Hearts),
		balatro.NewBasicCard(balatro.R4, balatro.Hearts),
		balatro.NewBasicCard(balatro.R6, balatro.Hearts),
	}))

	assert.NotEqual(t, balatro.Flush, *handT)
}

func TestHack(t *testing.T) {
	game := balatro.NewBasicGameState()
	j := balatro.NewJoker(balatro.Hack)
	game.Jokers = append(game.Jokers, j)

	round := game.NextRound()

	handT, s, m := round.ScoreHand(&game, balatro.Hand([]balatro.Card{
		balatro.NewBasicCard(balatro.A, balatro.Hearts),
		balatro.NewBasicCard(balatro.K, balatro.Hearts),
		balatro.NewBasicCard(balatro.R6, balatro.Hearts),
		balatro.NewCard(balatro.R4, balatro.Hearts, balatro.Normal{}, balatro.NewMult()),
		balatro.NewCard(balatro.R5, balatro.Hearts, balatro.Normal{}, balatro.NewGlass()),
	}))

	assert.Equal(t, balatro.Flush, *handT)
	assert.Equal(t, 35+11+10+6+4*2+5*2, s)
	assert.Equal(t, (4.0+4.0*2)*2*2, m)
}
