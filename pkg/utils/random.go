package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomCode generates a random numeric code of specified length
func GenerateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())

	code := ""
	for i := 0; i < length; i++ {
		code += string(rune('0' + rand.Intn(10)))
	}

	return code
}
