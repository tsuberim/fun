package main

import (
	"fmt"
	tree_sitter_fun "fun/tree-sitter-fun/bindings/go"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	//tree_sitter_fun "github.com/tsuberim/fun/tree-sitter-fun/bindings/go"
	//tree_sitter_fun "github.com/tsuberim/fun/bindings/go"
)

func main() {
	// Create a new LLVM IR module.
	m := ir.NewModule()
	hello := constant.NewCharArrayFromString("Hello, world!\n\x00")
	str := m.NewGlobalDef("str", hello)
	// Add external function declaration of puts.
	puts := m.NewFunc("puts", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")
	// Cast *[15]i8 to *i8.
	zero := constant.NewInt(types.I64, 0)
	gep := constant.NewGetElementPtr(hello.Typ, str, zero, zero)
	entry.NewCall(puts, gep)
	entry.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)

	code := []byte("hello")

	parser := tree_sitter.NewParser()
	defer parser.Close()
	err := parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_fun.Language()))
	if err != nil {
		panic(err)
	}

	tree := parser.Parse(code, nil)
	defer tree.Close()

	root := tree.RootNode()
	fmt.Println(root.ToSexp())
}
