package merge

import (
	"testing"

	"github.com/blak0p/splice/internal/ast"
)

func TestMergeNormalizedHeadingMatch(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "intro"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"original body"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "INTRO"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"modified body"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(got.Sections))
	}
	if got.Sections[0].Body.Lines()[0] != "modified body" {
		t.Fatalf("expected modified body, got %q", got.Sections[0].Body.Lines())
	}
}

func TestMergePreservesOriginalOnlySections(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "keep"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"original body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "remove"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"gone"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "keep"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"modified body"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	assertSectionHeading(t, got.Sections[0], "keep")
	assertSectionHeading(t, got.Sections[1], "remove")
	if got.Sections[0].Body.Lines()[0] != "modified body" {
		t.Fatalf("expected modified body for keep, got %q", got.Sections[0].Body.Lines())
	}
	if got.Sections[1].Body.Lines()[0] != "gone" {
		t.Fatalf("expected original body for remove, got %q", got.Sections[1].Body.Lines())
	}
}

func TestMergeModifiedBodyReplacesOriginal(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"a original"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "b"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"b original"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"a modified"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "b"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"b modified"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	if got.Sections[0].Body.Lines()[0] != "a modified" {
		t.Fatalf("expected a modified, got %q", got.Sections[0].Body.Lines())
	}
	if got.Sections[1].Body.Lines()[0] != "b modified" {
		t.Fatalf("expected b modified, got %q", got.Sections[1].Body.Lines())
	}
}

func TestMergeAddsNewSectionsFromModified(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "first"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"first body"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "first"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"first body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "second"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"second body"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	assertSectionHeading(t, got.Sections[0], "first")
	assertSectionHeading(t, got.Sections[1], "second")
}

func TestMergePreHeadingContent(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: nil, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"original pre"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"a body"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: nil, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"modified pre"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"a modified"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	if got.Sections[0].Heading != nil {
		t.Fatalf("expected implicit section first")
	}
	if got.Sections[0].Body.Lines()[0] != "modified pre" {
		t.Fatalf("expected modified pre, got %q", got.Sections[0].Body.Lines())
	}
	if got.Sections[1].Body.Lines()[0] != "a modified" {
		t.Fatalf("expected a modified, got %q", got.Sections[1].Body.Lines())
	}
}

func TestMergeDuplicateHeadingsByPosition(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"dup1 original"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"dup2 original"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"dup1 modified"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"dup2 modified"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	if got.Sections[0].Body.Lines()[0] != "dup1 modified" {
		t.Fatalf("expected dup1 modified, got %q", got.Sections[0].Body.Lines())
	}
	if got.Sections[1].Body.Lines()[0] != "dup2 modified" {
		t.Fatalf("expected dup2 modified, got %q", got.Sections[1].Body.Lines())
	}
}

func TestMergeEmptyDocuments(t *testing.T) {
	original := &ast.Document{}
	modified := &ast.Document{}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 0 {
		t.Fatalf("expected 0 sections, got %d", len(got.Sections))
	}
}

func assertSectionHeading(t *testing.T, s ast.Section, want string) {
	t.Helper()
	if s.Heading == nil {
		t.Fatalf("expected heading %q, got nil", want)
	}
	if s.Heading.Text != want {
		t.Fatalf("expected heading %q, got %q", want, s.Heading.Text)
	}
}
