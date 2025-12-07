package rand

import (
	"crypto/rand"
	"math/big"
)

func Int(max int) int {
	bigint, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(bigint.Int64())
}
