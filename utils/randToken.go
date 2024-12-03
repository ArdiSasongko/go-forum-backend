package utils

import (
	"crypto/rand"
	"math/big"
)

func GenToken() int {
	max := big.NewInt(999999 - 100000 + 1)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0
	}
	return int(n.Int64() + 100000)
}
