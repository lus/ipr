package token

import (
	"math/rand"

	"github.com/alexedwards/argon2id"
)

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_#+*")

// Generate generates a new machine token
func Generate() string {
	return randomString(64)
}

// Hash hashes a machine token
func Hash(tkn string) (string, error) {
	return argon2id.CreateHash(tkn, argon2id.DefaultParams)
}

// Check checks if a given raw string matches a specific machine token
func Check(hash, raw string) (bool, error) {
	return argon2id.ComparePasswordAndHash(raw, hash)
}

func randomString(n int) string {
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = characters[rand.Intn(len(characters))]
	}
	return string(runes)
}
