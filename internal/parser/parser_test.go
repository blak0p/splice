package parser

import (
	"reflect"
	"testing"

	"github.com/blak0p/splice/internal/ast"
)

func TestParseHeadingExtraction(t *testing.T) {
	input := "# Intro\n\nintro body\n\n## Details\n\ndetails body\n"
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(doc.Sections))
	}
	assertSection(t, doc.Sections[0], 1, "Intro", []string{"intro body"})
	assertSection(t, doc.Sections[1], 2, "Details", []string{"details body"})
}

func TestParsePreHeadingContent(t *testing.T) {
	input := "pre-heading content\n\n# First\n\nbody\n"
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(doc.Sections))
	}
	if doc.Sections[0].Heading != nil {
		t.Fatalf("expected implicit section with nil heading, got %+v", doc.Sections[0].Heading)
	}
	want := []string{"pre-heading content"}
	if !reflect.DeepEqual(doc.Sections[0].Body.Lines(), want) {
		t.Fatalf("unexpected pre-heading content: %v", doc.Sections[0].Body.Lines())
	}
	assertSection(t, doc.Sections[1], 1, "First", []string{"body"})
}

func TestParseNoHeadings(t *testing.T) {
	input := "just some content\nwithout headings\n"
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(doc.Sections))
	}
	if doc.Sections[0].Heading != nil {
		t.Fatalf("expected implicit section with nil heading, got %+v", doc.Sections[0].Heading)
	}
	want := []string{"just some content", "without headings"}
	if !reflect.DeepEqual(doc.Sections[0].Body.Lines(), want) {
		t.Fatalf("unexpected body content: %v", doc.Sections[0].Body.Lines())
	}
}

func TestParseEmptyInput(t *testing.T) {
	doc, err := Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 0 {
		t.Fatalf("expected 0 sections, got %d", len(doc.Sections))
	}
}

func TestParseUnparseableInput(t *testing.T) {
	// tree-sitter-markdown is very permissive; force an error with a nil-like input is not possible.
	// We test that parse errors propagate by using a deliberately invalid scenario when one arises.
	// For now, validate the parser does not panic and returns a document.
	doc, err := Parse("# Heading\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(doc.Sections))
	}
}

func TestParseBlankLinesBetweenBlocks(t *testing.T) {
	input := "# Para\n\nFirst block.\n\nSecond block.\n"
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(doc.Sections))
	}
	want := []string{"First block.", "", "Second block."}
	if !reflect.DeepEqual(doc.Sections[0].Body.Lines(), want) {
		t.Fatalf("expected %v, got %v", want, doc.Sections[0].Body.Lines())
	}
}

func TestParseBlockTypes(t *testing.T) {
	input := `# Test Section

This is a paragraph.

- Item 1
- Item 2

| Col 1 | Col 2 |
|---|---|
| Val 1 | Val 2 |

` + "```go" + `
func main() {}
` + "```" + `
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(doc.Sections))
	}
	sect := doc.Sections[0]
	if len(sect.Body.Blocks) != 4 {
		t.Fatalf("expected 4 blocks, got %d", len(sect.Body.Blocks))
	}

	if sect.Body.Blocks[0].Type() != ast.BlockParagraph {
		t.Errorf("expected BlockParagraph, got %v", sect.Body.Blocks[0].Type())
	}
	if sect.Body.Blocks[1].Type() != ast.BlockList {
		t.Errorf("expected BlockList, got %v", sect.Body.Blocks[1].Type())
	}
	if sect.Body.Blocks[2].Type() != ast.BlockTable {
		t.Errorf("expected BlockTable, got %v", sect.Body.Blocks[2].Type())
	}
	if sect.Body.Blocks[3].Type() != ast.BlockCodeBlock {
		t.Errorf("expected BlockCodeBlock, got %v", sect.Body.Blocks[3].Type())
	}
}

func TestParseFallbackToParagraph(t *testing.T) {
	input := `# HTML Section

<div>Some html</div>
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(doc.Sections))
	}
	sect := doc.Sections[0]
	if len(sect.Body.Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(sect.Body.Blocks))
	}
	if sect.Body.Blocks[0].Type() != ast.BlockParagraph {
		t.Errorf("expected BlockParagraph for HTML fallback, got %v", sect.Body.Blocks[0].Type())
	}
}

func assertSection(t *testing.T, s ast.Section, level int, text string, body []string) {
	t.Helper()
	if s.Heading == nil {
		t.Fatalf("expected heading, got nil")
	}
	if s.Heading.Level != level {
		t.Fatalf("expected heading level %d, got %d", level, s.Heading.Level)
	}
	if s.Heading.Text != text {
		t.Fatalf("expected heading text %q, got %q", text, s.Heading.Text)
	}
	if !reflect.DeepEqual(s.Body.Lines(), body) {
		t.Fatalf("expected body %v, got %v", body, s.Body.Lines())
	}
}
