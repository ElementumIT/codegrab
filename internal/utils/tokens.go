package utils

import "strings"

// EstimateTokens provides a rough estimate of the number of tokens in the given text.
// This is based on the heuristic that one token is approximately four characters.
func EstimateTokens(text string) int {
	if len(text) == 0 {
		return 0
	}
	cleanText := strings.ReplaceAll(text, "\n", " ")
	return len(cleanText) / 4
}
