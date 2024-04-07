package generator

import (
	"crypto"
)

type Option func(*Generator)

func WithAlphabet(alphabet string) Option {
	return func(g *Generator) {
		g.alphabet = []rune(alphabet)
	}
}

func WithHashFunc(hash crypto.Hash) Option {
	return func(g *Generator) {
		g.hashFunc = hash.HashFunc
	}
}

func WithLength(length int) Option {
	return func(g *Generator) {
		g.length = length
	}
}
