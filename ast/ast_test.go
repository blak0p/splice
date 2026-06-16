package ast_test

import (
	"reflect"
	"testing"

	"github.com/blak0p/splice/ast"
)

func TestBlockKindsAndInterface(t *testing.T) {
	var _ ast.Block = ast.Paragraph{}
	var _ ast.Block = ast.List{}
	var _ ast.Block = ast.Table{}
	var _ ast.Block = ast.CodeBlock{}

	p := ast.Paragraph{ContentLines: []string{"line1", "line2"}}
	if p.Kind() != ast.KindParagraph {
		t.Errorf("expected Paragraph kind, got %s", p.Kind())
	}
	if !reflect.DeepEqual(p.Lines(), []string{"line1", "line2"}) {
		t.Errorf("expected lines, got %v", p.Lines())
	}

	l := ast.List{ContentLines: []string{"- item 1"}}
	if l.Kind() != ast.KindList {
		t.Errorf("expected List kind, got %s", l.Kind())
	}

	tbl := ast.Table{ContentLines: []string{"| col |"}}
	if tbl.Kind() != ast.KindTable {
		t.Errorf("expected Table kind, got %s", tbl.Kind())
	}

	cb := ast.CodeBlock{ContentLines: []string{"```go"}}
	if cb.Kind() != ast.KindCodeBlock {
		t.Errorf("expected CodeBlock kind, got %s", cb.Kind())
	}
}

func TestBodyLinesHelper(t *testing.T) {
	// Case 1: Multiple blocks (checking empty line insertion)
	body1 := ast.Body{
		Blocks: []ast.Block{
			ast.Paragraph{ContentLines: []string{"para line"}},
			ast.List{ContentLines: []string{"list line"}},
		},
	}
	want1 := []string{"para line", "", "list line"}
	got1 := body1.Lines()
	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("expected %v, got %v", want1, got1)
	}

	// Case 2: Empty blocks
	body2 := ast.Body{
		Blocks: []ast.Block{},
	}
	got2 := body2.Lines()
	if len(got2) != 0 {
		t.Errorf("expected empty slice, got %v", got2)
	}

	// Case 3: Single block (no empty lines should be inserted)
	body3 := ast.Body{
		Blocks: []ast.Block{
			ast.Paragraph{ContentLines: []string{"single block line"}},
		},
	}
	want3 := []string{"single block line"}
	got3 := body3.Lines()
	if !reflect.DeepEqual(got3, want3) {
		t.Errorf("expected %v, got %v", want3, got3)
	}
}
