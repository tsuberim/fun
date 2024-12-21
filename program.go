package main

import (
	"fmt"
	tree_sitter_fun "fun/tree-sitter-fun/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	"os"
)

const RootModule = "<root>"

type Module struct {
	ImportPath string
	Expr
	Val
	Type *Scheme
}

func (m *Module) Pretty(indent int) string {
	return fmt.Sprintf("%s : %s", m.Val.Pretty(0), m.Type.Pretty(0))
}

type Program struct {
	parser    *tree_sitter.Parser
	evaluator *Evaluator
	inferer   *Inferrer
	Modules   map[string]*Module
}

func NewProgram() (*Program, error) {
	parser := tree_sitter.NewParser()
	err := parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_fun.Language()))
	if err != nil {
		return nil, err
	}

	p := &Program{
		parser:  parser,
		Modules: map[string]*Module{},
	}
	p.evaluator = NewEvaluator(p)
	p.inferer = NewInferrer(p)
	return p, nil
}

func (p *Program) Import(importPath string) (*Module, error) {
	mod, has := p.Modules[importPath]
	if has {
		return mod, nil
	}
	importedMod, err := p.importModule(importPath)
	if err != nil {
		return nil, err
	}
	mod = importedMod
	p.Modules[importPath] = mod
	return mod, nil
}

func (p *Program) importModule(importPath string) (*Module, error) {
	source, err := os.ReadFile(importPath)
	if err != nil {
		return nil, err
	}

	return p.Run(source, importPath)
}

func (p *Program) Run(source []byte, importPath string) (*Module, error) {
	tree := p.parser.Parse(source, nil)
	node := tree.RootNode()

	// parse
	expr, err := fromNode(node, source)
	if err != nil {
		return nil, err
	}

	// type check
	_, t, err := p.inferer.Infer(expr, typeEnv)
	if err != nil {
		return nil, err
	}
	scheme := generalize(t)

	// evaluate
	val, err := p.evaluator.Eval(expr, p.evaluator.stdlib)
	if err != nil {
		return nil, err
	}

	return &Module{
		ImportPath: importPath,
		Expr:       expr,
		Val:        val,
		Type:       scheme,
	}, nil
}
