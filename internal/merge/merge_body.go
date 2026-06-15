package merge

import (
	"math"
	"strings"
)

const similarityThreshold = 0.70

// mergeBody combines original and modified body lines. It preserves lines that
// exist in both documents at the same position, replaces original lines with
// similar modified lines using a fuzzy dice coefficient match, and appends new
// modified lines that have no match.
func mergeBody(origLines, modLines []string) []string {
	if len(origLines) == 0 && len(modLines) == 0 {
		return []string{}
	}
	if len(origLines) == 0 {
		return append([]string(nil), modLines...)
	}
	if len(modLines) == 0 {
		return []string{}
	}

	origMatched := make([]bool, len(origLines))
	modMatched := make([]bool, len(modLines))

	minLen := int(math.Min(float64(len(origLines)), float64(len(modLines))))
	for i := 0; i < minLen; i++ {
		origMatched[i] = true
		modMatched[i] = true
	}

	var result []string
	for i := 0; i < len(modLines); i++ {
		if modMatched[i] {
			result = append(result, modLines[i])
			continue
		}

		bestIdx := -1
		bestScore := 0.0
		for j, origLine := range origLines {
			if origMatched[j] {
				continue
			}
			score := dice(modLines[i], origLine)
			if score > bestScore {
				bestScore = score
				bestIdx = j
			}
		}

		if bestScore >= similarityThreshold {
			origMatched[bestIdx] = true
		}

		result = append(result, modLines[i])
	}

	return result
}

// dice returns the Dice coefficient between two strings based on character
// bigrams. Identical strings (after normalization) return 1.0.
func dice(a, b string) float64 {
	if a == b {
		return 1.0
	}

	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	if a == "" && b == "" {
		return 1.0
	}
	if a == "" || b == "" {
		return 0.0
	}

	bigramsA := bigrams(a)
	bigramsB := bigrams(b)

	intersection := 0
	for bg, countA := range bigramsA {
		if countB, ok := bigramsB[bg]; ok {
			if countA < countB {
				intersection += countA
			} else {
				intersection += countB
			}
		}
	}

	total := len(bigramsA) + len(bigramsB)
	if total == 0 {
		return 0.0
	}

	return 2.0 * float64(intersection) / float64(total)
}

// bigrams returns a multiset of overlapping 2-character substrings for s.
func bigrams(s string) map[string]int {
	result := make(map[string]int)
	for i := 0; i+1 < len(s); i++ {
		result[s[i:i+2]]++
	}
	return result
}
