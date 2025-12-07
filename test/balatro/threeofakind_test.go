package balatro

import (
	"golatro/pkg/balatro"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTOAKIncorrectLength(t *testing.T) {
	hand := make(balatro.Hand, 0, 5)
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	_, ok := balatro.CheckHandType(balatro.ThreeOfAKind, hand)
	assert.Equal(t, false, ok)
}

func TestTOAKIncorrectHandType(t *testing.T) {
	hand := make(balatro.Hand, 0, 5)
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.K, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.Q, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.R2, balatro.Clubs))
	_, ok := balatro.CheckHandType(balatro.ThreeOfAKind, hand)
	assert.Equal(t, false, ok)
}

func TestTOAKCorrect(t *testing.T) {
	hand := make(balatro.Hand, 0, 5)
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.K, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.A, balatro.Clubs))
	hand = append(hand, balatro.NewBasicCard(balatro.R2, balatro.Clubs))
	_, ok := balatro.CheckHandType(balatro.ThreeOfAKind, hand)
	assert.Equal(t, true, ok)
}
