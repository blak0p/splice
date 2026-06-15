package renderer

import (
	"strings"

	"github.com/blak0p/splice/internal/ast"
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
		if len(section.Body.Lines) > 0 {
			if section.Heading != nil {
				sb.WriteString("\n\n")
			}
			sb.WriteString(strings.Join(section.Body.Lines, "\n"))
		}
		sections = append(sections, sb.String())
	}

	return strings.Join(sections, "\n\n") + "\n"
}
