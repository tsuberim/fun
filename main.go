package main

import (
	"bufio"
	"fmt"
	tree_sitter_fun "fun/tree-sitter-fun/bindings/go"
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

func main() {
	if len(os.Args) > 1 {
		filename := os.Args[0]
		source, err := os.ReadFile(filename)
		parser := tree_sitter.NewParser()
		defer parser.Close()
		err = parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_fun.Language()))
		if err != nil {
			panic(err)
		}

		tree := parser.Parse(source, nil)
		node := tree.RootNode()
		expr, err := fromNode(node, source)
		if err != nil {
			panic(err)
		}

		fmt.Println(expr.Pretty(0))
	} else {
		err := repl()
		if err != nil {
			panic(err)
		}
	}

}

func repl() error {
	reader := bufio.NewReader(os.Stdin)

	parser := tree_sitter.NewParser()
	defer parser.Close()
	err := parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_fun.Language()))
	if err != nil {
		return err
	}

	for {
		print(">")
		bs, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		tree := parser.Parse(bs, nil)

		root := tree.RootNode()
		if root.HasError() {
			println("ERR")
		} else {
			println(root.ToSexp())
		}

		node := root.NamedChild(0)
		expr, err := fromNode(node, bs)
		tree.Close()

		if err != nil {
			println(err.Error())
			continue
		}

		println("Expr: ", expr.Pretty(0))

		val, err := Eval(expr, stdlib)
		if err != nil {
			println(err.Error())
			continue
		}

		println("Value: ", val.Pretty(0))

		inferrer := NewInferrer()
		_, t, err := inferrer.Infer(expr, &TypeEnv{Types: map[string]*Scheme{
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
		}})
		if err != nil {
			println(err.Error())
			continue
		}
		scheme := generalize(t)
		println("Type: ", scheme.Pretty(0))
	}
}
