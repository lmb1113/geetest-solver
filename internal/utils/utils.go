package utils

import (
	"math/rand"
)

func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
