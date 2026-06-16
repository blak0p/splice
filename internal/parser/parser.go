package parser

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/blak0p/splice/internal/ast"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/markdown"
)

// Parse converts a Markdown string into a Document by extracting sections
// using tree-sitter-markdown. Sections are flattened from the nested tree-sitter
// structure into a top-down ordered list.
func Parse(input string) (*ast.Document, error) {
	if len(input) == 0 {
		return &ast.Document{}, nil
	}

	if !utf8.ValidString(input) {
		return nil, fmt.Errorf("parse markdown: input is not valid UTF-8")
	}

	src := []byte(input)
	tree, err := markdown.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil, fmt.Errorf("parse markdown: %w", err)
	}

	doc := &ast.Document{}
	doc.Sections = walkNode(tree.BlockTree().RootNode(), src)
	return doc, nil
}

// walkNode recursively flattens nested section nodes into an ordered list.
// tree-sitter-markdown nests sections hierarchically (H1 contains H2 sections).
// We flatten them: each section appears in order, regardless of nesting depth.
func walkNode(node *sitter.Node, src []byte) []ast.Section {
	var sections []ast.Section

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		if child.Type() != "section" {
			continue
		}

		s := extractSection(child, src)
		if s != nil {
			sections = append(sections, *s)
		}

		// Collect child sections recursively (flattening)
		childSections := flattenChildSections(child, src)
		sections = append(sections, childSections...)
	}

	return sections
}

// flattenChildSections extracts nested section nodes from a section.
func flattenChildSections(node *sitter.Node, src []byte) []ast.Section {
	var sections []ast.Section

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() != "section" {
			continue
		}

		s := extractSection(child, src)
		if s != nil {
			sections = append(sections, *s)
		}

		more := flattenChildSections(child, src)
		sections = append(sections, more...)
	}

	return sections
}

// extractSection pulls heading and body from a section node.
// Returns nil if the section has no content (pure grouping node).
func extractSection(node *sitter.Node, src []byte) *ast.Section {
	var heading *ast.Heading
	var blocks []ast.Block

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		switch child.Type() {
		case "atx_heading", "setext_heading":
			heading = parseHeading(child, src)

		case "section":
			// Handled by flattenChildSections — skip here

		default:
			content := child.Content(src)
			content = strings.TrimSuffix(content, "\n")
			if content == "" {
				continue
			}

			blocks = append(blocks, mapNodeToBlock(child, src))
		}
	}

	if heading == nil && len(blocks) == 0 {
		return nil
	}

	return &ast.Section{
		Heading: heading,
		Body:    ast.Body{Blocks: blocks},
	}
}

func mapNodeToBlock(child *sitter.Node, src []byte) ast.Block {
	content := child.Content(src)
	content = strings.TrimSuffix(content, "\n")
	var lines []string
	if content != "" {
		lines = strings.Split(content, "\n")
	}

	switch child.Type() {
	case "paragraph":
		return ast.Paragraph{ContentLines: lines}
	case "list":
		return ast.List{ContentLines: lines}
	case "pipe_table", "table":
		return ast.Table{ContentLines: lines}
	case "fenced_code_block", "indented_code_block":
		return ast.CodeBlock{ContentLines: lines}
	default:
		return ast.Paragraph{ContentLines: lines}
	}
}

// parseHeading extracts level and text from a heading node.
// tree-sitter-markdown stores heading text inside an "inline" child node.
func parseHeading(node *sitter.Node, src []byte) *ast.Heading {
	level := 0
	text := ""

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		switch child.Type() {
		case "atx_h1_marker":
			level = 1
		case "atx_h2_marker":
			level = 2
		case "atx_h3_marker":
			level = 3
		case "atx_h4_marker":
			level = 4
		case "atx_h5_marker":
			level = 5
		case "atx_h6_marker":
			level = 6
		case "inline":
			text = strings.TrimSpace(child.Content(src))
		}
	}

	if level == 0 {
		level = 1
	}

	return &ast.Heading{
		Level: level,
		Text:  strings.TrimSpace(text),
	}
}