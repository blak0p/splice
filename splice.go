package splice

import (
	"context"
	"fmt"

	"github.com/blak0p/splice/ast"
	"github.com/blak0p/splice/internal/merge"
	"github.com/blak0p/splice/internal/parser"
	"github.com/blak0p/splice/internal/renderer"
)

// Option configures merge behavior.
type Option func(*config) error

type config struct {
	threshold       float64
	caseInsensitive bool
	blockMerger     func(orig, mod ast.Block) (ast.Block, bool)
}

func defaultConfig() *config {
	return &config{
		threshold: 0.7,
	}
}

// WithThreshold sets the similarity threshold for block matching (0.0 to 1.0).
func WithThreshold(t float64) Option {
	return func(c *config) error {
		c.threshold = t
		return nil
	}
}

// WithCaseInsensitive enables case-insensitive heading matching.
func WithCaseInsensitive(enabled bool) Option {
	return func(c *config) error {
		c.caseInsensitive = enabled
		return nil
	}
}

// WithBlockMerger sets a custom block-level merger function.
func WithBlockMerger(fn func(orig, mod ast.Block) (ast.Block, bool)) Option {
	return func(c *config) error {
		c.blockMerger = fn
		return nil
	}
}

// Parse compiles a markdown string into an AST document.
func Parse(input string) (*ast.Document, error) {
	return parser.Parse(input)
}

// Render serializes an AST document to a markdown string.
func Render(doc *ast.Document) string {
	return renderer.Render(doc)
}

// MergeAST merges two pre-parsed AST documents and returns the merged AST.
func MergeAST(origDoc, modDoc *ast.Document, opts ...Option) *ast.Document {
	cfg := defaultConfig()
	for _, opt := range opts {
		_ = opt(cfg)
	}
	return merge.MergeAST(origDoc, modDoc, toMergeConfig(cfg))
}

// Merge fusiona dos versiones de un documento Markdown a nivel de secciones.
// Toma el documento original y una versión modificada, y produce un resultado
// que preserva secciones no tocadas, aplica cambios donde hubo modificaciones,
// y agrega secciones nuevas (de modified) que no existían en original.
func Merge(ctx context.Context, original, modified string, opts ...Option) (string, error) {
	origDoc, err := parser.Parse(original)
	if err != nil {
		return "", fmt.Errorf("splice: parse original: %w", err)
	}

	modDoc, err := parser.Parse(modified)
	if err != nil {
		return "", fmt.Errorf("splice: parse modified: %w", err)
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		_ = opt(cfg)
	}

	merged := merge.MergeAST(origDoc, modDoc, toMergeConfig(cfg))
	return renderer.Render(merged), nil
}

func toMergeConfig(c *config) *merge.Config {
	return &merge.Config{
		Threshold:       c.threshold,
		CaseInsensitive: c.caseInsensitive,
		BlockMerger:     c.blockMerger,
	}
}
