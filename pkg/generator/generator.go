package generator

import (
	"crypto"

	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
)

const (
	_defaultHashFunc = crypto.SHA256
	_defaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	_defaultLength   = 10
)

type Generator struct {
	hashFunc func() crypto.Hash
	alphabet []rune
	length   int
}

func (g *Generator) GenerateShortURL(url string) string {
	// rand.Seed(time.Now().UnixNano())
	hasher := g.hashFunc().New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)

	var result string
	for i := 0; i < g.length; i++ {
		index := int(hash[i]) % len(g.alphabet)
		result += string(g.alphabet[index])
	}
	return result
}

func NewGenerator(opts ...Option) *Generator {
	g := &Generator{
		hashFunc: _defaultHashFunc.HashFunc, // MD5
		alphabet: []rune(_defaultAlphabet),
		length:   _defaultLength,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}
