package utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
func RandomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(len(alphabet))])

	}
	return sb.String()
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomOwnerName(min, max int) string {
	return RandomString(RandomInt(min, max))
}

func RandomCurrency(){

}