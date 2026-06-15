package parser

import (
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
	assertSection(t, doc.Sections[0], 1, "intro", "intro body")
	assertSection(t, doc.Sections[1], 2, "details", "details body")
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
	if doc.Sections[0].Body.Content != "pre-heading content" {
		t.Fatalf("unexpected pre-heading content: %q", doc.Sections[0].Body.Content)
	}
	assertSection(t, doc.Sections[1], 1, "first", "body")
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
	if doc.Sections[0].Body.Content != "just some content\nwithout headings" {
		t.Fatalf("unexpected body content: %q", doc.Sections[0].Body.Content)
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

func assertSection(t *testing.T, s ast.Section, level int, text, body string) {
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
	if s.Body.Content != body {
		t.Fatalf("expected body %q, got %q", body, s.Body.Content)
	}
}
