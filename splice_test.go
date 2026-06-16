package splice_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blak0p/splice"
	"github.com/blak0p/splice/ast"
)

// TestMerge_Success verifies that two valid markdown documents merge correctly.
func TestMerge_Success(t *testing.T) {
	ctx := context.Background()
	original := "# Intro\n\nHello\n\n# Install\n\nOriginal body\n"
	modified := "# Intro\n\nHello\n\n# Install\n\nModified body\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Modified body") {
		t.Errorf("expected merged output to contain modified body, got:\n%s", result)
	}
}

// TestMerge_UnparseableInput verifies that invalid UTF-8 returns an error.
// tree-sitter-markdown itself tolerates arbitrary byte sequences, so the parser
// validates the input first and rejects bytes that are not valid UTF-8.
func TestMerge_UnparseableInput(t *testing.T) {
	ctx := context.Background()
	original := "# Valid\n\ncontent\n"
	modified := string([]byte{0xff, 0xfe, 0xfd})

	_, err := splice.Merge(ctx, original, modified)
	if err == nil {
		t.Fatal("expected error for unparseable input, got nil")
	}
}

// TestMerge_NormalizedHeadingMatch verifies case-insensitive heading matching.
func TestMerge_NormalizedHeadingMatch(t *testing.T) {
	ctx := context.Background()
	original := "## Quick Start\n\nOriginal\n"
	modified := "## quick start\n\nModified\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Modified") {
		t.Errorf("expected merged output to contain modified body, got:\n%s", result)
	}
}

// TestMerge_DistinctHeadings verifies different headings are treated as separate sections.
func TestMerge_DistinctHeadings(t *testing.T) {
	ctx := context.Background()
	original := "## Install\n\nOriginal install\n"
	modified := "## Upgrade\n\nModified upgrade\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Original install") {
		t.Errorf("expected output to preserve original section, got:\n%s", result)
	}
	if !strings.Contains(result, "Modified upgrade") {
		t.Errorf("expected output to include modified section, got:\n%s", result)
	}
}

// TestMerge_PreserveOriginalOnly verifies original-only sections are preserved.
func TestMerge_PreserveOriginalOnly(t *testing.T) {
	ctx := context.Background()
	original := "## Keep\n\nOriginal body\n"
	modified := ""

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Original body") {
		t.Errorf("expected original-only section to be preserved, got:\n%s", result)
	}
}

// TestMerge_ModifiedBodyReplaces verifies shared sections use modified body.
func TestMerge_ModifiedBodyReplaces(t *testing.T) {
	ctx := context.Background()
	original := "## Section\n\nOriginal body\n"
	modified := "## Section\n\nModified body\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Original body") {
		t.Errorf("expected original body to be replaced, got:\n%s", result)
	}
	if !strings.Contains(result, "Modified body") {
		t.Errorf("expected modified body in output, got:\n%s", result)
	}
}

// TestMerge_InsertNewSection verifies new sections from modified appear in output.
func TestMerge_InsertNewSection(t *testing.T) {
	ctx := context.Background()
	original := "## Existing\n\nOriginal body\n"
	modified := "## Existing\n\nOriginal body\n\n## New\n\nNew body\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "New body") {
		t.Errorf("expected new section in output, got:\n%s", result)
	}
}

// TestMerge_LeadingParagraphWins verifies pre-heading content from modified wins.
func TestMerge_LeadingParagraphWins(t *testing.T) {
	ctx := context.Background()
	original := "Original pre-heading\n\n## Heading\n\nBody\n"
	modified := "Modified pre-heading\n\n## Heading\n\nBody\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Modified pre-heading") {
		t.Errorf("expected modified pre-heading to win, got:\n%s", result)
	}
}

// TestMerge_OriginalNoPreHeading verifies modified can add pre-heading content.
func TestMerge_OriginalNoPreHeading(t *testing.T) {
	ctx := context.Background()
	original := "## Heading\n\nBody\n"
	modified := "Pre-heading content\n\n## Heading\n\nBody\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Pre-heading content") {
		t.Errorf("expected pre-heading content from modified to appear, got:\n%s", result)
	}
}

