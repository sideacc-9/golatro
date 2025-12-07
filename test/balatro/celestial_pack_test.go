package balatro

import (
	"golatro/pkg/balatro"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelestialUpgrade(t *testing.T) {
	game := balatro.NewBasicGameState()
	pairLevel := game.HandLevels[balatro.Pair]
	sum, mult := pairLevel.GetSumMult()
	assert.Equal(t, 1, pairLevel.Level)
	assert.Equal(t, 10, sum)
	assert.Equal(t, 2.0, mult)

	upgradePairCard := balatro.CelestialCard(balatro.Pair)
	err := upgradePairCard.Apply(&game)
	assert.NoError(t, err)

	pairLevel = game.HandLevels[balatro.Pair]
	sum, mult = pairLevel.GetSumMult()
	assert.Equal(t, 2, pairLevel.Level)
	assert.Equal(t, 25, sum)
	assert.Equal(t, 3.0, mult)
}
