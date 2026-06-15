package splice

import (
	"context"
	"fmt"

	"github.com/blak0p/splice/internal/merge"
	"github.com/blak0p/splice/internal/parser"
	"github.com/blak0p/splice/internal/renderer"
)

// Merge fusiona dos versiones de un documento Markdown a nivel de secciones.
// Toma el documento original y una versión modificada, y produce un resultado
// que preserva secciones no tocadas, aplica cambios donde hubo modificaciones,
// y agrega secciones nuevas (de modified) que no existían en original.
func Merge(ctx context.Context, original string, modified string) (string, error) {
	origDoc, err := parser.Parse(original)
	if err != nil {
		return "", fmt.Errorf("splice: parse original: %w", err)
	}

	modDoc, err := parser.Parse(modified)
	if err != nil {
		return "", fmt.Errorf("splice: parse modified: %w", err)
	}

	merged := merge.MergeDocuments(origDoc, modDoc)
	return renderer.Render(merged), nil
}
