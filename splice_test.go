package splice_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/blak0p/splice"
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

// TestMerge_RealWorld performs a real-world merge using AGENTS.md files.
func TestMerge_RealWorld(t *testing.T) {
	ctx := context.Background()

	bakContent, err := os.ReadFile("/home/alejandro/dev/.entorno/splice/.config/opencode/AGENTS.md.bak")
	if err != nil {
		t.Fatalf("failed to read AGENTS.md.bak: %v", err)
	}

	mdContent, err := os.ReadFile("/home/alejandro/dev/.entorno/splice/.config/opencode/AGENTS.md")
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}

	result, err := splice.Merge(ctx, string(bakContent), string(mdContent))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == "" {
		t.Fatal("expected non-empty merged output")
	}

	mergedPath := "testdata/merged-agents.md"
	if err := os.MkdirAll("testdata", 0o755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}
	if err := os.WriteFile(mergedPath, []byte(result), 0o644); err != nil {
		t.Fatalf("failed to write merged output: %v", err)
	}
}
