package merge

import (
	"testing"

	"github.com/blak0p/splice/ast"
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

func TestMerge_InsertRelativeToDeletedSections(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "A"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"A body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "B"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"B body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "C"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"C body"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "A"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"A body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "X"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"X body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "C"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"C body"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 4 {
		t.Fatalf("expected 4 sections, got %d", len(got.Sections))
	}
	assertSectionHeading(t, got.Sections[0], "A")
	assertSectionHeading(t, got.Sections[1], "X")
	assertSectionHeading(t, got.Sections[2], "B")
	assertSectionHeading(t, got.Sections[3], "C")
}

func TestMerge_MultipleSequentialInsertions(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "A"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"A body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "B"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"B body"}}}}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "A"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"A body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "X"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"X body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "Y"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"Y body"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "B"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"B body"}}}}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 4 {
		t.Fatalf("expected 4 sections, got %d", len(got.Sections))
	}
	assertSectionHeading(t, got.Sections[0], "A")
	assertSectionHeading(t, got.Sections[1], "X")
	assertSectionHeading(t, got.Sections[2], "Y")
	assertSectionHeading(t, got.Sections[3], "B")
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

