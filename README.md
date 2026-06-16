# splice

[![Go Reference](https://pkg.go.dev/badge/github.com/blak0p/splice.svg)](https://pkg.go.dev/github.com/blak0p/splice)
[![Tests](https://github.com/blak0p/splice/actions/workflows/test.yml/badge.svg)](https://github.com/blak0p/splice/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**splice** is a Go library for section-aware 3-way merging of Markdown documents. It parses Markdown into an AST, matches sections by heading, applies changes at block level, and renders the result — all without treating your document as plain text.

## Installation

```bash
go get github.com/blak0p/splice
```

Requires Go 1.26+ and CGO (backed by [tree-sitter-markdown](https://github.com/tree-sitter-grammars/tree-sitter-markdown)).

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/blak0p/splice"
)

func main() {
    ctx := context.Background()

    original := "# Changelog\n\n## [1.0.0] — 2026-01-15\n\n### Added\n\n- Initial release\n"
    modified := "# Changelog\n\n## [1.0.0] — 2026-01-15\n\n### Added\n\n- Initial release\n\n## [1.1.0] — 2026-03-01\n\n### Added\n\n- Dark mode\n"

    result, err := splice.Merge(ctx, original, modified)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

The merge preserves the `[1.0.0]` section unchanged, applies any body changes from modified, and inserts the new `[1.1.0]` section after its nearest matched predecessor.

### With options

```go
result, err := splice.Merge(ctx, original, modified,
    splice.WithThreshold(0.5),
    splice.WithCaseInsensitive(true),
)
```

### Working with the AST

```go
origDoc, _ := splice.Parse(original)
modDoc, _ := splice.Parse(modified)

mergedDoc := splice.MergeAST(origDoc, modDoc)

output := splice.Render(mergedDoc)
```

## API

| Function | Description |
|----------|-------------|
| `Parse(input string) (*ast.Document, error)` | Parse Markdown into an AST. |
| `Render(doc *ast.Document) string` | Render an AST back to Markdown. |
| `Merge(ctx, original, modified string, opts ...Option) (string, error)` | 3-way merge two Markdown documents. |
| `MergeAST(origDoc, modDoc *ast.Document, opts ...Option) *ast.Document` | Merge two pre-parsed ASTs. |

### Options

- `WithThreshold(t float64)` — similarity threshold for block matching (0.0–1.0, default 0.7).
- `WithCaseInsensitive(enabled bool)` — case-insensitive heading matching.
- `WithBlockMerger(fn func(orig, mod ast.Block) (ast.Block, bool))` — custom per-block merger for advanced use cases (e.g., table-aware merge).

See the [package docs](https://pkg.go.dev/github.com/blak0p/splice) on pkg.go.dev for the full API reference.

## How it works

1. **Parse** — each document is parsed into a section-based AST via tree-sitter-markdown.
2. **Merge** — sections are matched by heading. Modified sections keep their changes; new sections are inserted after their nearest matched predecessor; untouched sections are preserved. Body blocks can be merged with custom strategies.
3. **Render** — the merged AST is serialized back to Markdown, preserving block structure (tables, lists, code blocks).

## License

MIT — see [LICENSE](LICENSE).
