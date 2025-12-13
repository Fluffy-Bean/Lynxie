package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func NewPassword() string {
	var pass string

	for range 8 {
		// if this panics im going to be impressed
		i, _ := rand.Int(rand.Reader, big.NewInt(2047))
		pass += " " + words[i.Int64()]
	}

	return strings.TrimSpace(pass)
}
