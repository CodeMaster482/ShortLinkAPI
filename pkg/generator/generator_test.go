package generator

import (
	"crypto"
	"testing"
)

func TestGenerateShortURL(t *testing.T) {
	t.Parallel()

	g := NewGenerator()
	urls := []struct {
		url          string
		expectedHash string
	}{
		{"https://www.example.com/page1", "VZvG8PSsrp"},
		{"https://www.example.com/page2", "rlhPMEN_Mc"},
		{"https://www.example.com/page3", "AvXq4g85WW"},
		{"lol", "hs_FjJxhv5"},
	}

	for _, tc := range urls {
		shortURL := g.GenerateShortURL(tc.url)
		if shortURL != tc.expectedHash {
			t.Errorf("Shortened URL for %s is %s, but expected %s", tc.url, shortURL, tc.expectedHash)
		}
	}
}

func TestCustomOptions(t *testing.T) {
	t.Parallel()
	g := NewGenerator(
		WithHashFunc(crypto.SHA256),
		WithAlphabet("abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKZXCVBNM"),
		WithLength(5),
	)

	urls := []struct {
		url          string
		expectedHash string
	}{
		{"https://www.example.com/page1", "UagZj"},
		{"https://www.example.com/page2", "HVYcl"},
		{"https://www.example.com/page3", "Qgkbf"},
	}

	for _, tc := range urls {
		shortURL := g.GenerateShortURL(tc.url)
		if shortURL != tc.expectedHash {
			t.Errorf("Shortened URL for %s is %s, but expected %s", tc.url, shortURL, tc.expectedHash)
		}
	}
}
