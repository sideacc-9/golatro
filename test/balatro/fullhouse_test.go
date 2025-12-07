package balatro

import (
	"golatro/pkg/balatro"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullHouseIncorrectLength(t *testing.T) {
	hand := make(balatro.Hand, 0, 5)
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	_, ok := balatro.CheckHandType(balatro.FullHouse, hand)
	assert.Equal(t, false, ok)
}

func TestFullHouseIncorrectHandType(t *testing.T) {
	hand := make(balatro.Hand, 0, 5)
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.J, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.K, balatro.Clubs))
	_, ok := balatro.CheckHandType(balatro.FullHouse, hand)
	assert.Equal(t, false, ok)
}

func TestFullHouseCorrect(t *testing.T) {
	hand := make(balatro.Hand, 0, 5)
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.J, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.J, balatro.Clubs))
	_, ok := balatro.CheckHandType(balatro.FullHouse, hand)
	assert.Equal(t, true, ok)
}
