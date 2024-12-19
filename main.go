package main

import (
	"fmt"
	tree_sitter_fun "fun/tree-sitter-fun/bindings/go"
	"github.com/maxott/go-repl"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	"os"
)

var stdlib = map[string]Val{
	"+": &Builtin{
		Name: "+",
		Impl: func(args []Val) (Val, error) {
			sum := 0
			for _, arg := range args {
				i, ok := arg.(*Int)
				if !ok {
					return nil, fmt.Errorf("invalid sum value type %t", arg)
				}

				sum += i.Value
			}

			return &Int{Value: sum}, nil
		},
	},
}

var typeEnv = &TypeEnv{Types: map[string]*Scheme{
	"+": {
		Forall: nil,
		Type: &TypeCons{
			Name: lambdaConsName,
			Args: []Type{
				&TypeCons{
					Name: "int",
					Args: nil,
				},
				&TypeCons{
					Name: "int",
					Args: nil,
				},
				&TypeCons{
					Name: "int",
					Args: nil,
				},
			},
		},
	},
}}

func main() {
	if len(os.Args) > 1 {
		filename := os.Args[1]
		source, err := os.ReadFile(filename)
		if err != nil {
			println(err.Error())
		}
		println(eval(source))
	} else {
		r := repl.NewRepl(&ReplHandler{})
		err := r.Loop()
		if err != nil {
			panic(err)
		}
	}
}

func eval(source []byte) string {
	parser := tree_sitter.NewParser()
	defer parser.Close()
	err := parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_fun.Language()))
	if err != nil {
		panic(err)
	}

	tree := parser.Parse(source, nil)
	node := tree.RootNode()
	expr, err := fromNode(node, source)
	if err != nil {
		return fmt.Sprintf("ParseError: %s", err)
	}

	inferrer := NewInferrer()
	_, t, err := inferrer.Infer(expr, typeEnv)
	if err != nil {
		return fmt.Sprintf("TypeError: %s", err)
	}
	scheme := generalize(t)

	val, err := Eval(expr, stdlib)
	if err != nil {
		return fmt.Sprintf("ValueError: %s", err)
	}

	return fmt.Sprintf("%s : %s", val.Pretty(0), scheme.Pretty(0))
}

type ReplHandler struct{}

func (r *ReplHandler) Prompt() string {
	return ">"
}

func (r *ReplHandler) Eval(buffer string) string {
	source := []byte(buffer)
	return eval(source)
}

func (r *ReplHandler) Tab(buffer string) string {
	// TODO: use tree-sitter to complete?
	return ""
}
