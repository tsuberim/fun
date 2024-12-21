package internal

import (
	"fmt"
	tree_sitter_fun "fun/tree-sitter-fun/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	"os"
	"path"
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
	parser      *tree_sitter.Parser
	evaluator   *Evaluator
	inferer     *Inferrer
	Modules     map[string]*Module
	env         *Env
	currentPath string
}

func NewProgram() (*Program, error) {
	parser := tree_sitter.NewParser()
	err := parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_fun.Language()))
	if err != nil {
		return nil, err
	}

	p := &Program{
		parser:      parser,
		Modules:     map[string]*Module{},
		currentPath: "",
	}
	p.evaluator = NewEvaluator(p)
	p.inferer = NewInferrer(p)
	p.env = NewStdEnv(p)
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
	source, err := os.ReadFile(path.Join(path.Dir(p.currentPath), importPath))
	if err != nil {
		return nil, err
	}

	return p.Run(source, importPath)
}

func (p *Program) Run(source []byte, importPath string) (*Module, error) {
	if importPath != RootModule {
		prevPath := p.currentPath
		p.currentPath = importPath
		defer func() {
			p.currentPath = prevPath
		}()
	}

	tree := p.parser.Parse(source, nil)
	node := tree.RootNode()

	// parse
	expr, err := fromNode(node, source)
	if err != nil {
		return nil, err
	}

	// type check
	_, t, err := p.inferer.Infer(expr, p.env.Types())
	if err != nil {
		return nil, err
	}
	scheme := generalize(t)

	// evaluate
	val, err := p.evaluator.Eval(expr, p.env.Values())
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
