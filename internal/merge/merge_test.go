package merge

import (
	"testing"

	"github.com/blak0p/splice/internal/ast"
)

func TestMergeNormalizedHeadingMatch(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "intro"}, Body: ast.Body{Content: "original body"}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "INTRO"}, Body: ast.Body{Content: "modified body"}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(got.Sections))
	}
	if got.Sections[0].Body.Content != "modified body" {
		t.Fatalf("expected modified body, got %q", got.Sections[0].Body.Content)
	}
}

func TestMergePreservesOriginalOnlySections(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "keep"}, Body: ast.Body{Content: "original body"}},
			{Heading: &ast.Heading{Level: 1, Text: "remove"}, Body: ast.Body{Content: "gone"}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "keep"}, Body: ast.Body{Content: "modified body"}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	assertSectionHeading(t, got.Sections[0], "keep")
	assertSectionHeading(t, got.Sections[1], "remove")
	if got.Sections[0].Body.Content != "modified body" {
		t.Fatalf("expected modified body for keep, got %q", got.Sections[0].Body.Content)
	}
	if got.Sections[1].Body.Content != "gone" {
		t.Fatalf("expected original body for remove, got %q", got.Sections[1].Body.Content)
	}
}

func TestMergeModifiedBodyReplacesOriginal(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Content: "a original"}},
			{Heading: &ast.Heading{Level: 1, Text: "b"}, Body: ast.Body{Content: "b original"}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Content: "a modified"}},
			{Heading: &ast.Heading{Level: 1, Text: "b"}, Body: ast.Body{Content: "b modified"}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	if got.Sections[0].Body.Content != "a modified" {
		t.Fatalf("expected a modified, got %q", got.Sections[0].Body.Content)
	}
	if got.Sections[1].Body.Content != "b modified" {
		t.Fatalf("expected b modified, got %q", got.Sections[1].Body.Content)
	}
}

func TestMergeAddsNewSectionsFromModified(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "first"}, Body: ast.Body{Content: "first body"}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "first"}, Body: ast.Body{Content: "first body"}},
			{Heading: &ast.Heading{Level: 1, Text: "second"}, Body: ast.Body{Content: "second body"}},
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
			{Heading: nil, Body: ast.Body{Content: "original pre"}},
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Content: "a body"}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: nil, Body: ast.Body{Content: "modified pre"}},
			{Heading: &ast.Heading{Level: 1, Text: "a"}, Body: ast.Body{Content: "a modified"}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	if got.Sections[0].Heading != nil {
		t.Fatalf("expected implicit section first")
	}
	if got.Sections[0].Body.Content != "modified pre" {
		t.Fatalf("expected modified pre, got %q", got.Sections[0].Body.Content)
	}
	if got.Sections[1].Body.Content != "a modified" {
		t.Fatalf("expected a modified, got %q", got.Sections[1].Body.Content)
	}
}

func TestMergeDuplicateHeadingsByPosition(t *testing.T) {
	original := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Content: "dup1 original"}},
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Content: "dup2 original"}},
		},
	}
	modified := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Content: "dup1 modified"}},
			{Heading: &ast.Heading{Level: 1, Text: "dup"}, Body: ast.Body{Content: "dup2 modified"}},
		},
	}

	got := MergeDocuments(original, modified)
	if len(got.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got.Sections))
	}
	if got.Sections[0].Body.Content != "dup1 modified" {
		t.Fatalf("expected dup1 modified, got %q", got.Sections[0].Body.Content)
	}
	if got.Sections[1].Body.Content != "dup2 modified" {
		t.Fatalf("expected dup2 modified, got %q", got.Sections[1].Body.Content)
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
