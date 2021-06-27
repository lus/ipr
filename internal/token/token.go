package token

import (
	"math/rand"

	"github.com/alexedwards/argon2id"
)

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_#+*")

type Token struct {
	raw  string
	hash string
}

func Generate() *Token {
	return &Token{
		raw: randomString(32),
	}
}

func (token *Token) Raw() string {
	return token.raw
}

func (token *Token) Hash() (string, error) {
	if token.hash != "" {
		return token.hash, nil
	}

	hash, err := argon2id.CreateHash(token.raw, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	token.hash = hash
	return hash, nil
}

func randomString(n int) string {
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = characters[rand.Intn(len(characters))]
	}
	return string(runes)
}
