package merge

import (
	"strings"

	"github.com/blak0p/splice/ast"
)

// Config controls merge behavior.
type Config struct {
	Threshold       float64
	CaseInsensitive bool
	BlockMerger     func(orig, mod ast.Block) (ast.Block, bool)
}

// DefaultConfig returns a Config with the standard similarity threshold.
func DefaultConfig() *Config {
	return &Config{Threshold: 0.7, CaseInsensitive: true}
}

// MergeAST combines original and modified documents. Sections that exist in both
// documents keep the modified body; sections only in original are preserved;
// sections only in modified are inserted after their nearest matched predecessor.
func MergeAST(original, modified *ast.Document, cfg *Config) *ast.Document {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if original == nil && modified == nil {
		return &ast.Document{}
	}
	if modified == nil {
		return shallowCopy(original)
	}
	if original == nil {
		return shallowCopy(modified)
	}

	matched := make([]bool, len(modified.Sections))
	modIdxToMergedIdx := make(map[int]int)
	var merged []ast.Section

	for i, origSection := range original.Sections {
		matchIdx := findMatch(origSection, i, original.Sections, modified.Sections, matched, cfg)
		if matchIdx == -1 {
			merged = append(merged, origSection)
			continue
		}

		modSection := modified.Sections[matchIdx]
		merged = append(merged, ast.Section{
			Heading: modSection.Heading,
			Body:    ast.Body{Blocks: mergeBody(origSection.Body.Blocks, modSection.Body.Blocks, cfg)},
		})
		modIdxToMergedIdx[matchIdx] = len(merged) - 1
		matched[matchIdx] = true
	}

	for i, modSection := range modified.Sections {
		if matched[i] {
			continue
		}

		insertIdx := findInsertIndex(i, matched, modIdxToMergedIdx)
		if insertIdx < 0 {
			merged = append(merged, modSection)
			modIdxToMergedIdx[i] = len(merged) - 1
		} else {
			for k, v := range modIdxToMergedIdx {
				if v >= insertIdx {
					modIdxToMergedIdx[k] = v + 1
				}
			}
			merged = append(merged, modSection)
			copy(merged[insertIdx+1:], merged[insertIdx:])
			merged[insertIdx] = modSection
			modIdxToMergedIdx[i] = insertIdx
		}
		matched[i] = true
	}

	return &ast.Document{Sections: merged}
}

// findMatch locates the section in modified that corresponds to original. It
// first tries exact heading match with positional bookkeeping, then falls back
// to the first unmatched implicit section for pre-heading content.
func findMatch(origSection ast.Section, origIndex int, originalSections, modifiedSections []ast.Section, matched []bool, cfg *Config) int {
	if origSection.Heading == nil {
		return findUnmatchedImplicit(modifiedSections, matched)
	}

	occurrence := 0
	for j := 0; j <= origIndex; j++ {
		if originalSections[j].Heading != nil &&
			headingEqual(originalSections[j].Heading.Text, origSection.Heading.Text, cfg) {
			occurrence++
		}
	}

	count := 0
	for i, ms := range modifiedSections {
		if ms.Heading == nil {
			continue
		}
		if headingEqual(ms.Heading.Text, origSection.Heading.Text, cfg) {
			count++
			if count == occurrence && !matched[i] {
				return i
			}
		}
	}

	return -1
}

func headingEqual(a, b string, cfg *Config) bool {
	if cfg != nil && cfg.CaseInsensitive {
		return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
	}
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}

func findUnmatchedImplicit(modifiedSections []ast.Section, matched []bool) int {
	for i, ms := range modifiedSections {
		if ms.Heading == nil && !matched[i] {
			return i
		}
	}
	return -1
}

func findInsertIndex(modIndex int, matched []bool, modIdxToMergedIdx map[int]int) int {
	for i := modIndex - 1; i >= 0; i-- {
		if matched[i] {
			if mergedIdx, ok := modIdxToMergedIdx[i]; ok {
				return mergedIdx + 1
			}
		}
	}
	return -1
}

func shallowCopy(doc *ast.Document) *ast.Document {
	if doc == nil {
		return &ast.Document{}
	}
	sections := make([]ast.Section, len(doc.Sections))
	copy(sections, doc.Sections)
	return &ast.Document{Sections: sections}
}
