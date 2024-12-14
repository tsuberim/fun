package tree_sitter_fun_test

import (
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_fun "github.com/tsuberim/fun/bindings/go"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_fun.Language())
	if language == nil {
		t.Errorf("Error loading Fun grammar")
	}
}
