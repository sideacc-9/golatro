package balatro

import (
	"fmt"
	"golatro/pkg/balatro"
	"testing"
)

func TestTargetPointsSequence(t *testing.T) {
	for ante := 1; ante <= 8; ante++ {
		for round := 1; round <= 3; round++ {
			fmt.Println("ante", ante, "round", round, balatro.TargetPoints(round, ante))
		}
	}
}
