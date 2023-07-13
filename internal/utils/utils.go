package utils

import (
	"crypto/rand"
	"errors"
	"net/url"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var errInvalidScheme = errors.New("invalid scheme")

func GenerateShortID(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)

	for i, value := range randomBytes {
		randomBytes[i] = charset[value%byte(len(charset))]
	}
	return string(randomBytes)
}

func ValidateLink(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", errInvalidScheme
	}

	return link, nil
}
