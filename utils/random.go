package utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int64) string {
	var sb strings.Builder

	for i := int64(0); i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(len(alphabet))])

	}
	return sb.String()
}

func RandomInt63(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomOwnerName(min, max int64) string {
	return RandomString(RandomInt63(min, max))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "JPY", "GBP"}
	return currencies[RandomInt63(0, int64(len(currencies)-1))]
}

