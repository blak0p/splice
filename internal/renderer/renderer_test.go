package renderer

import (
	"strings"
	"testing"

	"github.com/blak0p/splice/ast"
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
	assertSection(t, doc2.Sections[0], 1, "Intro", []string{"intro body"})
	assertSection(t, doc2.Sections[1], 2, "Details", []string{"details body"})
}

func TestRenderPreHeadingThenHeadings(t *testing.T) {
	doc := &ast.Document{
		Sections: []ast.Section{
			{Heading: nil, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"pre-heading content"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "First"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"body"}}}}},
		},
	}

	got := Render(doc)
	want := "pre-heading content\n\n# First\n\nbody\n"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderMultiSectionPreservesBoundaries(t *testing.T) {
	doc := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "One"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"body one"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "Two"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"body two"}}}}},
			{Heading: &ast.Heading{Level: 1, Text: "Three"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"body three"}}}}},
		},
	}

	got := Render(doc)
	want := "# One\n\nbody one\n\n# Two\n\nbody two\n\n# Three\n\nbody three\n"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderTrailingNewline(t *testing.T) {
	doc := &ast.Document{
		Sections: []ast.Section{
			{Heading: &ast.Heading{Level: 1, Text: "Title"}, Body: ast.Body{Blocks: []ast.Block{ast.Paragraph{ContentLines: []string{"body"}}}}},
		},
	}

	got := Render(doc)
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("expected trailing newline, got %q", got)
	}
}

func TestRenderMultipleBlocks(t *testing.T) {
	doc := &ast.Document{
		Sections: []ast.Section{
			{
				Heading: &ast.Heading{Level: 1, Text: "Blocks"},
				Body: ast.Body{
					Blocks: []ast.Block{
						ast.Paragraph{ContentLines: []string{"Paragraph line 1.", "Paragraph line 2."}},
						ast.List{ContentLines: []string{"- Item 1", "- Item 2"}},
						ast.Table{ContentLines: []string{"| Col 1 |", "|---|", "| Val 1 |"}},
						ast.CodeBlock{ContentLines: []string{"```go", "main()", "```"}},
					},
				},
			},
		},
	}

	got := Render(doc)
	want := `# Blocks

Paragraph line 1.
Paragraph line 2.

- Item 1
- Item 2

| Col 1 |
|---|
| Val 1 |

` + "```go" + `
main()
` + "```" + `
`
	if got != want {
		t.Fatalf("expected:\n%q\n\ngot:\n%q", want, got)
	}
}

func TestRenderEmptyBlocks(t *testing.T) {
	doc := &ast.Document{
		Sections: []ast.Section{
			{
				Heading: &ast.Heading{Level: 1, Text: "Empty"},
				Body: ast.Body{
					Blocks: []ast.Block{
						ast.Paragraph{ContentLines: []string{}},
						ast.Paragraph{ContentLines: []string{""}},
						ast.List{ContentLines: []string{"- item"}},
					},
				},
			},
		},
	}

	got := Render(doc)
	want := `# Empty

- item
`
	if got != want {
		t.Fatalf("expected:\n%q\n\ngot:\n%q", want, got)
	}
}

func assertSection(t *testing.T, s ast.Section, level int, text string, body []string) {
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
	if len(s.Body.Lines()) != len(body) {
		t.Fatalf("expected body %v, got %v", body, s.Body.Lines())
	}
	for i, line := range body {
		if s.Body.Lines()[i] != line {
			t.Fatalf("expected body line %q, got %q", line, s.Body.Lines()[i])
		}
	}
}
