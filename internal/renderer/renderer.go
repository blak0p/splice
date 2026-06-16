package renderer

import (
	"strings"

	"github.com/blak0p/splice/ast"
)

// Render converts a Document back into Markdown text.
// Sections are emitted in order, separated by blank lines. Sections with a
// heading emit the heading markers first, followed by the section body.
// Pre-heading sections (nil heading) emit only their body content.
func Render(doc *ast.Document) string {
	if doc == nil || len(doc.Sections) == 0 {
		return ""
	}

	var sections []string
	for _, section := range doc.Sections {
		var sb strings.Builder
		if section.Heading != nil {
			markers := strings.Repeat("#", section.Heading.Level)
			sb.WriteString(markers)
			sb.WriteByte(' ')
			sb.WriteString(section.Heading.Text)
		}

		var blockStrings []string
		for _, block := range section.Body.Blocks {
			if isBlockEmpty(block) {
				continue
			}
			blockStrings = append(blockStrings, strings.Join(block.Lines(), "\n"))
		}

		if len(blockStrings) > 0 {
			if section.Heading != nil {
				sb.WriteString("\n\n")
			}
			sb.WriteString(strings.Join(blockStrings, "\n\n"))
		}
		sections = append(sections, sb.String())
	}

	return strings.Join(sections, "\n\n") + "\n"
}

func isBlockEmpty(b ast.Block) bool {
	lines := b.Lines()
	if len(lines) == 0 {
		return true
	}
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			return false
		}
	}
	return true
}
