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

// Body holds the raw Markdown content that belongs to a section.
type Body struct {
	Content string
}
