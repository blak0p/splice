# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] — 2026-06-16

### Added

- Public API: `Parse`, `Render`, `Merge`, `MergeAST` with functional options.
- `ast` package: `Document`, `Section`, `Heading`, `Body`, `Block` interface, `BlockKind`.
- Block-level AST: `Paragraph`, `List`, `Table`, `CodeBlock` types.
- Merge options: `WithThreshold`, `WithCaseInsensitive`, `WithBlockMerger`.
- Section-aware 3-way merge with block-level body merging.
- Custom block merger support (e.g., table-aware merge).
- CGO-backed parser via tree-sitter-markdown.
- Full Go module with `go get github.com/blak0p/splice`.

[0.1.0]: https://github.com/blak0p/splice/releases/tag/v0.1.0
