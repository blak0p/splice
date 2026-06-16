package ast

// Document represents a parsed Markdown document as a sequence of sections.
type Document struct {
	Sections []Section
}

// Section is a logical chunk of the document. A nil Heading indicates
// pre-heading (implicit) content.
type Section struct {
	Heading *Heading
	Body    Body
}

// Heading stores the level (1-6) and normalized text of a Markdown heading.
type Heading struct {
	Level int
	Text  string
}

// Body holds the AST blocks of Markdown content that belong to a section.
type Body struct {
	Blocks []Block
}

// Lines returns the lines of all blocks in the body, separated by empty strings
// for backwards compatibility.
func (b Body) Lines() []string {
	var lines []string
	for i, block := range b.Blocks {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, block.Lines()...)
	}
	return lines
}

// BlockKind represents the specific semantic kind of Markdown block.
type BlockKind string

const (
	KindParagraph BlockKind = "Paragraph"
	KindList      BlockKind = "List"
	KindTable     BlockKind = "Table"
	KindCodeBlock BlockKind = "CodeBlock"
)

// Block is the interface implemented by all concrete AST blocks.
type Block interface {
	Kind() BlockKind
	Lines() []string
}

// Paragraph represents a block of narrative text.
type Paragraph struct {
	ContentLines []string
}

func (p Paragraph) Kind() BlockKind { return KindParagraph }
func (p Paragraph) Lines() []string { return p.ContentLines }

// List represents a bulleted or numbered list.
type List struct {
	ContentLines []string
}

func (l List) Kind() BlockKind { return KindList }
func (l List) Lines() []string { return l.ContentLines }

// Table represents a Markdown table.
type Table struct {
	ContentLines []string
}

func (t Table) Kind() BlockKind { return KindTable }
func (t Table) Lines() []string { return t.ContentLines }

// CodeBlock represents a fenced or indented block of code.
type CodeBlock struct {
	ContentLines []string
}

func (c CodeBlock) Kind() BlockKind { return KindCodeBlock }
func (c CodeBlock) Lines() []string { return c.ContentLines }
