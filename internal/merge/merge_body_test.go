package merge

import (
	"reflect"
	"testing"
)

func TestMergeBody_NewLineAppended(t *testing.T) {
	orig := []string{"Keep this line."}
	mod := []string{"Keep this line.", "New line at end."}
	want := []string{"Keep this line.", "New line at end."}
	got := mergeBody(orig, mod)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMergeBody_EditedLineByPosition(t *testing.T) {
	orig := []string{"Old text here."}
	mod := []string{"New text here."}
	want := []string{"New text here."}
	got := mergeBody(orig, mod)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMergeBody_DeletedLine(t *testing.T) {
	orig := []string{"Line one.", "Line two.", "Line three."}
	mod := []string{"Line one.", "Line three."}
	want := []string{"Line one.", "Line three."}
	got := mergeBody(orig, mod)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMergeBody_AllUnchanged(t *testing.T) {
	orig := []string{"Identical content."}
	mod := []string{"Identical content."}
	want := []string{"Identical content."}
	got := mergeBody(orig, mod)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMergeBody_FuzzyMatchAboveThreshold(t *testing.T) {
	orig := []string{"beta"}
	mod := []string{"Beta modified."}
	want := []string{"Beta modified."}
	got := mergeBody(orig, mod)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMergeBody_FuzzyMatchBelowThreshold(t *testing.T) {
	orig := []string{"Alpha."}
	mod := []string{" totally unrelated"}
	want := []string{" totally unrelated"}
	got := mergeBody(orig, mod)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMergeBody_EmptyBodies(t *testing.T) {
	got := mergeBody([]string{}, []string{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}

	got = mergeBody([]string{}, []string{"new"})
	want := []string{"new"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMergeBody_Mixed(t *testing.T) {
	orig := []string{"Alpha.", "Beta.", "Gamma."}
	mod := []string{"Alpha.", "Beta modified.", "Delta."}
	want := []string{"Alpha.", "Beta modified.", "Delta."}
	got := mergeBody(orig, mod)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
