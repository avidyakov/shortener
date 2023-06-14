package utils

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortID(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)

	for i, value := range randomBytes {
		randomBytes[i] = charset[value%byte(len(charset))]
	}
	return string(randomBytes)
}
