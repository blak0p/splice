# splice

[![Go Reference](https://pkg.go.dev/badge/github.com/blak0p/splice.svg)](https://pkg.go.dev/github.com/blak0p/splice)
[![Tests](https://github.com/blak0p/splice/actions/workflows/test.yml/badge.svg)](https://github.com/blak0p/splice/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**splice** is a Go library for section-aware 3-way merging of Markdown documents. It parses Markdown into an AST, matches sections by heading, applies changes at block level, and renders the result back — all without treating your document as plain text.

## Installation

```bash
go get github.com/blak0p/splice
```

Requires Go 1.26+ and CGO (backed by [tree-sitter-markdown](https://github.com/tree-sitter-grammars/tree-sitter-markdown)).

## API

### Parse

```go
doc, err := splice.Parse("# Hello\n\nThis is a paragraph.")
if err != nil {
    log.Fatal(err)
}
// doc.Sections[0].Heading.Text == "Hello"
// doc.Sections[0].Body.Blocks[0].Kind() == ast.KindParagraph
```

### Render

```go
doc, _ := splice.Parse("# Hello\n\nWorld")
output := splice.Render(doc)
fmt.Println(output)
// Output:
// # Hello
//
// World
```

### Merge

```go
original := "# Changelog\n\n## [1.0.0]\n\n- Initial release\n"
modified := "# Changelog\n\n## [1.0.0]\n\n- Initial release\n\n## [1.1.0]\n\n- Dark mode\n"

result, err := splice.Merge(context.Background(), original, modified)
if err != nil {
    log.Fatal(err)
}
// Sections only in modified are inserted after their nearest match.
// New section [1.1.0] appears right after [1.0.0].
```

### MergeAST

```go
origDoc, _ := splice.Parse("# A\n\nContent A")
modDoc, _ := splice.Parse("# A\n\nContent B\n\n## C\n\nNew section")

merged := splice.MergeAST(origDoc, modDoc)
output := splice.Render(merged)
// Original heading A keeps modified body ("Content B").
// New section C is inserted after its nearest match (A).
```

### Options

```go
result, err := splice.Merge(ctx, original, modified,
    splice.WithThreshold(0.5),
    splice.WithCaseInsensitive(true),
)
```

Custom block merger:

```go
result, err := splice.Merge(ctx, original, modified,
    splice.WithBlockMerger(func(orig, mod ast.Block) (ast.Block, bool) {
        // Merge tables row by row instead of replacing.
        if orig.Kind() == ast.KindTable && mod.Kind() == ast.KindTable {
            // custom table merge logic
            return mergedBlock, true
        }
        return nil, false // fallback to default
    }),
)
```

## How it works

1. **Parse** — each document is parsed into a section-based AST via tree-sitter-markdown.
2. **Merge** — sections are matched by heading. Modified sections keep their changes; new sections are inserted after their nearest matched predecessor; untouched sections are preserved. Body blocks can be merged with custom strategies.
3. **Render** — the merged AST is serialized back to Markdown, preserving block structure (tables, lists, code blocks).

## License

MIT — see [LICENSE](LICENSE).