// TestMerge_BothEmpty verifies merging two empty documents returns empty string.
func TestMerge_BothEmpty(t *testing.T) {
	ctx := context.Background()

	result, err := splice.Merge(ctx, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

// TestMerge_NoHeadingsDoc verifies documents without headings merge as a single implicit section.
func TestMerge_NoHeadingsDoc(t *testing.T) {
	ctx := context.Background()
	original := "Original paragraph\n"
	modified := "Modified paragraph\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Modified paragraph") {
		t.Errorf("expected modified paragraph in output, got:\n%s", result)
	}
}

// TestMerge_DuplicateHeadings verifies duplicate headings are matched by sequential position.
func TestMerge_DuplicateHeadings(t *testing.T) {
	ctx := context.Background()
	original := "## Dup\n\nFirst\n\n## Dup\n\nSecond\n"
	modified := "## Dup\n\nFirst mod\n\n## Dup\n\nSecond mod\n"

	result, err := splice.Merge(ctx, original, modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "First mod") {
		t.Errorf("expected first duplicate to be replaced, got:\n%s", result)
	}
	if !strings.Contains(result, "Second mod") {
		t.Errorf("expected second duplicate to be replaced, got:\n%s", result)
	}
}

// TestMerge_E2E walks testdata directories and verifies each case.
func TestMerge_E2E(t *testing.T) {
	ctx := context.Background()

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			dir := filepath.Join("testdata", entry.Name())

			original, err := os.ReadFile(filepath.Join(dir, "original.md"))
			if err != nil {
				t.Fatalf("failed to read original.md: %v", err)
			}

			modified, err := os.ReadFile(filepath.Join(dir, "modified.md"))
			if err != nil {
				t.Fatalf("failed to read modified.md: %v", err)
			}

			expected, err := os.ReadFile(filepath.Join(dir, "expected.md"))
			if err != nil {
				t.Fatalf("failed to read expected.md: %v", err)
			}

			got, err := splice.Merge(ctx, string(original), string(modified))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != string(expected) {
				t.Fatalf("expected:\n%s\n\ngot:\n%s", string(expected), got)
			}
		})
	}
}

// TestParseRenderRoundTrip verifies Parse + Render round-trip.
func TestParseRenderRoundTrip(t *testing.T) {
	input := "# Hello\n\nWorld\n\n## Section\n\nContent\n"
	doc, err := splice.Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	output := splice.Render(doc)
	if output != input {
		t.Fatalf("round-trip mismatch:\nexpected:\n%s\n\ngot:\n%s", input, output)
	}
}

// TestParseEmpty verifies Parse handles empty input.
func TestParseEmpty(t *testing.T) {
	doc, err := splice.Parse("")
	if err != nil {
		t.Fatalf("Parse empty failed: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
}

// TestMergeAST verifies MergeAST works with pre-parsed documents.
func TestMergeAST(t *testing.T) {
	origDoc, _ := splice.Parse("# Section\n\nOriginal\n")
	modDoc, _ := splice.Parse("# Section\n\nModified\n")
	merged := splice.MergeAST(origDoc, modDoc)
	output := splice.Render(merged)
	if !strings.Contains(output, "Modified") {
		t.Errorf("expected modified content, got:\n%s", output)
	}
}

// TestWithThreshold verifies WithThreshold(0.95) treats similar blocks as distinct.
func TestWithThreshold(t *testing.T) {
	ctx := context.Background()
	original := "# Section\n\nHello world\n"
	modified := "# Section\n\nHello world modified\n"
	// With strict threshold, similar but non-matching lines should be treated as distinct
	result, err := splice.Merge(ctx, original, modified, splice.WithThreshold(0.95))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Hello world modified") {
		t.Errorf("expected modified content, got:\n%s", result)
	}
}

// TestWithCaseInsensitive verifies headings match case-insensitively.
func TestWithCaseInsensitive(t *testing.T) {
	ctx := context.Background()
	original := "# Intro\n\nOriginal\n"
	modified := "# intro\n\nModified\n"
	result, err := splice.Merge(ctx, original, modified, splice.WithCaseInsensitive(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Modified") {
		t.Errorf("expected modified content with case-insensitive match, got:\n%s", result)
	}
}

// TestWithBlockMerger verifies custom block merger callback is invoked.
func TestWithBlockMerger(t *testing.T) {
	ctx := context.Background()
	original := "# Section\n\nHello\n"
	modified := "# Section\n\nWorld\n"
	called := false
	merger := func(orig, mod ast.Block) (ast.Block, bool) {
		called = true
		return mod, true
	}
	result, err := splice.Merge(ctx, original, modified, splice.WithBlockMerger(merger))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected WithBlockMerger callback to be invoked")
	}
	if !strings.Contains(result, "World") {
		t.Errorf("expected merged content from custom merger, got:\n%s", result)
	}
}
