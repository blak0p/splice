package merge

import (
	"strings"

	"github.com/blak0p/splice/ast"
)

const similarityThreshold = 0.70

// mergeBody combines original and modified AST body blocks.
func mergeBody(origBlocks, modBlocks []ast.Block) []ast.Block {
	var merged []ast.Block
	for i := 0; i < len(modBlocks); i++ {
		if i < len(origBlocks) {
			orig, mod := origBlocks[i], modBlocks[i]
	if orig.Kind() == mod.Kind() {
		switch orig.Kind() {
		case ast.KindParagraph:
			merged = append(merged, ast.Paragraph{ContentLines: mergeLines(orig.Lines(), mod.Lines())})
		case ast.KindList:
			merged = append(merged, ast.List{ContentLines: mergeLines(orig.Lines(), mod.Lines())})
		default:
			merged = append(merged, mod) // Atomic Table/CodeBlock merge
		}
	} else {
		merged = append(merged, mod) // Kind mismatch: atomic replace
	}
		} else {
			merged = append(merged, modBlocks[i]) // Append new blocks
		}
	}
	return merged
}

// mergeLines combines original and modified body lines. It preserves lines that
// exist in both documents at the same position, replaces original lines with
// similar modified lines using a fuzzy dice coefficient match, and appends new
// modified lines that have no match.
func mergeLines(origLines, modLines []string) []string {
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

	minLen := len(origLines)
	if len(modLines) < minLen {
		minLen = len(modLines)
	}
	for i := 0; i < minLen; i++ {
		if origLines[i] == modLines[i] {
			origMatched[i] = true
			modMatched[i] = true
		}
	}

	for i := 0; i < len(modLines); i++ {
		if modMatched[i] {
			continue
		}
		for j := 0; j < len(origLines); j++ {
			if origMatched[j] {
				continue
			}
			if modLines[i] == origLines[j] {
				origMatched[j] = true
				modMatched[i] = true
				break
			}
		}
	}

	for i := 0; i < len(modLines); i++ {
		if modMatched[i] {
			continue
		}
		bestIdx := -1
		bestScore := -1.0
		for j := 0; j < len(origLines); j++ {
			if origMatched[j] {
				continue
			}
			score := dice(modLines[i], origLines[j])
			if score > bestScore {
				bestScore = score
				bestIdx = j
			}
		}
		if bestIdx != -1 && bestScore >= similarityThreshold {
			origMatched[bestIdx] = true
			modMatched[i] = true
		}
	}

	result := make([]string, len(modLines))
	copy(result, modLines)
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
