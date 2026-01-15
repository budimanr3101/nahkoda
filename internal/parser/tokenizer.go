package parser

import (
	"errors"
	"strings"
	"unicode"
)

// tokenize converts a string into a list of tokens, respecting quotes.
// Example: `liat "pod keren"` -> ["liat", "pod keren"]
func tokenize(input string) ([]string, error) {
	var tokens []string
	var currentToken strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range input {
		switch {
		case inQuote:
			if r == quoteChar {
				inQuote = false
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			} else {
				currentToken.WriteRune(r)
			}
		case unicode.IsSpace(r):
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		case r == '"' || r == '\'':
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			inQuote = true
			quoteChar = r
		default:
			currentToken.WriteRune(r)
		}
	}

	if inQuote {
		return nil, errors.New("unclosed quote")
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens, nil
}
