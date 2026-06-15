package renderer

import (
	"testing"

	"github.com/blak0p/splice/internal/ast"
	"github.com/blak0p/splice/internal/parser"
)

func TestRenderRoundTrip(t *testing.T) {
	input := "# Intro\n\nintro body\n\n## Details\n\ndetails body\n"
	doc, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	rendered := Render(doc)
	doc2, err := parser.Parse(rendered)
	if err != nil {
		t.Fatalf("parse rendered: %v", err)
	}

	if len(doc2.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(doc2.Sections))
	}
	assertSection(t, doc2.Sections[0], 1, "intro", "intro body")
	assertSection(t, doc2.Sections[1], 2, "details", "details body")
}

func TestRenderPreHeadingThenHeadings(t *testing.T) {
	doc := &ast.Document{
		Sections: []ast.Section{
			{Heading: nil, Body: ast.Body{Content: "pre-heading content"}},
			{Heading: &ast.Heading{Level: 1, Text: "First"}, Body: ast.Body{Content: "body"}},
		},
	}

	got := Render(doc)
	want := "pre-heading content\n\n# First\nbody\n"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderMultiSectionPreservesBoundaries(t *testing.T) {
	doc := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "One"}, Body: ast.Body{Content: "body one"}},
			{Heading: &ast.Heading{Level: 1, Text: "Two"}, Body: ast.Body{Content: "body two"}},
			{Heading: &ast.Heading{Level: 1, Text: "Three"}, Body: ast.Body{Content: "body three"}},
		},
	}

	got := Render(doc)
	want := "# One\nbody one\n\n# Two\nbody two\n\n# Three\nbody three\n"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func assertSection(t *testing.T, s ast.Section, level int, text, body string) {
	t.Helper()
	if s.Heading == nil {
		t.Fatalf("expected heading, got nil")
	}
	if s.Heading.Level != level {
		t.Fatalf("expected level %d, got %d", level, s.Heading.Level)
	}
	if s.Heading.Text != text {
		t.Fatalf("expected text %q, got %q", text, s.Heading.Text)
	}
	if s.Body.Content != body {
		t.Fatalf("expected body %q, got %q", body, s.Body.Content)
	}
}
